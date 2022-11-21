package share

import (
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/ruleset/command"
	"github.com/tinode/chat/server/extra/store/model"
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
		Define: "input",
		Help:   `submit share`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return bots.StoreForm(ctx, types.FormMsg{
				ID:    inputFormID,
				Title: "Share Content",
				Field: []types.FormField{
					{
						Key:         "content",
						Type:        types.FormFieldTextarea,
						ValueType:   types.FormFieldValueString,
						Value:       "",
						Label:       "Content",
						Placeholder: "Input content",
						Required:    true,
					},
				},
			})
		},
	},
	{
		Define: "share [string]",
		Help:   `Share text`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			text, _ := tokens[1].Value.String()
			return bots.StorePage(ctx, model.PageShare, text, types.TextMsg{Text: text})
		},
	},
}
