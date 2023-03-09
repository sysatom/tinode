package gpt

import (
	"fmt"
	"github.com/tinode/chat/server/extra/ruleset/command"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/extra/utils"
	"github.com/tinode/chat/server/logs"
)

var commandRules = []command.Rule{
	{
		Define: "info",
		Help:   `Bot info`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return nil
		},
	},
	{
		Define: "key",
		Help:   `get api key`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			// get
			v, err := store.Chatbot.ConfigGet(ctx.AsUser, ctx.Original, ApiKey)
			if err != nil {
				logs.Err.Println("bot command key", err)
			}
			key, _ := v.String("value")

			return types.TextMsg{Text: fmt.Sprintf("key: %s", utils.Masker(key, 3))}
		},
	},
	{
		Define: "key [string]",
		Help:   `Set api key`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			key, _ := tokens[1].Value.String()

			// get
			v, err := store.Chatbot.ConfigGet(ctx.AsUser, ctx.Original, ApiKey)
			if err != nil {
				logs.Err.Println("bot command key [string]", err)
			}
			old, _ := v.String("value")

			// set
			err = store.Chatbot.ConfigSet(ctx.AsUser, ctx.Original, ApiKey, map[string]interface{}{
				"value": key,
			})
			if err != nil {
				logs.Err.Println("bot command key [string]", err)
			}

			return types.TextMsg{Text: fmt.Sprintf("%s --> %s", old, key)}
		},
	},
}
