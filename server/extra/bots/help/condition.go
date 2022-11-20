package help

import (
	"github.com/tinode/chat/server/extra/ruleset/condition"
	"github.com/tinode/chat/server/extra/types"
)

var conditionRules = []condition.Rule{
	{
		Condition: "TextMsg",
		Handler: func(ctx types.Context, forwarded types.MsgPayload) types.MsgPayload {
			return nil
		},
	},
}
