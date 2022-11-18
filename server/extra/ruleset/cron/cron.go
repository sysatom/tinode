package cron

import (
	"crypto/sha1"
	"fmt"
	"github.com/influxdata/cron"
	"github.com/tinode/chat/server/auth"
	"github.com/tinode/chat/server/extra/cache"
	"github.com/tinode/chat/server/extra/store"
	extraTypes "github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/logs"
	serverStore "github.com/tinode/chat/server/store"
	"github.com/tinode/chat/server/store/types"
	"time"
)

type Rule struct {
	Name   string
	When   string
	Action func(extraTypes.Context) []extraTypes.MsgPayload
}

type Ruleset struct {
	Type      string
	AuthLevel auth.Level
	outCh     chan result
	cronRules []Rule

	Send func(userUid, topicUid types.Uid, out extraTypes.MsgPayload)
}

type result struct {
	name    string
	ctx     extraTypes.Context
	payload extraTypes.MsgPayload
}

// NewCronRuleset New returns a cron rule set
func NewCronRuleset(name string, authLevel auth.Level, rules []Rule) *Ruleset {
	r := &Ruleset{
		Type:      name,
		AuthLevel: authLevel,
		cronRules: rules,
		outCh:     make(chan result, 100),
	}
	return r
}

func (r *Ruleset) Daemon() {
	logs.Info.Println("cron starting...")

	// process cron
	for rule := range r.cronRules {
		logs.Info.Printf("cron %s start...", r.cronRules[rule].Name)
		go r.ruleWorker(r.cronRules[rule])
	}

	// result pipeline
	go r.resultWorker()
}

func (r *Ruleset) ruleWorker(rule Rule) {
	p, err := cron.ParseUTC(rule.When)
	if err != nil {
		logs.Err.Println("cron worker", rule.Name, err)
		return
	}
	nextTime, err := p.Next(time.Now())
	if err != nil {
		logs.Err.Println("cron worker", rule.Name, err)
		return
	}
	for {
		if nextTime.Format("2006-01-02 15:04") == time.Now().Format("2006-01-02 15:04") {
			logs.Info.Printf("cron %s scheduled", rule.Name)
			msgs := func() []result {
				defer func() {
					if rc := recover(); rc != nil {
						logs.Warn.Printf("cron %s ruleWorker recover", rule.Name)
						if v, ok := rc.(error); ok {
							logs.Err.Println(v)
						}
					}
				}()

				items, err := store.Chatbot.OAuthGetAvailable(r.Type)
				if err != nil {
					logs.Err.Println("cron worker", rule.Name, err)
					return nil
				}
				if len(items) > 0 {
					var res []result
					for _, oauth := range items {
						ctx := extraTypes.Context{
							Original: oauth.Topic,
							AsUser:   types.ParseUserId(oauth.Uid),
							Token:    oauth.Token,
						}
						ra := rule.Action(ctx)
						for i := range ra {
							res = append(res, result{
								name:    rule.Name,
								ctx:     ctx,
								payload: ra[i],
							})
						}
					}
					return res
				}

				var res []result
				ra := rule.Action(extraTypes.Context{}) // fixme
				for i := range ra {
					res = append(res, result{
						name:    rule.Name,
						ctx:     extraTypes.Context{}, // fixme
						payload: ra[i],
					})
				}

				return res
			}()
			if len(msgs) > 0 {
				for _, item := range msgs {
					r.outCh <- item
				}
			}
		}
		nextTime, err = p.Next(time.Now())
		if err != nil {
			logs.Err.Println("cron worker", rule.Name, err)
			continue
		}
		time.Sleep(2 * time.Second)
	}
}

func (r *Ruleset) resultWorker() {
	for out := range r.outCh {
		// filter
		res := r.filter(out)
		// pipeline
		r.pipeline(res)
	}
}

func (r *Ruleset) filter(res result) result {
	// user auth record
	_, authLvl, _, _, _ := serverStore.Users.GetAuthRecord(res.ctx.AsUser, "basic")
	switch r.AuthLevel {
	case auth.LevelRoot:
		if authLvl != auth.LevelRoot {
			return result{}
		}
	}

	filterKey := []byte(fmt.Sprintf("cron:%s:%s:filter", res.name, res.ctx.AsUser.UserId()))

	// content hash
	d := un(res.payload)
	s := sha1.New()
	_, _ = s.Write(d)
	hash := s.Sum(nil)

	state := cache.DB.SIsMember(filterKey, hash)
	if state {
		return result{}
	}

	_ = cache.DB.SAdd(filterKey, hash)
	return res
}

func (r *Ruleset) pipeline(res result) {
	if res.payload == nil {
		return
	}
	r.Send(res.ctx.AsUser, types.ParseUserId(res.ctx.Original), res.payload)
}

func un(payload extraTypes.MsgPayload) []byte {
	switch v := payload.(type) {
	case extraTypes.InfoMsg:
		return []byte(v.Title)
	}
	return nil
}
