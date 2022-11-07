package model

import (
	"time"
)

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
