package share

import (
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/pkg/parser"
	"github.com/tinode/chat/server/extra/ruleset/command"
	"github.com/tinode/chat/server/extra/store/model"
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
		Define: "input",
		Help:   `submit share`,
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
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
					},
				},
			})
		},
	},
	{
		Define: "share [string]",
		Help:   `Share text`,
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
			text, _ := tokens[1].Value.String()
			return bots.StorePage(ctx, model.PageShare, text, types.TextMsg{Text: text})
		},
	},
}
