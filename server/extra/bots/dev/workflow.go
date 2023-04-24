package dev

import (
	"github.com/tinode/chat/server/extra/ruleset/workflow"
	"github.com/tinode/chat/server/extra/types"
)

const (
	exampleWorkflowId = "example_workflow"
)

var workflowRules = []workflow.Rule{
	{
		Id: exampleWorkflowId,
		Trigger: workflow.Trigger{
			Type:   types.TriggerCommandType,
			Define: "example [string]",
		},
		Step: []workflow.Step{
			{
				Action: func(ctx types.Context, head map[string]interface{}, content interface{}) types.MsgPayload {
					return types.TextMsg{Text: "step 1"}
				},
			},
			{
				Action: func(ctx types.Context, head map[string]interface{}, content interface{}) types.MsgPayload {
					return types.TextMsg{Text: "step 2"}
				},
			},
		},
	},
}
