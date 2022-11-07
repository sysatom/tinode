package bots

import (
	"encoding/json"
	"errors"
	"github.com/tinode/chat/server/extra/command"
	"github.com/tinode/chat/server/extra/types"
)

const BotNameSuffix = "_bot"

type Handler interface {
	// Init initializes the bot.
	Init(jsonconf json.RawMessage) error

	// IsReady сhecks if the bot is initialized.
	IsReady() bool

	// Run return bot result
	Run(ctx types.Context, head map[string]interface{}, content interface{}) ([]map[string]interface{}, []interface{}, error)
}

type configType struct {
	Name string `json:"name"`
}

var handlers map[string]Handler

func Register(name string, bot Handler) {
	if handlers == nil {
		handlers = make(map[string]Handler)
	}

	if bot == nil {
		panic("Register: bot is nil")
	}
	if _, dup := handlers[name]; dup {
		panic("Register: called twice for bot " + name)
	}
	handlers[name] = bot
}

func RunCommand(commandRules []command.Rule, ctx types.Context, _ map[string]interface{}, content interface{}) ([]map[string]interface{}, []interface{}, error) {
	in, ok := content.(string)
	if !ok {
		return nil, nil, nil
	}
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

// Init initializes registered handlers.
func Init(jsonconf json.RawMessage) error {
	var config []json.RawMessage

	if err := json.Unmarshal(jsonconf, &config); err != nil {
		return errors.New("failed to parse config: " + err.Error())
	}

	configMap := make(map[string]json.RawMessage)
	for _, cc := range config {
		var item configType
		if err := json.Unmarshal(cc, &item); err != nil {
			return errors.New("failed to parse config: " + err.Error())
		}

		configMap[item.Name] = cc
	}
	for name, bot := range handlers {
		var configItem json.RawMessage
		if v, ok := configMap[name]; ok {
			configItem = v
		} else {
			configItem = []byte(`{"enabled": true}`)
		}
		if err := bot.Init(configItem); err != nil {
			return err
		}
	}

	return nil
}

func List() map[string]Handler {
	return handlers
}
