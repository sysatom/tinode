package bots

import (
	"encoding/json"
)

const BotNameSuffix = "_bot"

type Handler interface {
	// Init initializes the bot.
	Init() error

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
func Init() error {
	for _, bot := range handlers {
		if err := bot.Init(); err != nil {
			return err
		}
	}

	return nil
}

func List() map[string]Handler {
	return handlers
}
