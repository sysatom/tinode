package help

import (
	"encoding/json"
	"errors"
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/logs"
)

const Name = "help"

var handler bot

type bot struct {
	initialized bool
}

type configType struct {
	Enabled bool `json:"enabled"`
}

func (bot) Init(jsonconf json.RawMessage) error {

	// Check if the handler is already initialized
	if handler.initialized {
		return errors.New("already initialized")
	}

	var config configType
	if err := json.Unmarshal(jsonconf, &config); err != nil {
		return errors.New("failed to parse config: " + err.Error())
	}

	if !config.Enabled {
		logs.Info.Printf("bot %s disabled", Name)
		return nil
	}

	handler.initialized = true

	return nil
}

func (bot) IsReady() bool {
	return handler.initialized
}

func (b bot) Run(ctx types.Context, head map[string]interface{}, content interface{}) ([]map[string]interface{}, []interface{}, error) {
	if !b.IsReady() {
		// todo error message
		logs.Info.Printf("bot %s unavailable", Name)
		return nil, nil, nil
	}

	return bots.RunCommand(commandRules, ctx, head, content)
}

func init() {
	bots.Register(Name, &handler)
}
