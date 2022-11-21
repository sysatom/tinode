package share

import (
	"fmt"
	"github.com/tinode/chat/server/extra/ruleset/form"
	"github.com/tinode/chat/server/extra/types"
)

const (
	inputFormID = "input_form"
)

var formRules = []form.Rule{
	{
		Id:         inputFormID,
		IsLongTerm: true,
		Handler: func(ctx types.Context, values map[string]interface{}) types.MsgPayload {
			return types.TextMsg{Text: fmt.Sprintf("%s", values["content"])}
		},
	},
}
