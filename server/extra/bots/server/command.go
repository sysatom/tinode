package server

import (
	"fmt"
	"github.com/tinode/chat/server/extra"
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/ruleset/command"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/types"
	"runtime"
)

var commandRules = []command.Rule{
	{
		Define: "version",
		Help:   `Version`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return types.TextMsg{Text: fmt.Sprintf("Chatbot framework v%s", extra.Version)}
		},
	},
	{
		Define: "vars",
		Help:   `vars url`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return types.LinkMsg{
				Title: "Vars Url",
				Url:   fmt.Sprintf("%s/debug/vars", types.AppUrl()),
			}
		},
	},
	{
		Define: "mem stats",
		Help:   `App memory stats`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)

			return types.InfoMsg{
				Title: "Memory stats",
				Model: memStats,
			}
		},
	},
	{
		Define: "golang stats",
		Help:   `App golang stats`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			numGoroutine := runtime.NumGoroutine()

			return types.InfoMsg{
				Title: "Golang stats",
				Model: map[string]interface{}{
					"NumGoroutine": numGoroutine,
				},
			}
		},
	},
	{
		Define: "server stats",
		Help:   `Server stats`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			data, err := store.Chatbot.DataGet(ctx.AsUser, ctx.Original, "stats")
			if err != nil {
				return types.TextMsg{Text: "Empty server stats"}
			}

			return types.InfoMsg{
				Title: "Server stats",
				Model: data,
			}
		},
	},
	{
		Define: "agent",
		Help:   `agent url`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return bots.AgentURI(ctx)
		},
	},
}
