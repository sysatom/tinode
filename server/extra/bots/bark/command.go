package bark

import (
	"fmt"
	"github.com/tinode/chat/server/extra/bark"
	"github.com/tinode/chat/server/extra/command"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/logs"
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
		Define: "key",
		Help:   `get device key`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			// get
			v, err := store.Chatbot.ConfigGet(ctx.AsUser, "", bark.BarkDeviceKey)
			if err != nil {
				logs.Err.Println("bot command key", err)
			}
			key, _ := v.String("value")

			return types.TextMsg{Text: fmt.Sprintf("key: %s", key)}
		},
	},
	{
		Define: "key [string]",
		Help:   `Set device key`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			key, _ := tokens[1].Value.String()

			// get
			v, err := store.Chatbot.ConfigGet(ctx.AsUser, "", bark.BarkDeviceKey)
			if err != nil {
				logs.Err.Println("bot command key [string]", err)
			}
			old, _ := v.String("value")

			// set
			err = store.Chatbot.ConfigSet(ctx.AsUser, "", bark.BarkDeviceKey, map[string]interface{}{
				"value": key,
			})
			if err != nil {
				logs.Err.Println("bot command key [string]", err)
			}

			return types.TextMsg{Text: fmt.Sprintf("%s --> %s", old, key)}
		},
	},
}
