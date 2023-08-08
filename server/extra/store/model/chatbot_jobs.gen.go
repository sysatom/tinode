// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameJob = "chatbot_jobs"

// Job mapped from table <chatbot_jobs>
type Job struct {
	ID         int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UID        string    `gorm:"column:uid;not null" json:"uid"`
	Topic      string    `gorm:"column:topic;not null" json:"topic"`
	WorkflowID int32     `gorm:"column:workflow_id;not null" json:"workflow_id"`
	DagID      int32     `gorm:"column:dag_id;not null" json:"dag_id"`
	TriggerID  int32     `gorm:"column:trigger_id;not null" json:"trigger_id"`
	State      JobState  `gorm:"column:state;not null" json:"state"`
	StartedAt  time.Time `gorm:"column:started_at;not null" json:"started_at"`
	FinishedAt time.Time `gorm:"column:finished_at;not null" json:"finished_at"`
	CreatedAt  time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;not null" json:"updated_at"`
	Steps      []*Step   `gorm:"foreignKey:job_id" json:"steps"`
}

// TableName Job's table name
func (*Job) TableName() string {
	return TableNameJob
}
