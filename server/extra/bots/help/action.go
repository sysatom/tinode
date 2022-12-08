package help

import (
	"fmt"
	"github.com/tinode/chat/server/extra/ruleset/action"
	"github.com/tinode/chat/server/extra/types"
)

const (
	helpActionID = "help_action"
)

var actionRules = []action.Rule{
	{
		Id: helpActionID,
		Handler: func(ctx types.Context, values map[string]interface{}) types.MsgPayload {
			fmt.Println(values)
			return types.TextMsg{Text: fmt.Sprintf("ok, action [%s: %d]", ctx.ActionRuleId, ctx.SeqId)}
		},
	},
}
