package notion

import (
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/ruleset/command"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/types"
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
		Define: "config",
		Help:   `Config`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			c1, _ := store.Chatbot.ConfigGet(ctx.AsUser, ctx.Original, TokenKey)
			tokenValue, _ := c1.String("value")

			return bots.StoreForm(ctx, types.FormMsg{
				ID:    configFormID,
				Title: "Config",
				Field: []types.FormField{
					{
						Type:        types.FormFieldText,
						Key:         "token",
						Value:       tokenValue,
						ValueType:   types.FormFieldValueString,
						Label:       "Internal Integration Token",
						Placeholder: "Input token",
					},
				},
			})
		},
	},
}
