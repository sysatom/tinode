package types

type TriggerType string

const (
	TriggerCommandType TriggerType = "command"
)

type StepType string

const (
	ActionStep   StepType = "action"
	CommandStep  StepType = "command"
	FormStep     StepType = "form"
	InstructStep StepType = "instruct"
	SessionStep  StepType = "session"
)

type Step struct {
	Type StepType
	Bot  Bot
	Flag string
	Args []string
}

type Trigger struct {
	Type   TriggerType
	Define string
}

type WorkflowOperate string

const (
	WorkflowCommandTriggerOperate WorkflowOperate = "command_trigger"
	WorkflowProcessOperate        WorkflowOperate = "workflow_process"
)

type Bot string
