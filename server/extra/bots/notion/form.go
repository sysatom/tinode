package notion

import (
	"fmt"
	"github.com/tinode/chat/server/extra/ruleset/form"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/logs"
)

const (
	configFormID = "config_form"
)

var formRules = []form.Rule{
	{
		Id: configFormID,
		Handler: func(ctx types.Context, values map[string]interface{}) types.MsgPayload {
			for key, value := range values {
				switch key {
				case "token":
					err := store.Chatbot.ConfigSet(ctx.AsUser, ctx.Original, TokenKey, map[string]interface{}{
						"value": value.(string),
					})
					if err != nil {
						logs.Err.Println(err)
					}
				}
			}
			return types.TextMsg{Text: fmt.Sprintf("ok, form [%s]", ctx.FormId)}
		},
	},
}
