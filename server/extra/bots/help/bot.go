package help

import (
	"encoding/json"
	"errors"
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/logs"
	serverTypes "github.com/tinode/chat/server/store/types"
)

const Name = "help"

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
		logs.Info.Printf("bot %s disabled", Name)
		return nil
	}

	handler.initialized = true

	return nil
}

func (bot) IsReady() bool {
	return handler.initialized
}

func (b bot) Help() (map[string][]string, error) {
	return bots.Help(commandRules, agentRules, cronRules)
}

func (b bot) Input(_ types.Context, _ map[string]interface{}, _ interface{}) (types.MsgPayload, error) {
	return types.TextMsg{Text: "Input"}, nil
}

func (b bot) Command(ctx types.Context, content interface{}) (types.MsgPayload, error) {
	return bots.RunCommand(commandRules, ctx, content)
}

func (b bot) Form(ctx types.Context, values map[string]interface{}) (types.MsgPayload, error) {
	return bots.RunForm(formRules, ctx, values)
}

func (b bot) Action(ctx types.Context, option string) (types.MsgPayload, error) {
	return bots.RunAction(actionRules, ctx, option)
}

func (b bot) Cron(send func(rcptTo string, uid serverTypes.Uid, out types.MsgPayload)) error {
	return bots.RunCron(cronRules, Name, b.AuthLevel(), send)
}

func (b bot) Condition(ctx types.Context, forwarded types.MsgPayload) (types.MsgPayload, error) {
	return bots.RunCondition(conditionRules, ctx, forwarded)
}

func (b bot) Group(_ types.Context, _ map[string]interface{}, _ interface{}) (types.MsgPayload, error) {
	return types.TextMsg{Text: "Group"}, nil
}

func (b bot) Agent(ctx types.Context, content interface{}) (types.MsgPayload, error) {
	return bots.RunAgent(agentRules, ctx, content)
}

func (b bot) Session(ctx types.Context, content interface{}) (types.MsgPayload, error) {
	return bots.RunSession(sessionRules, ctx, content)
}
