package demo

import (
	"context"
	"errors"
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/command"
	"github.com/tinode/chat/server/extra/types"
)

var handler bot
var config configType

type bot struct {
	initialized bool
}

type configType struct {
	Enabled bool `json:"enabled"`
}

func (bot) Init() error {

	// Check if the handler is already initialized
	if handler.initialized {
		return errors.New("already initialized")
	}

	handler.initialized = true

	if !config.Enabled {
		return nil
	}

	return nil
}

func (bot) IsReady() bool {
	return handler.initialized
}

func (bot) Run(_ map[string]interface{}, content interface{}) ([]map[string]interface{}, []interface{}, error) {
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
	config.Enabled = true
	bots.Register("demo", &handler)
}
