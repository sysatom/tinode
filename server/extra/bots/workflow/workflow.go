package workflow

import (
	"github.com/tinode/chat/server/extra/ruleset/workflow"
	"github.com/tinode/chat/server/extra/types/schema"
)

const (
	exampleWorkflowId = "example_workflow"
)

var workflowRules = []workflow.Rule{
	{
		Id:      exampleWorkflowId,
		Help:    "example workflow",
		Trigger: schema.CommandTrigger("example [string]"),
		Step: schema.Step(
			schema.Session("guess_session", "100"),
			schema.Instruct("dev_example"),
			schema.Command(schema.Bot("dev"), "rand", "1", "100"),
			schema.Action("dev_action"),
			schema.Form("dev_form"),
		),
	},
}
