package url

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/store/model"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/extra/utils"
	"github.com/tinode/chat/server/logs"
	"gorm.io/gorm"
	"strings"
)

const Name = "url"

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

func (b bot) Input(_ types.Context, _ map[string]interface{}, content interface{}) (types.MsgPayload, error) {
	text := ""
	if m, ok := content.(map[string]interface{}); ok {
		if t, ok := m["txt"]; ok {
			if s, ok := t.(string); ok {
				text = s
			}
		}
	} else if s, ok := content.(string); ok {
		text = s
	}
	if utils.IsUrl(text) {
		url, err := store.Chatbot.UrlGetByUrl(text)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return types.TextMsg{Text: "query url error"}, nil
		}
		if url.ID > 0 {
			return types.LinkMsg{Url: fmt.Sprintf("%s/u/%s", types.AppUrl(), url.Flag)}, nil
		}
		flag := strings.ToLower(types.Id().String())
		err = store.Chatbot.UrlCreate(model.Url{
			Flag:  flag,
			Url:   text,
			State: model.UrlStateEnable,
		})
		if err != nil {
			return types.TextMsg{Text: "create error"}, nil
		}
		return types.LinkMsg{Url: fmt.Sprintf("%s/u/%s", types.AppUrl(), flag)}, nil
	} else {
		url, err := store.Chatbot.UrlGetByFlag(text)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return types.TextMsg{Text: "query url error"}, nil
		}
		if url.ID > 0 {
			return types.LinkMsg{Url: url.Url}, nil
		}
		return types.TextMsg{Text: "empty"}, nil
	}
}

func (b bot) Group(_ types.Context, _ map[string]interface{}, _ interface{}) (types.MsgPayload, error) {
	return nil, nil
}

func (b bot) Command(ctx types.Context, content interface{}) (types.MsgPayload, error) {
	return bots.RunCommand(commandRules, ctx, content)
}
