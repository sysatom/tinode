package demo

import (
	"encoding/json"
	"errors"
	"github.com/tinode/chat/server/extra/bot"
)

var handler demoBot

type demoBot struct {
	initialized bool
}

type configType struct {
	Enabled bool `json:"enabled"`
}

func (demoBot) Init(jsonconf string) error {

	// Check if the handler is already initialized
	if handler.initialized {
		return errors.New("already initialized")
	}

	var config configType
	if err := json.Unmarshal([]byte(jsonconf), &config); err != nil {
		return errors.New("failed to parse config: " + err.Error())
	}

	handler.initialized = true

	if !config.Enabled {
		return nil
	}

	return nil
}

func (demoBot) IsReady() bool {
	return handler.initialized
}

func (demoBot) Run(head map[string]interface{}, content interface{}) (map[string]interface{}, interface{}, error) {
	return head, content, nil
}

func init() {
	bot.Register("demo", &handler)
}
