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
			schema.Form("dev_form"),
			schema.Action("dev_action"),
			schema.Instruct("dev_example"),
			schema.Session("guess_session"),
			schema.Command(schema.Bot("dev"), "rand", "1", "2"),
			schema.Condition(schema.Bot("dev"), "RepoMsg"),
		),
	},
}
