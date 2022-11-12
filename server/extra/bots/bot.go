package bots

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tinode/chat/server/extra/ruleset/command"
	"github.com/tinode/chat/server/extra/ruleset/form"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/store/model"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/logs"
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

	// Command return bot result
	Command(ctx types.Context, content interface{}) (map[string]interface{}, interface{}, error)

	// Form return bot form result
	Form(ctx types.Context, content interface{}) (map[string]interface{}, interface{}, error)

	// Cron cron script daemon
	Cron(send func(userUid, topicUid serverTypes.Uid, out types.MsgPayload)) error
}

type Base struct{}

func (Base) Command(_ types.Context, _ interface{}) (map[string]interface{}, interface{}, error) {
	return nil, nil, nil
}

func (Base) Form(_ types.Context, _ interface{}) (map[string]interface{}, interface{}, error) {
	return nil, nil, nil
}

func (Base) Cron(_ func(userUid, topicUid serverTypes.Uid, out types.MsgPayload)) error {
	return nil
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

	var id string
	values := make(map[string]interface{})
	if m, ok := msg.Ent[0].Data.Val.(map[string]interface{}); ok {
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

	// check form
	exForm, err := store.Chatbot.FormGet(ctx.FormId)
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
	err = store.Chatbot.FormSet(ctx.FormId, model.Form{Values: values, State: model.FormStateSubmitSuccess})
	if err != nil {
		return nil, nil, err
	}

	if payload == nil {
		return nil, nil, nil
	}
	heads, contents := payload.Convert()
	return heads, contents, nil
}

func StoreForm(ctx types.Context, payload types.MsgPayload) types.MsgPayload {
	formId := types.Id()
	d, err := json.Marshal(payload)
	if err != nil {
		logs.Err.Println(err)
		return types.TextMsg{Text: "store form error"}
	}
	schema := model.JSON{}
	err = schema.Scan(d)
	if err != nil {
		logs.Err.Println(err)
		return types.TextMsg{Text: "store form error"}
	}

	var values model.JSON = make(map[string]interface{})
	if v, ok := payload.(types.FormMsg); ok {
		for _, field := range v.Field {
			values[field.Key] = nil
		}
	}

	// store form
	err = store.Chatbot.FormSet(formId, model.Form{
		FormId: formId,
		Uid:    ctx.AsUser.UserId(),
		Topic:  ctx.Original,
		Schema: schema,
		Values: values,
		State:  model.FormStateCreated,
	})
	if err != nil {
		logs.Err.Println(err)
		return types.TextMsg{Text: "store form error"}
	}

	// store page
	err = store.Chatbot.PageSet(formId, model.Page{
		PageId: formId,
		Uid:    ctx.AsUser.UserId(),
		Topic:  ctx.Original,
		Type:   model.PageForm,
		Schema: schema,
	})
	if err != nil {
		logs.Err.Println(err)
		return types.TextMsg{Text: "store form error"}
	}

	return types.LinkMsg{
		Title: fmt.Sprintf("Form [%s]", formId),
		Url:   fmt.Sprintf("http://127.0.0.1:6060/extra/page/%s", formId), // fixme
	}
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

// Cron registered handlers
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
