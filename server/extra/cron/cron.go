package cron

import (
	"github.com/influxdata/cron"
	"github.com/tinode/chat/server/extra/store"
	extraTypes "github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/logs"
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
	outCh     chan extraTypes.MsgPayload
	cronRules []Rule
}

// NewCronRuleset New returns a cron rule set
func NewCronRuleset(name string, rules []Rule) *Ruleset {
	r := &Ruleset{
		Type:      name,
		cronRules: rules,
		outCh:     make(chan extraTypes.MsgPayload, 100),
	}
	return r
}

func (r *Ruleset) Daemon() {
	logs.Info.Println("cron starting...")

	// process cron
	ctx := extraTypes.Context{}
	for rule := range r.cronRules {
		logs.Info.Println("cron " + r.cronRules[rule].Name + ": start...")
		go r.ruleWorker(ctx, r.cronRules[rule])
	}

	// result pipeline
	go r.resultWorker(ctx)
}

func (r *Ruleset) ruleWorker(ctx extraTypes.Context, rule Rule) {
	p, err := cron.ParseUTC(rule.When)
	if err != nil {
		logs.Err.Println(err)
		return
	}
	nextTime, err := p.Next(time.Now())
	if err != nil {
		logs.Err.Println(err)
		return
	}
	for {
		if nextTime.Format("2006-01-02 15:04") == time.Now().Format("2006-01-02 15:04") {
			logs.Info.Println("cron " + rule.Name + ": scheduled")
			msgs := func() []extraTypes.MsgPayload {
				defer func() {
					if rc := recover(); rc != nil {
						logs.Warn.Println("ruleWorker recover " + rule.Name)
						if v, ok := rc.(error); ok {
							logs.Err.Println(v)
						}
					}
				}()

				items, err := store.Chatbot.OAuthGetAvailable(r.Type)
				if err != nil {
					logs.Err.Println(err)
					return nil
				}
				if len(items) > 0 {
					var result []extraTypes.MsgPayload
					for _, oauth := range items {
						ra := rule.Action(extraTypes.Context{
							Original: oauth.Topic,
							AsUser:   types.ParseUserId(oauth.Uid),
							Token:    oauth.Token,
						})
						result = append(result, ra...)
					}
					return result
				}

				return rule.Action(ctx)
			}()
			if len(msgs) > 0 {
				for _, item := range msgs {
					r.outCh <- item
				}
			}
		}
		nextTime, err = p.Next(time.Now())
		if err != nil {
			logs.Err.Println(err)
			continue
		}
		time.Sleep(2 * time.Second)
	}
}

func (r *Ruleset) resultWorker(ctx extraTypes.Context) {
	for out := range r.outCh {
		// pipeline
		r.pipeline(ctx, out)
	}
}

func (r *Ruleset) pipeline(_ extraTypes.Context, res extraTypes.MsgPayload) {
	if res == nil {
		return
	}
	// todo send message
}
