package main

import (
	"encoding/json"
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/pkg/cache"
	"github.com/tinode/chat/server/extra/pkg/channels"
	"github.com/tinode/chat/server/extra/pkg/queue"
	"github.com/tinode/chat/server/extra/pkg/route"
	extraStore "github.com/tinode/chat/server/extra/store"
	extraMysql "github.com/tinode/chat/server/extra/store/mysql"
	extraTypes "github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/extra/vendors"
	"github.com/tinode/chat/server/logs"
	"github.com/tinode/chat/server/store"
	"github.com/tinode/chat/server/store/types"
	"net/http"
	"strings"

	// bots
	_ "github.com/tinode/chat/server/extra/bots/anki"
	_ "github.com/tinode/chat/server/extra/bots/attendance"
	_ "github.com/tinode/chat/server/extra/bots/bark"
	_ "github.com/tinode/chat/server/extra/bots/clipboard"
	_ "github.com/tinode/chat/server/extra/bots/cloudflare"
	_ "github.com/tinode/chat/server/extra/bots/dev"
	_ "github.com/tinode/chat/server/extra/bots/download"
	_ "github.com/tinode/chat/server/extra/bots/finance"
	_ "github.com/tinode/chat/server/extra/bots/genshin"
	_ "github.com/tinode/chat/server/extra/bots/github"
	_ "github.com/tinode/chat/server/extra/bots/gpt"
	_ "github.com/tinode/chat/server/extra/bots/iot"
	_ "github.com/tinode/chat/server/extra/bots/leetcode"
	_ "github.com/tinode/chat/server/extra/bots/linkit"
	_ "github.com/tinode/chat/server/extra/bots/markdown"
	_ "github.com/tinode/chat/server/extra/bots/mtg"
	_ "github.com/tinode/chat/server/extra/bots/notion"
	_ "github.com/tinode/chat/server/extra/bots/obsidian"
	_ "github.com/tinode/chat/server/extra/bots/okr"
	_ "github.com/tinode/chat/server/extra/bots/pocket"
	_ "github.com/tinode/chat/server/extra/bots/qr"
	_ "github.com/tinode/chat/server/extra/bots/queue"
	_ "github.com/tinode/chat/server/extra/bots/rust"
	_ "github.com/tinode/chat/server/extra/bots/search"
	_ "github.com/tinode/chat/server/extra/bots/server"
	_ "github.com/tinode/chat/server/extra/bots/share"
	_ "github.com/tinode/chat/server/extra/bots/subscribe"
	_ "github.com/tinode/chat/server/extra/bots/url"
	_ "github.com/tinode/chat/server/extra/bots/web"
	_ "github.com/tinode/chat/server/extra/bots/webhook"
	_ "github.com/tinode/chat/server/extra/bots/workflow"

	// push
	_ "github.com/tinode/chat/server/extra/pkg/bark"

	// cache
	_ "github.com/tinode/chat/server/extra/pkg/cache"
)

// hook

func hookMux() *http.ServeMux {
	// Webservice
	wc := route.NewContainer()
	for _, bot := range bots.List() {
		if ws := bot.Webservice(); ws != nil {
			wc.Add(ws)
		}
	}
	route.AddSwagger(wc)
	mux := wc.ServeMux

	mux.Handle("/extra/", newRouter())
	mux.Handle("/app/", newWebappRouter())
	mux.Handle("/u/", newUrlRouter())
	mux.Handle("/d/", newDownloadRouter())

	return mux
}

func hookStore() {
	// init cache
	cache.InitCache()
	// init database
	extraMysql.Init()
	extraStore.Init()
	err := extraStore.Store.Open()
	if err != nil {
		panic(err)
	}
}

func hookBot(jsconfig json.RawMessage, vc json.RawMessage) {
	// set vendors configs
	vendors.Configs = vc

	// init bots
	err := bots.Init(jsconfig)
	if err != nil {
		logs.Err.Fatal("Failed to initialize bot:", err)
	}

	// bootstrap bots
	err = bots.Bootstrap()
	if err != nil {
		logs.Err.Fatal("Failed to bootstrap bot:", err)
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
	globals.cronRuleset, err = bots.Cron(botSend)
	if err != nil {
		logs.Err.Fatal("Failed to bot cron:", err)
	}

	// stats register
	statsRegisterInt("BotTotal")
	statsRegisterInt("BotRunInputTotal")
	statsRegisterInt("BotRunGroupTotal")
	statsRegisterInt("BotRunAgentTotal")
	statsRegisterInt("BotRunCommandTotal")
	statsRegisterInt("BotRunConditionTotal")
	statsRegisterInt("BotRunCronTotal")
	statsRegisterInt("BotRunFormTotal")
	statsRegisterInt("BotTriggerPipelineTotal")

	statsSet("BotTotal", int64(len(bots.List())))
}

func hookChannel() {
	err := channels.Init()
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

func hookHandleIncomingMessage(t *Topic, msg *ClientComMessage) {
	// update online status
	onlineStatus(msg.AsUser)
	// check grp or p2p
	if strings.HasPrefix(msg.Original, "grp") {
		groupIncomingMessage(t, msg, extraTypes.GroupEventReceive)
	} else {
		botIncomingMessage(t, msg)
	}
}

func hookHandleGroupEvent(t *Topic, msg *ClientComMessage, event int) {
	if strings.HasPrefix(msg.Original, "grp") {
		switch extraTypes.GroupEvent(event) {
		case extraTypes.GroupEventJoin:
			msg.AsUser = msg.Set.MsgSetQuery.Sub.User
		case extraTypes.GroupEventExit:
			msg.AsUser = msg.Del.User
		}
		user, err := store.Users.Get(types.ParseUserId(msg.AsUser))
		if err != nil {
			logs.Err.Println(err)
		}
		// Current user is bot
		if isBotUser(user) {
			return
		}
		groupIncomingMessage(t, msg, extraTypes.GroupEvent(event))
	}
}

func hookMounted() {
	// notify after reboot
	go notifyAfterReboot()
}

func hookQueue() {
	queue.InitMessageQueue(NewAsyncMessageConsumer())
}

func hookEvent() {
	onSendEvent()
	onPushInstruct()
}
