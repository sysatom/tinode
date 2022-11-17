package pocket

import (
	"encoding/json"
	"errors"
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/logs"
)

const Name = "pocket"

var handler bot
var Config configType

type bot struct {
	initialized bool
	bots.Base
}

type configType struct {
	Enabled     bool   `json:"enabled"`
	ConsumerKey string `json:"consumer_key"`
}

func (bot) Init(jsonconf json.RawMessage) error {

	// Check if the handler is already initialized
	if handler.initialized {
		return errors.New("already initialized")
	}

	if err := json.Unmarshal(jsonconf, &Config); err != nil {
		return errors.New("failed to parse config: " + err.Error())
	}

	if !Config.Enabled {
		logs.Info.Printf("bot %s disabled", Name)
		return nil
	}

	handler.initialized = true

	return nil
}

func (bot) IsReady() bool {
	return handler.initialized
}

func (b bot) Command(ctx types.Context, content interface{}) (types.MsgPayload, error) {
	return bots.RunCommand(commandRules, ctx, content)
}

func init() {
	bots.Register(Name, &handler)
}
