package gpt

import (
	"encoding/json"
	"errors"
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/pkg/flog"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/extra/vendors/openai"
)

const Name = "gpt"

const ApiKey = "openai_key"

var handler bot

func init() {
	bots.Register(Name, &handler)
}

type bot struct {
	initialized bool
	bots.Base
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
		flog.Info("bot %s disabled", Name)
		return nil
	}

	handler.initialized = true

	return nil
}

func (bot) IsReady() bool {
	return handler.initialized
}

func (b bot) Input(ctx types.Context, _ types.KV, context interface{}) (types.MsgPayload, error) {
	// key
	v, err := store.Chatbot.ConfigGet(ctx.AsUser, ctx.Original, ApiKey)
	if err != nil {
		flog.Error(err)
	}
	key, _ := v.String("value")

	// input
	text := ""
	if v, ok := context.(string); ok {
		text = v
	}
	if text == "" {
		return types.TextMsg{Text: "input error"}, nil
	}

	client := openai.NewOpenAI(key)
	resp, err := client.Chat(text)
	if err != nil || resp == nil {
		return types.TextMsg{Text: "api error"}, nil
	}

	if len(resp.Choices) > 0 {
		return types.TextMsg{Text: resp.Choices[0].Message.Content}, nil
	}

	return types.TextMsg{Text: "error"}, nil
}

func (b bot) Command(ctx types.Context, content interface{}) (types.MsgPayload, error) {
	return bots.RunCommand(commandRules, ctx, content)
}
