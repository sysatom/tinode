package bark

import (
	"context"
	"github.com/tinode/chat/server/extra/command"
	"github.com/tinode/chat/server/extra/types"
)

var commandRules = []command.Rule{
	{
		Define: "version",
		Help:   `Version`,
		Handler: func(ctx context.Context, tokens []*command.Token) []types.MsgPayload {
			return []types.MsgPayload{types.TextMsg{Text: "V1"}}
		},
	},
}
