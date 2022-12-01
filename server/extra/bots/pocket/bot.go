package pocket

import (
	"encoding/json"
	"errors"
	"github.com/tinode/chat/server/drafty"
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/extra/utils"
	"github.com/tinode/chat/server/extra/vendors/pocket"
	"github.com/tinode/chat/server/logs"
	serverTypes "github.com/tinode/chat/server/store/types"
	"gorm.io/gorm"
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

func (b bot) Cron(send func(rcptTo string, uid serverTypes.Uid, out types.MsgPayload)) error {
	return bots.RunCron(cronRules, Name, b.AuthLevel(), send)
}

func (b bot) Input(ctx types.Context, head map[string]interface{}, content interface{}) (types.MsgPayload, error) {
	text, err := drafty.PlainText(content)
	if err != nil {
		return nil, err
	}

	if utils.IsUrl(text) {
		url := text
		oauth, err := store.Chatbot.OAuthGet(ctx.AsUser, ctx.Original, Name)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			logs.Err.Println("bot command pocket oauth", err)
		}
		if oauth.Token == "" {
			return types.TextMsg{Text: "App is unauthorized"}, nil
		}

		provider := pocket.NewPocket(Config.ConsumerKey, "", "", oauth.Token)
		_, err = provider.Add(url)
		if err != nil {
			logs.Err.Println(err)
			return types.TextMsg{Text: "Add error"}, nil
		}

		return types.TextMsg{Text: "ok"}, nil
	}

	return nil, nil
}

func init() {
	bots.Register(Name, &handler)
}
