package model

import "time"

type Objective struct {
	Id           int64     `json:"id,omitempty"  gorm:"primaryKey"`
	Uid          string    `json:"uid,omitempty" `
	Topic        string    `json:"topic,omitempty" `
	Sequence     int64     `json:"sequence,omitempty" `
	Title        string    `json:"title,omitempty" `
	Memo         string    `json:"memo,omitempty" `
	Motive       string    `json:"motive,omitempty" `
	Feasibility  string    `json:"feasibility,omitempty" `
	IsPlan       bool      `json:"is_plan,omitempty" `
	PlanStart    int64     `json:"plan_start,omitempty" `
	PlanEnd      int64     `json:"plan_end,omitempty" `
	TotalValue   int32     `json:"total_value,omitempty" `
	CurrentValue int32     `json:"current_value,omitempty" `
	CreatedAt    time.Time `json:"created_at,omitempty" `
	UpdatedAt    time.Time `json:"updated_at,omitempty" `
	Tag          string    `json:"tag,omitempty" gorm:"-"`
}

func (Objective) TableName() string {
	return "chatbot_objectives"
}

type KeyResult struct {
	Id           int64         `json:"id,omitempty"  gorm:"primaryKey"`
	Uid          string        `json:"uid,omitempty" `
	Topic        string        `json:"topic,omitempty" `
	ObjectiveId  int64         `json:"objective_id,omitempty" `
	Sequence     int64         `json:"sequence,omitempty" `
	Title        string        `json:"title,omitempty" `
	Memo         string        `json:"memo,omitempty" `
	InitialValue int32         `json:"initial_value,omitempty" `
	TargetValue  int32         `json:"target_value,omitempty" `
	CurrentValue int32         `json:"current_value,omitempty" `
	ValueMode    ValueModeType `json:"value_mode,omitempty" `
	CreatedAt    time.Time     `json:"created_at,omitempty" `
	UpdatedAt    time.Time     `json:"updated_at,omitempty" `
	Tag          string        `json:"tag,omitempty" gorm:"-"`
}

func (KeyResult) TableName() string {
	return "chatbot_key_results"
}

type ValueModeType string

const (
	ValueSumMode  ValueModeType = "sum"
	ValueLastMode ValueModeType = "last"
	ValueAvgMode  ValueModeType = "avg"
	ValueMaxMode  ValueModeType = "max"
)

type KeyResultValue struct {
	Id          int64     `json:"id,omitempty"  gorm:"primaryKey"`
	KeyResultId int64     `json:"key_result_id,omitempty" `
	Value       int32     `json:"value,omitempty" `
	CreatedAt   time.Time `json:"created_at,omitempty" `
	UpdatedAt   time.Time `json:"updated_at,omitempty" `
}

func (KeyResultValue) TableName() string {
	return "chatbot_key_result_values"
}

type Todo struct {
	Id             int64     `json:"id,omitempty" gorm:"primaryKey"`
	Uid            string    `json:"uid,omitempty" `
	Topic          string    `json:"topic,omitempty" `
	Sequence       int64     `json:"sequence,omitempty" `
	Content        string    `json:"content,omitempty" `
	Category       string    `json:"category,omitempty" `
	Remark         string    `json:"remark,omitempty" `
	Priority       int64     `json:"priority,omitempty" `
	IsRemindAtTime bool      `json:"is_remind_at_time,omitempty" `
	RemindAt       int64     `json:"remind_at,omitempty" `
	RepeatMethod   string    `json:"repeat_method,omitempty" `
	RepeatRule     string    `json:"repeat_rule,omitempty" `
	RepeatEndAt    int64     `json:"repeat_end_at,omitempty" `
	Complete       bool      `json:"complete" `
	CreatedAt      time.Time `json:"created_at,omitempty" `
	UpdatedAt      time.Time `json:"updated_at,omitempty" `
}

func (Todo) TableName() string {
	return "chatbot_todos"
}

type Counter struct {
	Id        int64     `json:"id,omitempty" gorm:"primaryKey"`
	Uid       string    `json:"uid,omitempty" `
	Topic     string    `json:"topic,omitempty" `
	Flag      string    `json:"flag,omitempty"`
	Digit     int64     `json:"digit,omitempty"`
	Status    int32     `json:"status,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func (Counter) TableName() string {
	return "chatbot_counters"
}
