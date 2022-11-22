package agent

import "github.com/tinode/chat/server/extra/types"

type Rule struct {
	Id      string
	Handler func(ctx types.Context, content interface{}) types.MsgPayload
}

type Ruleset []Rule

func (r Ruleset) ProcessCondition(ctx types.Context, content interface{}) (types.MsgPayload, error) {
	var result types.MsgPayload
	for _, rule := range r {
		if rule.Id == ctx.AgentId {
			result = rule.Handler(ctx, content)
		}
	}
	return result, nil
}
