package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tinode/chat/server/extra/ruleset/action"
	"github.com/tinode/chat/server/extra/ruleset/session"
	"github.com/tinode/chat/server/extra/ruleset/workflow"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/adjust/rmq/v5"
	"github.com/redis/go-redis/v9"
	"github.com/tinode/chat/server/auth"
	"github.com/tinode/chat/server/extra/bots"
	botGithub "github.com/tinode/chat/server/extra/bots/github"
	botPocket "github.com/tinode/chat/server/extra/bots/pocket"
	"github.com/tinode/chat/server/extra/pkg/cache"
	extraStore "github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/store/model"
	extraTypes "github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/extra/utils"
	"github.com/tinode/chat/server/extra/vendors"
	"github.com/tinode/chat/server/extra/vendors/dropbox"
	"github.com/tinode/chat/server/extra/vendors/github"
	"github.com/tinode/chat/server/extra/vendors/pocket"
	"github.com/tinode/chat/server/logs"
	"github.com/tinode/chat/server/store"
	"github.com/tinode/chat/server/store/types"
)

const BotFather = "BotFather"

func isBot(subs types.Subscription) bool {
	// normal bot user
	if subs.GetState() != types.StateOK {
		return false
	}
	// verified
	trusted := subs.GetTrusted()
	if trusted == nil {
		return false
	}
	if !isVerified(trusted) {
		return false
	}
	// check name
	public := subs.GetPublic()
	if public == nil {
		return false
	}
	name := utils.Fn(public)
	if !strings.HasSuffix(name, bots.BotNameSuffix) {
		return false
	}

	return true
}

func isBotUser(user *types.User) bool {
	if user == nil {
		return false
	}
	// normal bot user
	if user.State != types.StateOK {
		return false
	}
	// verified
	if !isVerified(user.Trusted) {
		return false
	}
	// check name
	name := utils.Fn(user.Public)
	if !strings.HasSuffix(name, bots.BotNameSuffix) {
		return false
	}

	return true
}

func isVerified(trusted interface{}) bool {
	if v, ok := trusted.(map[string]interface{}); ok {
		if b, ok := v["verified"]; ok {
			if vv, ok := b.(bool); ok {
				return vv
			}
		}
	}
	return false
}

func botName(subs types.Subscription) string {
	public := subs.GetPublic()
	if public == nil {
		return ""
	}
	name := utils.Fn(public)
	name = strings.ReplaceAll(name, bots.BotNameSuffix, "")
	return name
}

// botSend bot send message, rcptTo: user uid: bot
func botSend(rcptTo string, uid types.Uid, out extraTypes.MsgPayload, option ...interface{}) {
	if out == nil {
		return
	}

	t := globals.hub.topicGet(rcptTo)
	if t == nil {
		var original = ""
		switch types.GetTopicCat(rcptTo) {
		case types.TopicCatP2P:
			u1, u2, err := types.ParseP2P(rcptTo)
			if err != nil {
				logs.Err.Println(err)
				return
			}
			if u1 == uid {
				original = u2.UserId()
			} else {
				original = u1.UserId()
			}
		default:
			original = uid.UserId() // initTopicP2P: userID2 := types.ParseUserId(t.xoriginal)
		}

		sess := &Session{
			uid:     uid,
			authLvl: auth.LevelAuth,
			subs:    make(map[string]*Subscription),
			send:    make(chan interface{}, sendQueueLimit+32),
			stop:    make(chan interface{}, 1),
			detach:  make(chan string, 64),
		}
		msg := &ClientComMessage{
			Sub: &MsgClientSub{
				Topic:   uid.UserId(),
				Get:     &MsgGetQuery{},
				Created: false,
				Newsub:  false,
			},
			Original:  original,
			RcptTo:    rcptTo,
			AsUser:    uid.UserId(),
			AuthLvl:   int(auth.LevelAuth),
			MetaWhat:  0,
			Timestamp: time.Now(),
			sess:      sess,
			init:      true,
		}
		globals.hub.join <- msg
		// wait sometime
		time.Sleep(200 * time.Millisecond)

		t = globals.hub.topicGet(rcptTo)
	}

	if t == nil {
		logs.Err.Printf("topic %s error, Failed to send", rcptTo)
		return
	}

	heads, contents := extraTypes.Convert([]extraTypes.MsgPayload{out})
	if !(len(heads) > 0 && len(contents) > 0) {
		logs.Err.Printf("topic %s convert error, Failed to send", rcptTo)
		return
	}
	head, content := heads[0], contents[0]

	// set head context
	if len(option) > 0 {
		for _, item := range option {
			switch v := item.(type) {
			case extraTypes.Context:
				if head != nil {
					if v.WorkflowFlag != "" {
						head["x-workflow-flag"] = v.WorkflowFlag
					}
					if v.WorkflowVersion > 0 {
						head["x-workflow-version"] = v.WorkflowVersion
					}
				}
			}
		}
	}

	msg := &ClientComMessage{
		Pub: &MsgClientPub{
			Topic:   rcptTo,
			Head:    head,
			Content: content,
		},
		AsUser:    uid.UserId(),
		Timestamp: types.TimeNow(),
	}
	if strings.HasPrefix(rcptTo, "grp") {
		msg.Original = rcptTo
		msg.RcptTo = rcptTo
	}
	t.handleClientMsg(msg)
}

func newProvider(category string) vendors.OAuthProvider {
	var provider vendors.OAuthProvider

	switch category {
	case pocket.ID:
		provider = pocket.NewPocket(botPocket.Config.ConsumerKey, "", "", "")
	case github.ID:
		provider = github.NewGithub(botGithub.Config.ID, botGithub.Config.Secret, "", "")
	case dropbox.ID:
		provider = dropbox.NewDropbox("", "", "", "")
	default:
		return nil
	}

	return provider
}

func botIncomingMessage(t *Topic, msg *ClientComMessage) {
	// check topic owner user
	if msg.AsUser == msg.Pub.Topic {
		return
	}
	if msg.Original == "" || msg.RcptTo == "" {
		return
	}

	subs, err := store.Topics.GetUsers(msg.Pub.Topic, nil)
	if err != nil {
		logs.Err.Println("hook bot incoming", err)
		return
	}

	uid := types.ParseUserId(msg.AsUser)
	ctx := extraTypes.Context{
		Id:        msg.Id,
		Original:  msg.Original,
		RcptTo:    msg.RcptTo,
		AsUser:    uid,
		AuthLvl:   msg.AuthLvl,
		MetaWhat:  msg.MetaWhat,
		Timestamp: msg.Timestamp,
	}

	// behavior
	bots.Behavior(uid, bots.MessageBotIncomingBehavior, 1)

	// user auth record
	_, authLvl, _, _, _ := store.Users.GetAuthRecord(uid, "basic")

	// bot
	for _, sub := range subs {
		if !isBot(sub) {
			continue
		}

		// bot name
		name := botName(sub)
		handle, ok := bots.List()[name]
		if !ok {
			continue
		}

		if !handle.IsReady() {
			logs.Info.Printf("bot %s unavailable", t.name)
			continue
		}

		var payload extraTypes.MsgPayload

		switch handle.AuthLevel() {
		case auth.LevelRoot:
			if authLvl != auth.LevelRoot {
				payload = extraTypes.TextMsg{Text: "Unauthorized"}
			}
		}

		// auth
		if payload == nil {
			// session
			if sess, ok := sessionCurrent(uid, msg.Original); ok && sess.State == model.SessionStart {
				// session cancel command
				isCancel := false
				if msg.Pub.Head == nil {
					if v, ok := msg.Pub.Content.(string); ok {
						if v == "cancel" {
							_ = extraStore.Chatbot.SessionState(ctx.AsUser, ctx.Original, model.SessionCancel)
							payload = extraTypes.TextMsg{Text: "session cancel"}
							isCancel = true
						}
					}
				}
				if !isCancel {
					ctx.SessionRuleId = sess.RuleID
					ctx.SessionInitValues = extraTypes.KV(sess.Init)
					ctx.SessionLastValues = extraTypes.KV(sess.Values)

					// get action handler
					var botHandler bots.Handler
					for _, handler := range bots.List() {
						for _, item := range handler.Rules() {
							switch v := item.(type) {
							case []session.Rule:
								for _, rule := range v {
									if rule.Id == sess.RuleID {
										botHandler = handler
									}
								}
							}
						}
					}
					if botHandler == nil {
						payload = extraTypes.TextMsg{Text: "error session"}
					} else {
						payload, err = botHandler.Session(ctx, msg.Pub.Content)
						if err != nil {
							logs.Warn.Printf("topic[%s]: failed to run bot: %v", t.name, err)
						}
					}
				}
			}
			// action
			if payload == nil {
				if msg.Pub.Head != nil {
					var cm extraTypes.ChatMessage
					d, err := json.Marshal(msg.Pub.Content)
					if err != nil {
						logs.Err.Println(err)
					}
					err = json.Unmarshal(d, &cm)
					if err != nil {
						logs.Err.Println(err)
					}
					var seq float64
					var option string
					for _, ent := range cm.Ent {
						if ent.Tp == "EX" {
							if m, ok := ent.Data.Val.(map[string]interface{}); ok {
								if v, ok := m["seq"]; ok {
									seq = v.(float64)
								}
								if v, ok := m["resp"]; ok {
									values := v.(map[string]interface{})
									for s := range values {
										option = s
									}
								}
							}
						}
					}
					if seq > 0 {
						message, err := extraStore.Chatbot.GetMessage(msg.RcptTo, int(seq))
						if err != nil {
							logs.Err.Println(err)
						}
						actionRuleId := ""
						if src, ok := extraTypes.KV(message.Content).Map("src"); ok {
							if id, ok := src["id"]; ok {
								actionRuleId = id.(string)
							}
						}
						ctx.SeqId = int(seq)
						ctx.ActionRuleId = actionRuleId

						// get action handler
						var botHandler bots.Handler
						for _, handler := range bots.List() {
							for _, item := range handler.Rules() {
								switch v := item.(type) {
								case []action.Rule:
									for _, rule := range v {
										if rule.Id == actionRuleId {
											botHandler = handler
										}
									}
								}
							}
						}
						if botHandler == nil {
							payload = extraTypes.TextMsg{Text: "error action"}
						} else {
							payload, err = botHandler.Action(ctx, option)
							if err != nil {
								logs.Warn.Printf("topic[%s]: failed to run bot: %v", t.name, err)
							}

							if payload != nil {
								botUid := types.ParseUserId(msg.Original)
								botSend(msg.RcptTo, botUid, payload, extraTypes.WithContext(ctx))

								// workflow action step
								workflowFlag, _ := extraTypes.KV(message.Head).String("x-workflow-flag")
								workflowVersion, _ := extraTypes.KV(message.Head).Int64("x-workflow-version")
								nextWorkflow(ctx, workflowFlag, int(workflowVersion), msg.RcptTo, botUid)
								return
							}
						}
					}
				}
			}
			// command
			if payload == nil {
				var content interface{}
				if msg.Pub.Head == nil {
					content = msg.Pub.Content
				} else {
					// Compatible with drafty
					if m, ok := msg.Pub.Content.(map[string]interface{}); ok {
						if txt, ok := m["txt"]; ok {
							content = txt
						}
					}
				}
				// check "/" prefix
				if in, ok := content.(string); ok && strings.HasPrefix(in, "/") {
					in = strings.Replace(in, "/", "", 1)
					payload, err = handle.Command(ctx, in)
					if err != nil {
						logs.Warn.Printf("topic[%s]: failed to run bot: %v", t.name, err)
					}

					// stats
					statsInc("BotRunCommandTotal", 1)

					// error message
					if payload == nil {
						payload = extraTypes.TextMsg{Text: "error command"}
					}
				}
			}
			// workflow command trigger
			if payload == nil {
				var content interface{}
				if msg.Pub.Head == nil {
					content = msg.Pub.Content
				} else {
					// Compatible with drafty
					if m, ok := msg.Pub.Content.(map[string]interface{}); ok {
						if txt, ok := m["txt"]; ok {
							content = txt
						}
					}
				}
				// check "~" prefix
				if in, ok := content.(string); ok && strings.HasPrefix(in, "~") {
					var workflowFlag string
					var workflowVersion int
					in = strings.Replace(in, "~", "", 1)
					payload, workflowFlag, workflowVersion, err = handle.Workflow(ctx, msg.Pub.Head, in, extraTypes.WorkflowCommandTriggerOperate)
					if err != nil {
						logs.Warn.Printf("topic[%s]: failed to run bot: %v", t.name, err)
					}
					ctx.WorkflowFlag = workflowFlag
					ctx.WorkflowVersion = workflowVersion

					// stats
					statsInc("BotTriggerWorkflowTotal", 1)

					// error message
					if payload == nil {
						payload = extraTypes.TextMsg{Text: "error workflow"}
					}
				}
			}
			// condition
			if payload == nil {
				if msg.Pub.Head != nil {
					fUid := ""
					fSeq := int64(0)
					if v, ok := msg.Pub.Head["forwarded"]; ok {
						if s, ok := v.(string); ok {
							f := strings.Split(s, ":")
							if len(f) == 2 {
								fUid = f[0]
								fSeq, _ = strconv.ParseInt(f[1], 10, 64)
							}
						}
					}

					if fUid != "" && fSeq > 0 {
						uid2 := types.ParseUserId(fUid)
						topic := uid.P2PName(uid2)
						message, err := extraStore.Chatbot.GetMessage(topic, int(fSeq))
						if err != nil {
							logs.Err.Println(err)
						}

						if message.ID > 0 {
							src, _ := extraTypes.KV(message.Content).Map("src")
							tye, _ := extraTypes.KV(message.Content).String("tye")
							d, _ := json.Marshal(src)
							pl := extraTypes.ToPayload(tye, d)
							ctx.Condition = tye
							payload, err = handle.Condition(ctx, pl)
							if err != nil {
								logs.Warn.Printf("topic[%s]: failed to run bot: %v", t.name, err)
							}

							// stats
							statsInc("BotRunConditionTotal", 1)
						}
					}
				}
			}
			// input
			if payload == nil {
				payload, err = handle.Input(ctx, msg.Pub.Head, msg.Pub.Content)
				if err != nil {
					logs.Warn.Printf("topic[%s]: failed to run bot: %v", t.name, err)
					continue
				}

				// stats
				statsInc("BotRunInputTotal", 1)
			}
		}

		// send message
		if payload == nil {
			continue
		}

		botUid := types.ParseUserId(msg.Original)
		botSend(msg.RcptTo, botUid, payload, extraTypes.WithContext(ctx))
	}
}

func groupIncomingMessage(t *Topic, msg *ClientComMessage, event extraTypes.GroupEvent) {
	subs, err := store.Topics.GetUsers(msg.Original, nil)
	if err != nil {
		logs.Err.Println("hook bot incoming", err)
		return
	}
	// check bot user incoming
	for _, sub := range subs {
		if !isBot(sub) {
			continue
		}
		if strings.TrimPrefix(msg.AsUser, "usr") == sub.User {
			return
		}
	}

	uid := types.ParseUserId(msg.AsUser)
	ctx := extraTypes.Context{
		Id:        msg.Id,
		Original:  msg.Original,
		RcptTo:    msg.RcptTo,
		AsUser:    uid,
		AuthLvl:   msg.AuthLvl,
		MetaWhat:  msg.MetaWhat,
		Timestamp: msg.Timestamp,
	}

	// behavior
	bots.Behavior(uid, bots.MessageGroupIncomingBehavior, 1)

	// user auth record
	_, authLvl, _, _, _ := store.Users.GetAuthRecord(uid, "basic")

	// bot
	for _, sub := range subs {
		if !isBot(sub) {
			continue
		}

		// bot name
		name := botName(sub)
		handle, ok := bots.List()[name]
		if !ok {
			continue
		}

		if !handle.IsReady() {
			logs.Info.Printf("bot %s unavailable", t.name)
			continue
		}

		var payload extraTypes.MsgPayload

		switch handle.AuthLevel() {
		case auth.LevelRoot:
			if authLvl != auth.LevelRoot {
				payload = extraTypes.TextMsg{Text: "Unauthorized"}
			}
		}

		// auth
		if payload == nil {
			// condition
			if msg.Pub != nil && msg.Pub.Head != nil {
				fUid := ""
				fSeq := int64(0)
				if v, ok := msg.Pub.Head["forwarded"]; ok {
					if s, ok := v.(string); ok {
						f := strings.Split(s, ":")
						if len(f) == 2 {
							fUid = f[0]
							fSeq, _ = strconv.ParseInt(f[1], 10, 64)
						}
					}
				}

				if fUid != "" && fSeq > 0 {
					uid2 := types.ParseUserId(fUid)
					topic := uid.P2PName(uid2)
					message, err := extraStore.Chatbot.GetMessage(topic, int(fSeq))
					if err != nil {
						logs.Err.Println(err)
					}

					if message.ID > 0 {
						src, _ := extraTypes.KV(message.Content).Map("src")
						tye, _ := extraTypes.KV(message.Content).String("tye")
						d, _ := json.Marshal(src)
						pl := extraTypes.ToPayload(tye, d)
						ctx.Condition = tye
						payload, err = handle.Condition(ctx, pl)
						if err != nil {
							logs.Warn.Printf("topic[%s]: failed to run bot: %v", t.name, err)
						}

						// stats
						statsInc("BotRunConditionTotal", 1)
					}
				}
			}
		}

		// group
		if payload == nil {
			ctx.GroupEvent = event
			var head map[string]any
			var content any
			if msg.Pub != nil {
				head = msg.Pub.Head
				content = msg.Pub.Content
			}
			payload, err = handle.Group(ctx, head, content)
			if err != nil {
				logs.Warn.Printf("topic[%s]: failed to run group bot: %v", t.name, err)
				continue
			}

			// stats
			statsInc("BotRunGroupTotal", 1)
		}

		// send message
		if payload == nil {
			continue
		}

		botUid := types.ParseUid(sub.User)
		botSend(msg.RcptTo, botUid, payload)
	}
}

func nextWorkflow(ctx extraTypes.Context, workflowFlag string, workflowVersion int, rcptTo string, botUid types.Uid) {
	if workflowFlag != "" && workflowVersion > 0 {
		workflowData, err := extraStore.Chatbot.WorkflowGet(ctx.AsUser, ctx.Original, workflowFlag)
		if err != nil {
			logs.Err.Println(err)
			return
		}
		for _, handler := range bots.List() {
			for _, item := range handler.Rules() {
				switch v := item.(type) {
				case []workflow.Rule:
					for _, rule := range v {
						if rule.Id == workflowData.RuleID {
							ctx.WorkflowFlag = workflowFlag
							ctx.WorkflowVersion = workflowVersion
							ctx.WorkflowRuleId = workflowData.RuleID
							ctx.WorkflowStepIndex = int(workflowData.Step)
							payload, _, _, err := handler.Workflow(ctx, nil, nil, extraTypes.WorkflowNextOperate)
							if err != nil {
								logs.Err.Println(err)
								return
							}
							botSend(rcptTo, botUid, payload, extraTypes.WithContext(ctx))
						}
					}
				}
			}
		}
	}
}

func notifyAfterReboot() {
	botUid, _, _, _, err := store.Users.GetAuthUniqueRecord("basic", fmt.Sprintf("server%s", bots.BotNameSuffix))
	if err != nil {
		logs.Err.Println(err)
		return
	}

	creds, err := extraStore.Chatbot.GetCredentials()
	if err != nil {
		logs.Err.Println(err)
		return
	}

	for _, cred := range creds {
		_, level, _, _, err := store.Users.GetAuthRecord(store.EncodeUid(cred.Userid), "basic")
		if err != nil {
			logs.Err.Println(err)
			continue
		}
		if level != auth.LevelRoot {
			continue
		}
		rcptTo := store.EncodeUid(cred.Userid).P2PName(botUid)
		if rcptTo != "" {
			botSend(rcptTo, botUid, extraTypes.TextMsg{Text: "reboot"})
		}
	}
}

func onlineStatus(usrStr string) {
	uid := types.ParseUserId(usrStr)
	user, err := store.Users.Get(uid)
	if err != nil {
		return
	}
	if isBotUser(user) {
		return
	}

	ctx := context.Background()
	key := fmt.Sprintf("online:%s", usrStr)
	_, err = cache.DB.Get(ctx, key).Result()
	if err == redis.Nil {
		cache.DB.Set(ctx, key, time.Now().Unix(), 30*time.Minute)
	} else if err != nil {
		return
	} else {
		cache.DB.Expire(ctx, key, 30*time.Minute)
	}
}

func sessionCurrent(uid types.Uid, topic string) (model.Session, bool) {
	sess, err := extraStore.Chatbot.SessionGet(uid, topic)
	if err != nil {
		return model.Session{}, false
	}
	return sess, true
}

func errorResponse(rw http.ResponseWriter, text string) {
	rw.WriteHeader(http.StatusBadRequest)
	_, _ = rw.Write([]byte(text))
}

type AsyncMessageConsumer struct {
	name string
}

func NewAsyncMessageConsumer() *AsyncMessageConsumer {
	return &AsyncMessageConsumer{name: "consumer"}
}

func (c *AsyncMessageConsumer) Consume(delivery rmq.Delivery) {
	payload := delivery.Payload()

	var qp extraTypes.QueuePayload
	err := json.Unmarshal([]byte(payload), &qp)
	if err != nil {
		if err := delivery.Reject(); err != nil {
			logs.Err.Printf("failed to reject %s: %s\n", payload, err)
			return
		}
		return
	}

	uid := types.ParseUserId(qp.Uid)
	msg := extraTypes.ToPayload(qp.Type, qp.Msg)
	botSend(qp.RcptTo, uid, msg)

	if err := delivery.Ack(); err != nil {
		logs.Err.Printf("failed to ack %s: %s\n", payload, err)
		return
	}
}
