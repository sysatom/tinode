package dev

import (
	"fmt"
	"github.com/tinode/chat/server/extra/ruleset/action"
	"github.com/tinode/chat/server/extra/types"
)

const (
	devActionID = "dev_action"
)

var actionRules = []action.Rule{
	{
		Id: devActionID,
		Handler: map[string]func(ctx types.Context) types.MsgPayload{
			"do1": func(ctx types.Context) types.MsgPayload {
				return types.TextMsg{Text: fmt.Sprintf("do 1 something, action [%s: %d]", ctx.ActionRuleId, ctx.SeqId)}
			},
			"do2": func(ctx types.Context) types.MsgPayload {
				return types.TextMsg{Text: fmt.Sprintf("do 2 something, action [%s: %d]", ctx.ActionRuleId, ctx.SeqId)}
			},
		},
	},
}
