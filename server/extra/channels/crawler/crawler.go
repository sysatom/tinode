package crawler

import (
	"fmt"
	"github.com/flower-corp/rosedb"
	"github.com/influxdata/cron"
	"github.com/tinode/chat/server/logs"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

type Crawler struct {
	db    *rosedb.RoseDB
	jobs  map[string]Rule
	outCh chan Result

	Send func(channel, name string, out [][]byte)
}

func New() *Crawler {
	path := filepath.Join("./tmp", "rosedb")
	opts := rosedb.DefaultOptions(path)
	db, err := rosedb.Open(opts)
	if err != nil {
		logs.Err.Fatal("open rosedb err: %v", err)
	}

	return &Crawler{
		db:    db,
		jobs:  make(map[string]Rule),
		outCh: make(chan Result, 10),
	}
}

func (s *Crawler) Init(rules ...Rule) error {
	for _, r := range rules {
		// check
		if r.Name == "" {
			continue
		}
		if r.When == "" {
			continue
		}
		if !IsUrl(r.Page.URL) {
			continue
		}

		s.jobs[r.Name] = r
	}
	return nil
}

func (s *Crawler) Run() {
	logs.Info.Println("crawler starting...")

	for name, job := range s.jobs {
		go s.ruleWorker(name, job)
	}

	go s.resultWorker()
}

func (s *Crawler) ruleWorker(name string, r Rule) {
	logs.Info.Println("crawler : crawl...", name)
	p, err := cron.ParseUTC(r.When)
	if err != nil {
		logs.Err.Println(err, name)
		return
	}
	nextTime, err := p.Next(time.Now())
	if err != nil {
		logs.Err.Println(err, name)
		return
	}
	for {
		if nextTime.Format("2006-01-02 15:04") == time.Now().Format("2006-01-02 15:04") {
			logs.Info.Println("crawler : scheduled", name)
			result := func() [][]byte {
				defer func() {
					if r := recover(); r != nil {
						logs.Warn.Println("ruleWorker recover "+name, name)
						if v, ok := r.(error); ok {
							logs.Err.Println(v, name)
						}
					}
				}()
				return r.Run()
			}()
			logs.Info.Println("crawler: result", len(result))
			if len(result) > 0 {
				s.outCh <- Result{
					Name:    name,
					Channel: r.Channel,
					Mode:    r.Mode,
					Result:  result,
				}
			}
		}
		nextTime, err = p.Next(time.Now())
		if err != nil {
			logs.Err.Println(err, name)
			time.Sleep(30 * time.Second)
			continue
		}
		time.Sleep(2 * time.Second)
	}
}

func (s *Crawler) resultWorker() {
	for out := range s.outCh {
		// filter
		diff := s.filter(out.Name, out.Mode, out.Result)
		// send
		s.Send(out.Channel, out.Name, diff)
	}
}

func (s *Crawler) filter(name, mode string, latest [][]byte) [][]byte {
	sentKey := fmt.Sprintf("crawler:%s:sent", name)
	todoKey := fmt.Sprintf("crawler:%s:todo", name)
	sendTimeKey := fmt.Sprintf("crawler:%s:sendtime", name)

	// sent
	old, err := s.db.SMembers([]byte(sentKey))
	if err != nil {
		return [][]byte{}
	}

	// to do
	todo, err := s.db.SMembers([]byte(todoKey))
	if err != nil {
		return [][]byte{}
	}

	// merge
	old = append(old, todo...)

	// diff
	diff := StringSliceDiff(latest, old)

	switch mode {
	case "instant":
		_ = s.db.Set([]byte(sendTimeKey), []byte(strconv.FormatInt(time.Now().Unix(), 10)))
	case "daily":
		sendString, err := s.db.Get([]byte(sendTimeKey))
		if err != nil {
			return [][]byte{}
		}
		oldSend := int64(0)
		if len(sendString) != 0 {
			oldSend, _ = strconv.ParseInt(string(sendString), 10, 64)
		}

		if time.Now().Unix()-oldSend < 24*60*60 {
			for _, item := range diff {
				_ = s.db.SAdd([]byte(todoKey), item)
			}

			return [][]byte{}
		}

		diff = append(diff, todo...)

		_ = s.db.Set([]byte(sendTimeKey), []byte(strconv.FormatInt(time.Now().Unix(), 10)))
	default:
		return [][]byte{}
	}

	// add data
	for _, item := range diff {
		_ = s.db.SAdd([]byte(sentKey), item)
	}

	// clear to do
	_ = s.db.Delete([]byte(todoKey))

	return diff
}

const (
	UrlRegex = `https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`
)

func IsUrl(text string) bool {
	re := regexp.MustCompile("^" + UrlRegex + "$")
	return re.MatchString(text)
}

func StringSliceDiff(s1, s2 [][]byte) [][]byte {
	if len(s1) == 0 {
		return s2
	}
	mb := make(map[string]struct{}, len(s2))
	for _, x := range s2 {
		mb[string(x)] = struct{}{}
	}
	var diff [][]byte
	for _, x := range s1 {
		if _, ok := mb[string(x)]; !ok {
			diff = append(diff, x)
		}
	}
	return diff
}
