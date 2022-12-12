package main

import (
	"encoding/json"
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/cache"
	"github.com/tinode/chat/server/extra/channels"
	extraStore "github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/logs"
	"net/http"
	"strings"

	// bots
	_ "github.com/tinode/chat/server/extra/bots/anki"
	_ "github.com/tinode/chat/server/extra/bots/bark"
	_ "github.com/tinode/chat/server/extra/bots/cloudflare"
	_ "github.com/tinode/chat/server/extra/bots/finance"
	_ "github.com/tinode/chat/server/extra/bots/genshin"
	_ "github.com/tinode/chat/server/extra/bots/github"
	_ "github.com/tinode/chat/server/extra/bots/help"
	_ "github.com/tinode/chat/server/extra/bots/iot"
	_ "github.com/tinode/chat/server/extra/bots/mtg"
	_ "github.com/tinode/chat/server/extra/bots/notion"
	_ "github.com/tinode/chat/server/extra/bots/okr"
	_ "github.com/tinode/chat/server/extra/bots/pocket"
	_ "github.com/tinode/chat/server/extra/bots/qr"
	_ "github.com/tinode/chat/server/extra/bots/search"
	_ "github.com/tinode/chat/server/extra/bots/server"
	_ "github.com/tinode/chat/server/extra/bots/share"
	_ "github.com/tinode/chat/server/extra/bots/subscribe"
	_ "github.com/tinode/chat/server/extra/bots/url"
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
	mux.Handle("/u/", newUrlRouter())
}

func hookStore() {
	// init cache
	cache.InitCache()
	// open database
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
	statsRegisterInt("BotRunInputTotal")
	statsRegisterInt("BotRunGroupTotal")
	statsRegisterInt("BotRunAgentTotal")
	statsRegisterInt("BotRunCommandTotal")
	statsRegisterInt("BotRunConditionTotal")
	statsRegisterInt("BotRunCronTotal")
	statsRegisterInt("BotRunFormTotal")

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
	if strings.HasPrefix(msg.Pub.Topic, "grp") {
		groupIncomingMessage(t, msg)
	} else {
		botIncomingMessage(t, msg)
	}
}

func hookMounted() {
	// notify after reboot
	go notifyAfterReboot()
}
