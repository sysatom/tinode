package github

import (
	"encoding/json"
	"errors"
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/cron"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/logs"
	serverTypes "github.com/tinode/chat/server/store/types"
)

const Name = "github"

var handler bot
var Config configType

type bot struct {
	initialized bool
}

type configType struct {
	Enabled bool   `json:"enabled"`
	ID      string `json:"id"`
	Secret  string `json:"secret"`
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

func (b bot) Run(ctx types.Context, content interface{}) (map[string]interface{}, interface{}, error) {
	if !b.IsReady() {
		logs.Info.Printf("bot %s unavailable", Name)
		return nil, nil, nil
	}

	return bots.RunCommand(commandRules, ctx, content)
}

func (b bot) Form(_ types.Context, _ interface{}) (map[string]interface{}, interface{}, error) {
	return nil, nil, nil
}

func (bot) Cron(send func(userUid, topicUid serverTypes.Uid, out types.MsgPayload)) error {
	ruleset := cron.NewCronRuleset(Name, cronRules)
	ruleset.Send = send
	ruleset.Daemon()
	return nil
}

func init() {
	bots.Register(Name, &handler)
}
