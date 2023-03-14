package helper

import (
	"github.com/tinode/chat/server/extra/ruleset/form"
	"github.com/tinode/chat/server/extra/types"
)

const (
	exampleFormID = "helper_example_form"
)

var formRules = []form.Rule{
	{
		Id: exampleFormID,
		Handler: func(ctx types.Context, values map[string]interface{}) types.MsgPayload {
			return nil
		},
	},
}
