package help

import (
	"fmt"
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
			fmt.Println(values)
			return types.TextMsg{Text: fmt.Sprintf("ok, form [%s]", ctx.FormId)}
		},
	},
}
