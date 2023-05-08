package model

import (
	"encoding/json"
	"time"
)

type User struct {
	ID        uint64     `gorm:"primaryKey"`
	State     int        `gorm:"column:state"`
	Stateat   *time.Time `gorm:"column:stateat"`
	UserAgent string     `gorm:"column:useragent"`
	From      uint64
	Access    JSON        `gorm:"type:json"`
	Public    JSON        `gorm:"type:json"`
	Trusted   JSON        `gorm:"type:json"`
	Tags      interface{} `gorm:"type:json"`
	CreatedAt time.Time   `gorm:"column:createdat"`
	UpdatedAt time.Time   `gorm:"column:updatedat"`
	LastSeen  *time.Time  `gorm:"column:lastseen"`

	// bot
	Fn       string `json:"fn,omitempty"`
	Verified bool   `json:"verified,omitempty"`
}

func (User) TableName() string {
	return "users"
}

type Topic struct {
	ID        uint64     `gorm:"primaryKey"`
	State     int        `gorm:"column:state"`
	Stateat   *time.Time `gorm:"column:stateat"`
	Name      string
	UseBt     bool `gorm:"column:usebt"`
	Owner     uint64
	SeqId     int         `gorm:"column:seqid"`
	DelId     int         `gorm:"column:delid"`
	Access    JSON        `gorm:"type:json"`
	Public    JSON        `gorm:"type:json"`
	Trusted   JSON        `gorm:"type:json"`
	Tags      interface{} `gorm:"type:json"`
	CreatedAt time.Time   `gorm:"column:createdat"`
	UpdatedAt time.Time   `gorm:"column:updatedat"`
	TouchedAt time.Time   `gorm:"column:touchedat"`

	Fn string `json:"fn,omitempty"`
	// channel
	Verified bool `json:"verified,omitempty"`
}

func (Topic) TableName() string {
	return "topics"
}

type Message struct {
	ID        uint64 `gorm:"primaryKey"`
	DelId     int    `gorm:"column:delid"`
	SeqId     int    `gorm:"column:seqid"`
	Topic     string
	From      uint64
	Head      JSON      `gorm:"type:json"`
	Content   JSON      `gorm:"type:json"`
	CreatedAt time.Time `gorm:"column:createdat"`
	UpdatedAt time.Time `gorm:"column:updatedat"`
	DeletedAt time.Time `gorm:"column:deletedat"`

	// search
	Txt string          `json:"txt,omitempty"`
	Raw json.RawMessage `json:"raw,omitempty"`
}

func (Message) TableName() string {
	return "messages"
}

type Credential struct {
	ID        uint64 `gorm:"primaryKey"`
	UserId    int64  `gorm:"column:userid"`
	Method    string
	Value     string
	Synthetic string
	Resp      string
	Done      bool
	Retries   int
	CreatedAt time.Time `gorm:"column:createdat"`
	UpdatedAt time.Time `gorm:"column:updatedat"`
	DeletedAt time.Time `gorm:"column:deletedat"`
}

func (Credential) TableName() string {
	return "credentials"
}

type Config struct {
	ID        uint `gorm:"primaryKey"`
	Uid       string
	Topic     string
	Key       string
	Value     JSON `gorm:"type:json"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Config) TableName() string {
	return "chatbot_configs"
}

type OAuth struct {
	ID        uint64 `gorm:"primaryKey"`
	Uid       string
	Topic     string
	Name      string
	Type      string
	Token     string
	Extra     JSON `gorm:"type:json"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (OAuth) TableName() string {
	return "chatbot_oauth"
}

type Form struct {
	ID        uint64 `gorm:"primaryKey"`
	FormId    string
	Uid       string
	Topic     string
	Schema    JSON `gorm:"type:json"`
	Values    JSON `gorm:"type:json"`
	Extra     JSON `gorm:"type:json"`
	State     FormState
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Form) TableName() string {
	return "chatbot_form"
}

type FormState int

const (
	FormStateUnknown FormState = iota
	FormStateCreated
	FormStateSubmitSuccess
	FormStateSubmitFailed
)

type Action struct {
	ID        uint64 `gorm:"primaryKey"`
	Uid       string
	Topic     string
	SeqId     int `gorm:"column:seqid"`
	Value     string
	State     ActionState
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Action) TableName() string {
	return "chatbot_action"
}

type ActionState int

const (
	ActionStateUnknown ActionState = iota
	ActionStateLongTerm
	ActionStateSubmitSuccess
	ActionStateSubmitFailed
)

type Session struct {
	ID        uint64 `gorm:"primaryKey"`
	Uid       string
	Topic     string
	RuleId    string
	Init      JSON `gorm:"type:json"`
	Values    JSON `gorm:"type:json"`
	State     SessionState
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Session) TableName() string {
	return "chatbot_session"
}

type SessionState int

const (
	SessionStateUnknown SessionState = iota
	SessionStart
	SessionDone
	SessionCancel
)

type Page struct {
	ID        uint64 `gorm:"primaryKey"`
	PageId    string
	Uid       string
	Topic     string
	Type      PageType
	Schema    JSON `gorm:"type:json"`
	State     PageState
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Page) TableName() string {
	return "chatbot_page"
}

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

type Data struct {
	ID        uint64 `gorm:"primaryKey"`
	Uid       string
	Topic     string
	Key       string
	Value     JSON `gorm:"type:json"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Data) TableName() string {
	return "chatbot_data"
}

type Url struct {
	ID        uint64 `gorm:"primaryKey"`
	Flag      string
	Url       string
	State     UrlState
	ViewCount int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Url) TableName() string {
	return "chatbot_url"
}

type UrlState int

const (
	UrlStateUnknown UrlState = iota
	UrlStateEnable
	UrlStateDisable
)

type Behavior struct {
	ID        uint `gorm:"primaryKey"`
	Uid       string
	Flag      string
	Count     int
	Extra     *JSON `gorm:"type:json"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Behavior) TableName() string {
	return "chatbot_behavior"
}

type Instruct struct {
	ID        int `gorm:"primaryKey"`
	No        string
	Uid       string
	Object    InstructObject
	Bot       string
	Flag      string
	Content   JSON `gorm:"type:json"`
	Priority  InstructPriority
	State     InstructState
	ExpireAt  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Instruct) TableName() string {
	return "chatbot_instruct"
}

type InstructState int

const (
	InstructStateUnknown InstructState = iota
	InstructCreate
	InstructDone
	InstructCancel
)

type InstructObject string

const (
	InstructObjectHelper InstructObject = "helper"
)

type InstructPriority int

const (
	InstructPriorityHigh    InstructPriority = 3
	InstructPriorityDefault InstructPriority = 2
	InstructPriorityLow     InstructPriority = 1
)

type Workflow struct {
	ID        uint64 `gorm:"primaryKey"`
	Uid       string
	Topic     string
	Flag      string
	RuleId    string
	Version   int
	Step      int
	Values    JSON `gorm:"type:json"`
	State     WorkflowState
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Workflow) TableName() string {
	return "chatbot_workflow"
}

type WorkflowState int

const (
	WorkflowStateUnknown WorkflowState = iota
	WorkflowStart
	WorkflowDone
	WorkflowCancel
)

type Parameter struct {
	ID        uint64 `gorm:"primaryKey"`
	Flag      string
	Params    JSON `gorm:"type:json"`
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiredAt time.Time
}

func (Parameter) TableName() string {
	return "chatbot_parameter"
}
