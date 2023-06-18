package cron

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/influxdata/cron"
	"github.com/tinode/chat/server/auth"
	"github.com/tinode/chat/server/extra/pkg/cache"
	"github.com/tinode/chat/server/extra/store"
	extraTypes "github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/logs"
	serverStore "github.com/tinode/chat/server/store"
	"github.com/tinode/chat/server/store/types"
	"gorm.io/gorm"
	"time"
)

type Rule struct {
	Name   string
	Help   string
	When   string
	Action func(extraTypes.Context) []extraTypes.MsgPayload
}

type Ruleset struct {
	Type      string
	AuthLevel auth.Level
	outCh     chan result
	cronRules []Rule

	Send extraTypes.SendFunc
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
	// process cron
	for rule := range r.cronRules {
		logs.Info.Printf("cron %s start", r.cronRules[rule].Name)
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
			msgs := func() []result {
				defer func() {
					if rc := recover(); rc != nil {
						logs.Warn.Printf("cron %s ruleWorker recover", rule.Name)
						if v, ok := rc.(error); ok {
							logs.Err.Println(v)
						}
					}
				}()

				// bot user
				botUid, _, _, _, _ := serverStore.Users.GetAuthUniqueRecord("basic", fmt.Sprintf("%s_bot", r.Type))

				// all normal users
				users, err := store.Chatbot.GetNormalUsers()
				if err != nil {
					logs.Err.Println(err)
					return nil
				}

				var res []result
				for _, user := range users {
					// check subscription
					uid := serverStore.EncodeUid(int64(user.ID))
					topic := uid.P2PName(botUid)
					sub, err := serverStore.Subs.Get(topic, uid, false)
					if err != nil {
						continue
					}
					if sub == nil || sub.Topic == "" {
						continue
					}

					// get oauth token
					oauth, err := store.Chatbot.OAuthGet(uid, botUid.UserId(), r.Type)
					if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
						continue
					}

					// ctx
					ctx := extraTypes.Context{
						Original: botUid.UserId(),
						AsUser:   uid,
						Token:    oauth.Token,
						RcptTo:   topic,
					}

					// run action
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

	filterKey := fmt.Sprintf("cron:%s:%s:filter", res.name, res.ctx.AsUser.UserId())

	// content hash
	d := un(res.payload)
	s := sha1.New()
	_, _ = s.Write(d)
	hash := s.Sum(nil)

	ctx := context.Background()
	state := cache.DB.SIsMember(ctx, filterKey, hash).Val()
	if state {
		return result{}
	}

	_ = cache.DB.SAdd(ctx, filterKey, hash)
	return res
}

func (r *Ruleset) pipeline(res result) {
	if res.payload == nil {
		return
	}
	r.Send(res.ctx.RcptTo, types.ParseUserId(res.ctx.Original), res.payload)
}

func un(payload extraTypes.MsgPayload) []byte {
	switch v := payload.(type) {
	case extraTypes.TextMsg:
		return []byte(v.Text)
	case extraTypes.InfoMsg:
		return []byte(v.Title)
	case extraTypes.RepoMsg:
		return []byte(*v.FullName)
	case extraTypes.LinkMsg:
		return []byte(v.Url)
	}
	return nil
}
