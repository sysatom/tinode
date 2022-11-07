package demo

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/command"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/logs"
)

const Name = "demo"

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

func (b bot) Run(_ map[string]interface{}, content interface{}) ([]map[string]interface{}, []interface{}, error) {
	if !b.IsReady() {
		// todo error message
		logs.Info.Printf("bot %s unavailable", Name)
		return nil, nil, nil
	}

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
	bots.Register(Name, &handler)
}
