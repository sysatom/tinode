package bots

import (
	"encoding/json"
	"errors"
)

const BotNameSuffix = "_bot"

type Handler interface {
	// Init initializes the bot.
	Init(jsonconf string) error

	// IsReady сhecks if the bot is initialized.
	IsReady() bool

	// Run return bot result
	Run(head map[string]interface{}, content interface{}) ([]map[string]interface{}, []interface{}, error)
}

type configType struct {
	Name   string          `json:"name"`
	Config json.RawMessage `json:"config"`
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

// Init initializes registered handlers.
func Init(jsconfig string) error {
	var config []configType

	if err := json.Unmarshal([]byte(jsconfig), &config); err != nil {
		return errors.New("failed to parse config: " + err.Error())
	}

	for _, cc := range config {
		if bot := handlers[cc.Name]; bot != nil {
			if err := bot.Init(string(cc.Config)); err != nil {
				return err
			}
		}
	}

	return nil
}

func List() map[string]Handler {
	return handlers
}
