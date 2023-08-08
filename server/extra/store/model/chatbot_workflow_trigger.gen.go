// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameWorkflowTrigger = "chatbot_workflow_trigger"

// WorkflowTrigger mapped from table <chatbot_workflow_trigger>
type WorkflowTrigger struct {
	ID         int32       `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UID        string      `gorm:"column:uid;not null" json:"uid"`
	Topic      string      `gorm:"column:topic;not null" json:"topic"`
	WorkflowID int32       `gorm:"column:workflow_id;not null" json:"workflow_id"`
	Type       TriggerType `gorm:"column:type;not null" json:"type"`
	Rule       string      `gorm:"column:rule;not null" json:"rule"`
	Count_     int32       `gorm:"column:count;not null" json:"count"`
	State      int32       `gorm:"column:state;not null" json:"state"`
	CreatedAt  time.Time   `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt  time.Time   `gorm:"column:updated_at;not null" json:"updated_at"`
}

// TableName WorkflowTrigger's table name
func (*WorkflowTrigger) TableName() string {
	return TableNameWorkflowTrigger
}
