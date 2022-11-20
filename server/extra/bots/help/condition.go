package help

import (
	"github.com/tinode/chat/server/extra/ruleset/condition"
	"github.com/tinode/chat/server/extra/types"
)

var conditionRules = []condition.Rule{
	{
		Condition: "RepoMsg",
		Handler: func(ctx types.Context, forwarded types.MsgPayload) types.MsgPayload {
			repo, _ := forwarded.(types.RepoMsg)
			return types.TextMsg{Text: *repo.FullName}
		},
	},
}
