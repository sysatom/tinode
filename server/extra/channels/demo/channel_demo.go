package demo

import (
	"encoding/json"
	"errors"
	"github.com/tinode/chat/server/extra/channels"
)

var publisher demoChannel

type demoChannel struct {
	id          string
	initialized bool
}

type configType struct {
	Enabled bool   `json:"enabled"`
	Id      string `json:"id"`
}

func (demoChannel) Init(jsonconf string) error {

	// Check if the handler is already initialized
	if publisher.initialized {
		return errors.New("already initialized")
	}

	var config configType
	if err := json.Unmarshal([]byte(jsonconf), &config); err != nil {
		return errors.New("failed to parse config: " + err.Error())
	}

	publisher.initialized = true

	if !config.Enabled {
		return nil
	}

	if len(config.Id) < 0 {
		return nil
	}
	publisher.id = config.Id

	return nil
}

func (demoChannel) IsReady() bool {
	return publisher.initialized
}

func (demoChannel) Id() string {
	return publisher.id
}

func init() {
	channels.Register("demo", &publisher)
}
