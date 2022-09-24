package demo

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/tinode/chat/server/extra/bot"
	"github.com/tinode/chat/server/extra/command"
	"github.com/tinode/chat/server/extra/types"
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

func (demoBot) Run(_ map[string]interface{}, content interface{}) ([]map[string]interface{}, []interface{}, error) {
	in, ok := content.(string)
	if !ok {
		return nil, nil, nil
	}
	ctx := context.Background()
	rs := command.Ruleset(commandRules)
	payloads, err := rs.Help(in)
	if err != nil {
		return nil, nil, err
	}
	if len(payloads) > 0 {
		heads, contents := types.Convert(payloads)
		return heads, contents, nil
	}

	payloads, err = rs.ProcessCommand(ctx, in)
	if err != nil {
		return nil, nil, err
	}

	heads, contents := types.Convert(payloads)
	return heads, contents, nil
}

func init() {
	bot.Register("demo", &handler)
}
