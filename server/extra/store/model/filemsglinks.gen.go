// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameFilemsglink = "filemsglinks"

// Filemsglink mapped from table <filemsglinks>
type Filemsglink struct {
	ID        int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Createdat time.Time `gorm:"column:createdat;not null" json:"createdat"`
	Fileid    int64     `gorm:"column:fileid;not null" json:"fileid"`
	Msgid     int32     `gorm:"column:msgid" json:"msgid"`
	Topic     string    `gorm:"column:topic" json:"topic"`
	Userid    int64     `gorm:"column:userid" json:"userid"`
}

// TableName Filemsglink's table name
func (*Filemsglink) TableName() string {
	return TableNameFilemsglink
}
