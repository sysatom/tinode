package server

import (
	"fmt"
	"github.com/tinode/chat/server/extra/ruleset/command"
	"github.com/tinode/chat/server/extra/types"
	"runtime"
)

var commandRules = []command.Rule{
	{
		Define: "version",
		Help:   `Version`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return types.TextMsg{Text: "V1"}
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
}
