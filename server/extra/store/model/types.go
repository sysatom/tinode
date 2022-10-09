package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Config struct {
	ID        uint `gorm:"primaryKey"`
	Uid       int64
	Topic     string
	Key       string
	Value     JSON `gorm:"type:json"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Config) TableName() string {
	return "chatbot_configs"
}

type JSON map[string]interface{}

func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := make(map[string]interface{})
	err := json.Unmarshal(bytes, &result)
	*j = result
	return err
}

func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j JSON) String(key string) (string, bool) {
	if v, ok := j.get(key); ok {
		if t, ok := v.(string); ok {
			return t, ok
		}
	}
	return "", false
}

func (j JSON) Int64(key string) (int64, bool) {
	if v, ok := j.get(key); ok {
		if t, ok := v.(int64); ok {
			return t, ok
		}
	}
	return 0, false
}

func (j JSON) get(key string) (interface{}, bool) {
	v, ok := j[key]
	return v, ok
}
