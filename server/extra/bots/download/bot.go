package download

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/queue"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/extra/utils"
	"github.com/tinode/chat/server/logs"
	serverTypes "github.com/tinode/chat/server/store/types"
)

const Name = "download"

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

func (b bot) Input(ctx types.Context, _ map[string]interface{}, content interface{}) (types.MsgPayload, error) {
	text := types.ExtractText(content)
	if utils.IsUrl(text) {
		go func() {
			originalName, filename, err := fileDownload(text)
			if err != nil {
				_ = queue.AsyncMessage(ctx.RcptTo, ctx.Original, types.TextMsg{Text: err.Error()})
				return
			}
			_ = queue.AsyncMessage(ctx.RcptTo, ctx.Original, types.LinkMsg{
				Title: originalName,
				Url:   fmt.Sprintf("%s/d/%s", types.AppUrl(), filename),
			})
		}()
		return types.TextMsg{Text: "background"}, nil
	}
	return nil, nil
}

func (b bot) Group(_ types.Context, _ map[string]interface{}, _ interface{}) (types.MsgPayload, error) {
	return nil, nil
}

func (b bot) Command(ctx types.Context, content interface{}) (types.MsgPayload, error) {
	return bots.RunCommand(commandRules, ctx, content)
}

func (b bot) Cron(send func(rcptTo string, uid serverTypes.Uid, out types.MsgPayload)) error {
	return bots.RunCron(cronRules, Name, b.AuthLevel(), send)
}
