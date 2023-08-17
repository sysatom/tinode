// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameStep = "chatbot_steps"

// Step mapped from table <chatbot_steps>
type Step struct {
	ID         int32      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UID        string     `gorm:"column:uid;not null" json:"uid"`
	Topic      string     `gorm:"column:topic;not null" json:"topic"`
	JobID      int32      `gorm:"column:job_id;not null" json:"job_id"`
	Action     JSON       `gorm:"column:action;not null" json:"action"`
	Name       string     `gorm:"column:name;not null" json:"name"`
	Describe   string     `gorm:"column:describe;not null" json:"describe"`
	NodeID     string     `gorm:"column:node_id;not null" json:"node_id"`
	Depend     []string   `gorm:"column:depend;type:json;not null;serializer:json" json:"depend"`
	Input      JSON       `gorm:"column:input" json:"input"`
	Output     JSON       `gorm:"column:output" json:"output"`
	Error      string     `gorm:"column:error" json:"error"`
	State      StepState  `gorm:"column:state;not null" json:"state"`
	StartedAt  *time.Time `gorm:"column:started_at" json:"started_at"`
	FinishedAt *time.Time `gorm:"column:finished_at" json:"finished_at"`
	CreatedAt  time.Time  `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"column:updated_at;not null" json:"updated_at"`
}

// TableName Step's table name
func (*Step) TableName() string {
	return TableNameStep
}
