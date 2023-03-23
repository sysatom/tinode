package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

type JSON map[string]interface{}

func (j *JSON) Scan(value interface{}) error {
	if bytes, ok := value.([]byte); ok {
		result := make(map[string]interface{})
		err := json.Unmarshal(bytes, &result)
		if err != nil {
			return err
		}
		*j = result
		return nil
	}
	if result, ok := value.(map[string]interface{}); ok {
		*j = result
		return nil
	}
	return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
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
		switch n := v.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			return reflect.ValueOf(n).Convert(reflect.TypeOf(int64(0))).Int(), true
		}
	}
	return 0, false
}

func (j JSON) Uint64(key string) (uint64, bool) {
	if v, ok := j.get(key); ok {
		if t, ok := v.(float64); ok {
			return uint64(t), ok
		}
	}
	return 0, false
}

func (j JSON) Float64(key string) (float64, bool) {
	if v, ok := j.get(key); ok {
		if t, ok := v.(float64); ok {
			return t, ok
		}
	}
	return 0, false
}

func (j JSON) Map(key string) (map[string]interface{}, bool) {
	if v, ok := j.get(key); ok {
		if t, ok := v.(map[string]interface{}); ok {
			return t, ok
		}
	}
	return nil, false
}

func (j JSON) get(key string) (interface{}, bool) {
	v, ok := j[key]
	return v, ok
}
