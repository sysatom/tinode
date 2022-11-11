package bots

import (
	"encoding/json"
	"errors"
	"github.com/tinode/chat/server/extra/ruleset/command"
	"github.com/tinode/chat/server/extra/ruleset/form"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/store/model"
	"github.com/tinode/chat/server/extra/types"
	serverTypes "github.com/tinode/chat/server/store/types"
	"gorm.io/gorm"
	"strings"
)

const BotNameSuffix = "_bot"

type Handler interface {
	// Init initializes the bot.
	Init(jsonconf json.RawMessage) error

	// IsReady сhecks if the bot is initialized.
	IsReady() bool

	// Run return bot result
	Run(ctx types.Context, content interface{}) (map[string]interface{}, interface{}, error)

	// Form return bot form result
	Form(ctx types.Context, content interface{}) (map[string]interface{}, interface{}, error)

	// Cron cron script daemon
	Cron(send func(userUid, topicUid serverTypes.Uid, out types.MsgPayload)) error
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

func RunCommand(commandRules []command.Rule, ctx types.Context, content interface{}) (map[string]interface{}, interface{}, error) {
	in, ok := content.(string)
	if !ok {
		return nil, nil, nil
	}
	rs := command.Ruleset(commandRules)
	payload, err := rs.Help(in)
	if err != nil {
		return nil, nil, err
	}
	if payload != nil {
		heads, contents := payload.Convert()
		return heads, contents, nil
	}

	payload, err = rs.ProcessCommand(ctx, in)
	if err != nil {
		return nil, nil, err
	}
	if payload == nil {
		return nil, nil, nil
	}

	heads, contents := payload.Convert()
	return heads, contents, nil
}

func RunForm(formRules []form.Rule, ctx types.Context, content interface{}) (map[string]interface{}, interface{}, error) {
	var msg types.ChatMessage
	d, err := json.Marshal(content)
	if err != nil {
		return nil, nil, err
	}
	err = json.Unmarshal(d, &msg)
	if err != nil {
		return nil, nil, err
	}

	if len(msg.Ent) > 0 {
		if msg.Ent[0].Tp != "EX" {
			return nil, nil, nil
		}
	}
	var seq int
	var id string
	values := make(map[string]interface{})
	if m, ok := msg.Ent[0].Data.Val.(map[string]interface{}); ok {
		if v, ok := m["seq"]; ok {
			if vv, ok := v.(float64); ok {
				seq = int(vv)
			}
		}
		if v, ok := m["resp"]; ok {
			if vv, ok := v.(map[string]interface{}); ok {
				for s := range vv {
					ss := strings.Split(s, "|")
					if len(ss) == 2 {
						id = ss[0]
						values[ss[1]] = vv[s]
					}
				}
			}
		}
	}

	ctx.FormId = id
	ctx.SeqId = seq

	// check form
	exForm, err := store.Chatbot.FormGet(ctx.AsUser, ctx.Original, ctx.SeqId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, err
	}
	if exForm.ID > 0 {
		return nil, nil, nil
	}

	// process form
	rs := form.Ruleset(formRules)
	payload, err := rs.ProcessForm(ctx, values)
	if err != nil {
		return nil, nil, err
	}

	// store form
	err = store.Chatbot.FormSet(ctx.AsUser, ctx.Original, ctx.SeqId, values, int(model.FormStateSuccess))
	if err != nil {
		return nil, nil, err
	}

	if payload == nil {
		return nil, nil, nil
	}
	heads, contents := payload.Convert()
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

func Cron(send func(userUid, topicUid serverTypes.Uid, out types.MsgPayload)) error {
	for _, bot := range handlers {
		if err := bot.Cron(send); err != nil {
			return err
		}
	}
	return nil
}

func List() map[string]Handler {
	return handlers
}
