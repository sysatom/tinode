package main

import (
	"encoding/json"
	"github.com/tinode/chat/server/auth"
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/channels"
	extraStore "github.com/tinode/chat/server/extra/store"
	extraTypes "github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/logs"
	"github.com/tinode/chat/server/store"
	"github.com/tinode/chat/server/store/types"
	"net/http"
	"strconv"
	"strings"

	// bots
	_ "github.com/tinode/chat/server/extra/bots/bark"
	_ "github.com/tinode/chat/server/extra/bots/cloudflare"
	_ "github.com/tinode/chat/server/extra/bots/finance"
	_ "github.com/tinode/chat/server/extra/bots/github"
	_ "github.com/tinode/chat/server/extra/bots/help"
	_ "github.com/tinode/chat/server/extra/bots/notion"
	_ "github.com/tinode/chat/server/extra/bots/okr"
	_ "github.com/tinode/chat/server/extra/bots/pocket"
	_ "github.com/tinode/chat/server/extra/bots/server"
	_ "github.com/tinode/chat/server/extra/bots/subscribe"
	_ "github.com/tinode/chat/server/extra/bots/webhook"

	// push
	_ "github.com/tinode/chat/server/extra/bark"

	// store
	_ "github.com/tinode/chat/server/extra/store/mysql"

	// cache
	_ "github.com/tinode/chat/server/extra/cache"
)

// hook

func hookMux(mux *http.ServeMux) {
	mux.Handle("/extra/", newRouter())
}

func hookStore() {
	err := extraStore.Store.Open()
	if err != nil {
		panic(err)
	}
}

func hookBot(jsconfig json.RawMessage) {
	// init bots
	err := bots.Init(jsconfig)
	if err != nil {
		logs.Err.Fatal("Failed to initialize bot:", err)
	}

	// bot father
	err = initializeBotFather()
	if err != nil {
		logs.Err.Fatal("Failed to create or update bot father:", err)
	}

	// bot users
	err = initializeBotUsers()
	if err != nil {
		logs.Err.Fatal("Failed to create or update bot users:", err)
	}

	// bot cron
	err = bots.Cron(botSend)
	if err != nil {
		logs.Err.Fatal("Failed to bot cron:", err)
	}

	// stats register
	statsRegisterInt("BotTotal")
	statsRegisterInt("BotRunTotal")

	statsSet("BotTotal", int64(len(bots.List())))
}

func hookChannel(jsconfig json.RawMessage) {
	err := channels.Init(jsconfig)
	if err != nil {
		logs.Err.Fatal("Failed to initialize channel:", err)
	}

	err = initializeChannels()
	if err != nil {
		logs.Err.Fatal("Failed to create or update channels:", err)
	}

	err = initializeCrawler()
	if err != nil {
		logs.Err.Fatal("Failed to initialize crawler:", err)
	}

	// stats register
	statsRegisterInt("ChannelTotal")
	statsRegisterInt("ChannelPublishTotal")

	statsSet("ChannelTotal", int64(len(channels.List())))
}

func hookHandleBotIncomingMessage(t *Topic, msg *ClientComMessage) {
	// check topic owner user
	_, u2, _ := types.ParseP2P(msg.Pub.Topic)
	if !u2.IsZero() && u2.Compare(types.ParseUserId(msg.AsUser)) == 0 {
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
			// command
			if msg.Pub.Head == nil {
				payload, err = handle.Command(ctx, msg.Pub.Content)
				if err != nil {
					logs.Warn.Printf("topic[%s]: failed to run bot: %v", t.name, err)
				}

				// stats
				statsInc("BotRunTotal", 1)
			}

			if payload == nil {
				// condition
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
							src, _ := message.Content.Map("src")
							tye, _ := message.Content.String("tye")
							d, _ := json.Marshal(src)
							pl := extraTypes.ToPayload(tye, d)
							ctx.Condition = tye
							payload, err = handle.Condition(ctx, pl)
							if err != nil {
								logs.Warn.Printf("topic[%s]: failed to run bot: %v", t.name, err)
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
				}
			}
		}

		// send  message
		if payload == nil {
			continue
		}

		uid2 := types.ParseUserId(msg.Original)
		botSend(uid, uid2, payload)
	}
}
