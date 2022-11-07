package main

import (
	"encoding/json"
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/channels"
	"github.com/tinode/chat/server/extra/router"
	extraStore "github.com/tinode/chat/server/extra/store"
	extraTypes "github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/logs"
	"github.com/tinode/chat/server/store"
	"github.com/tinode/chat/server/store/types"
	"net/http"
	// bots
	_ "github.com/tinode/chat/server/extra/bots/bark"
	_ "github.com/tinode/chat/server/extra/bots/help"
	_ "github.com/tinode/chat/server/extra/bots/subscribe"

	// push
	_ "github.com/tinode/chat/server/extra/bark"

	// store
	_ "github.com/tinode/chat/server/extra/store/mysql"
)

// hot

func hookMux(mux *http.ServeMux) {
	mux.Handle("/extra/", http.HandlerFunc(router.ServeExtra))
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
	subs, err := store.Topics.GetUsers(msg.Pub.Topic, nil)
	if err != nil {
		logs.Err.Println(err)
		return
	}

	ctx := extraTypes.Context{
		Id:        msg.Id,
		Original:  msg.Original,
		RcptTo:    msg.RcptTo,
		AsUser:    types.ParseUserId(msg.AsUser),
		AuthLvl:   msg.AuthLvl,
		MetaWhat:  msg.MetaWhat,
		Timestamp: msg.Timestamp,
	}

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
		heads, contents, err := handle.Run(ctx, msg.Pub.Head, msg.Pub.Content)
		if err != nil {
			logs.Warn.Printf("topic[%s]: failed to run bot: %v", t.name, err)
			continue
		}
		// multiple messages
		for i, content := range contents {
			head := heads[i]
			if content == nil {
				continue
			}

			// stats
			statsInc("BotRunTotal", 1)

			now := types.TimeNow()
			if err := store.Messages.Save(
				&types.Message{
					ObjHeader: types.ObjHeader{CreatedAt: now},
					SeqId:     t.lastID + 1,
					Topic:     t.name,
					From:      sub.User,
					Head:      head,
					Content:   content,
				}, nil, true); err != nil {
				logs.Warn.Printf("topic[%s]: failed to save bot message: %v", t.name, err)
				continue
			}

			t.lastID++
			t.touched = now

			data := &ServerComMessage{
				Data: &MsgServerData{
					Topic:     msg.Original,
					From:      sub.User,
					Timestamp: now,
					SeqId:     t.lastID,
					Head:      head,
					Content:   content,
				},
				// Internal-only values.
				Id:        msg.Id,
				RcptTo:    msg.RcptTo,
				AsUser:    sub.User,
				Timestamp: now,
				sess:      msg.sess,
			}

			t.broadcastToSessions(data)

			asUid := types.ParseUid(sub.User)

			// sendPush will update unread message count and send push notification.
			if pushRcpt := t.pushForData(asUid, data.Data); pushRcpt != nil {
				sendPush(pushRcpt)
			}
		}
	}
}
