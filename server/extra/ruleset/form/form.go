package form

import "github.com/tinode/chat/server/extra/types"

type Rule struct {
	Id      string
	Handler func(ctx types.Context, values map[string]interface{}) types.MsgPayload
}

type Ruleset []Rule

func (r Ruleset) ProcessForm(ctx types.Context, values map[string]interface{}) (types.MsgPayload, error) {
	var result types.MsgPayload
	for _, rule := range r {
		if rule.Id == ctx.FormRuleId {
			result = rule.Handler(ctx, values)
		}
	}
	return result, nil
}
