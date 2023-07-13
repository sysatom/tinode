package subscribe

import (
	"github.com/tinode/chat/server/extra/pkg/channels"
	"github.com/tinode/chat/server/extra/pkg/parser"
	"github.com/tinode/chat/server/extra/ruleset/command"
	"github.com/tinode/chat/server/extra/types"
)

var commandRules = []command.Rule{
	{
		Define: "info",
		Help:   `Bot info`,
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
			return nil
		},
	},
	{
		Define: "list",
		Help:   `List subscribe`,
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
			return types.InfoMsg{
				Title: "Subscribes",
				Model: channels.List(),
			}
		},
	},
}
