package model

import (
	"time"
)

type Message struct {
	ID        uint `gorm:"primaryKey"`
	DelId     int  `gorm:"column:delid"`
	SeqId     int  `gorm:"column:seqid"`
	Topic     string
	From      int64
	Head      JSON      `gorm:"type:json"`
	Content   JSON      `gorm:"type:json"`
	CreatedAt time.Time `gorm:"column:createdat"`
	UpdatedAt time.Time `gorm:"column:updatedat"`
	DeletedAt time.Time `gorm:"column:deletedat"`
}

func (Message) TableName() string {
	return "messages"
}

type Credential struct {
	ID        uint  `gorm:"primaryKey"`
	UserId    int64 `gorm:"column:userid"`
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
	ID        uint `gorm:"primaryKey"`
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
	ID        uint `gorm:"primaryKey"`
	FormId    string
	Uid       string
	Topic     string
	Schema    JSON `gorm:"type:json"`
	Values    JSON `gorm:"type:json"`
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

type Page struct {
	ID        uint `gorm:"primaryKey"`
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
	PageForm  PageType = "form"
	PageChart PageType = "chart"
	PageTable PageType = "table"
	PageOkr   PageType = "okr"
	PageShare PageType = "share"
)

type PageState int

const (
	PageStateUnknown PageState = iota
	PageStateCreated
	PageStateProcessedSuccess
	PageStateProcessedFailed
)

type Data struct {
	ID        uint `gorm:"primaryKey"`
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
