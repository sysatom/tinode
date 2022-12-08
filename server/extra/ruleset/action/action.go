package action

import "github.com/tinode/chat/server/extra/types"

type Rule struct {
	Id         string
	IsLongTerm bool
	Handler    func(ctx types.Context, values map[string]interface{}) types.MsgPayload
}

type Ruleset []Rule

func (r Ruleset) ProcessAction(ctx types.Context, values map[string]interface{}) (types.MsgPayload, error) {
	var result types.MsgPayload
	for _, rule := range r {
		if rule.Id == ctx.ActionRuleId {
			result = rule.Handler(ctx, values)
		}
	}
	return result, nil
}
