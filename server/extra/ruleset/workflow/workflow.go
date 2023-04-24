package workflow

import "github.com/tinode/chat/server/extra/types"

type Rule struct {
	Id      string
	Version int
	Trigger Trigger
	Step    []Step
}

type Step struct {
	Action func(ctx types.Context, head map[string]interface{}, content interface{}) types.MsgPayload
}

type Trigger struct {
	Type   types.TriggerType
	Define string
}

type Ruleset []Rule

func (r Ruleset) ProcessWorkflow(ctx types.Context, head map[string]interface{}, content interface{}) (types.MsgPayload, error) {
	return nil, nil
}
