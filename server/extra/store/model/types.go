package model

type FormState int

const (
	FormStateUnknown FormState = iota
	FormStateCreated
	FormStateSubmitSuccess
	FormStateSubmitFailed
)

type ActionState int

const (
	ActionStateUnknown ActionState = iota
	ActionStateLongTerm
	ActionStateSubmitSuccess
	ActionStateSubmitFailed
)

type SessionState int

const (
	SessionStateUnknown SessionState = iota
	SessionStart
	SessionDone
	SessionCancel
)

type PageType string

const (
	PageForm     PageType = "form"
	PageChart    PageType = "chart"
	PageTable    PageType = "table"
	PageOkr      PageType = "okr"
	PageShare    PageType = "share"
	PageJson     PageType = "json"
	PageHtml     PageType = "html"
	PageMarkdown PageType = "markdown"
)

type PageState int

const (
	PageStateUnknown PageState = iota
	PageStateCreated
	PageStateProcessedSuccess
	PageStateProcessedFailed
)

type UrlState int

const (
	UrlStateUnknown UrlState = iota
	UrlStateEnable
	UrlStateDisable
)

type InstructState int

const (
	InstructStateUnknown InstructState = iota
	InstructCreate
	InstructDone
	InstructCancel
)

type InstructObject string

const (
	InstructObjectLinkit InstructObject = "linkit"
)

type InstructPriority int

const (
	InstructPriorityHigh    InstructPriority = 3
	InstructPriorityDefault InstructPriority = 2
	InstructPriorityLow     InstructPriority = 1
)

type WorkflowState int

const (
	WorkflowStateUnknown WorkflowState = iota
	WorkflowStart
	WorkflowDone
	WorkflowCancel
)

type ValueModeType string

const (
	ValueSumMode  ValueModeType = "sum"
	ValueLastMode ValueModeType = "last"
	ValueAvgMode  ValueModeType = "avg"
	ValueMaxMode  ValueModeType = "max"
)
