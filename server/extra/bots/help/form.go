package help

import (
	"github.com/tinode/chat/server/extra/ruleset/form"
	"github.com/tinode/chat/server/extra/types"
)

const (
	helpFormID = "help_form"
)

var formRules = []form.Rule{
	{
		Id: helpFormID,
		Handler: func(ctx types.Context, values map[string]interface{}) types.MsgPayload {
			return types.TextMsg{Text: "ok form"}
		},
	},
}
