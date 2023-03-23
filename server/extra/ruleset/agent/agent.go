package agent

import (
	"errors"
	"github.com/tinode/chat/server/extra/types"
)

type Rule struct {
	Id      string
	Help    string
	Handler func(ctx types.Context, content interface{}) types.MsgPayload
}

type Ruleset []Rule

func (r Ruleset) ProcessAgent(agentVersion int, ctx types.Context, content interface{}) (types.MsgPayload, error) {
	if agentVersion > ctx.AgentVersion {
		return nil, errors.New("agent version too low")
	}
	var result types.MsgPayload
	for _, rule := range r {
		if rule.Id == ctx.AgentId {
			result = rule.Handler(ctx, content)
		}
	}
	return result, nil
}
