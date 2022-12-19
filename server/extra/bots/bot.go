package bots

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tinode/chat/server/auth"
	"github.com/tinode/chat/server/extra/ruleset/action"
	"github.com/tinode/chat/server/extra/ruleset/agent"
	"github.com/tinode/chat/server/extra/ruleset/command"
	"github.com/tinode/chat/server/extra/ruleset/condition"
	"github.com/tinode/chat/server/extra/ruleset/cron"
	"github.com/tinode/chat/server/extra/ruleset/form"
	"github.com/tinode/chat/server/extra/ruleset/session"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/store/model"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/extra/utils"
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

	AuthLevel() auth.Level

	// Input return input result
	Input(ctx types.Context, head map[string]interface{}, content interface{}) (types.MsgPayload, error)

	// Command return bot result
	Command(ctx types.Context, content interface{}) (types.MsgPayload, error)

	// Form return bot form result
	Form(ctx types.Context, values map[string]interface{}) (types.MsgPayload, error)

	// Action return bot action result
	Action(ctx types.Context, option string) (types.MsgPayload, error)

	// Session return bot session result
	Session(ctx types.Context, content interface{}) (types.MsgPayload, error)

	// Cron cron script daemon
	Cron(send func(rcptTo string, uid serverTypes.Uid, out types.MsgPayload)) error

	// Condition run conditional process
	Condition(ctx types.Context, forwarded types.MsgPayload) (types.MsgPayload, error)

	// Group return group result
	Group(ctx types.Context, head map[string]interface{}, content interface{}) (types.MsgPayload, error)

	// Agent return group result
	Agent(ctx types.Context, content interface{}) (types.MsgPayload, error)
}

type Base struct{}

func (Base) AuthLevel() auth.Level {
	return auth.LevelAuth
}

func (Base) Input(_ types.Context, _ map[string]interface{}, _ interface{}) (types.MsgPayload, error) {
	return nil, nil
}

func (Base) Command(_ types.Context, _ interface{}) (types.MsgPayload, error) {
	return nil, nil
}

func (Base) Form(_ types.Context, _ map[string]interface{}) (types.MsgPayload, error) {
	return nil, nil
}

func (Base) Action(_ types.Context, _ string) (types.MsgPayload, error) {
	return nil, nil
}

func (Base) Session(_ types.Context, _ interface{}) (types.MsgPayload, error) {
	return nil, nil
}

func (Base) Cron(_ func(rcptTo string, uid serverTypes.Uid, out types.MsgPayload)) error {
	return nil
}

func (Base) Condition(_ types.Context, _ types.MsgPayload) (types.MsgPayload, error) {
	return nil, nil
}

func (Base) Group(_ types.Context, _ map[string]interface{}, _ interface{}) (types.MsgPayload, error) {
	return nil, nil
}

func (Base) Agent(_ types.Context, _ interface{}) (types.MsgPayload, error) {
	return nil, nil
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

func RunCommand(commandRules []command.Rule, ctx types.Context, content interface{}) (types.MsgPayload, error) {
	in, ok := content.(string)
	if !ok {
		return nil, nil
	}
	rs := command.Ruleset(commandRules)
	payload, err := rs.Help(in)
	if err != nil {
		return nil, err
	}
	if payload != nil {
		return payload, nil
	}

	payload, err = rs.ProcessCommand(ctx, in)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func RunForm(formRules []form.Rule, ctx types.Context, values map[string]interface{}) (types.MsgPayload, error) {
	// check form
	exForm, err := store.Chatbot.FormGet(ctx.FormId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if exForm.ID == 0 {
		return nil, nil
	}
	if exForm.State > model.FormStateCreated {
		return nil, nil
	}

	// process form
	rs := form.Ruleset(formRules)
	payload, err := rs.ProcessForm(ctx, values)
	if err != nil {
		return nil, err
	}

	// is long term
	isLongTerm := false
	for _, rule := range rs {
		if rule.Id == ctx.FormRuleId {
			isLongTerm = rule.IsLongTerm
		}
	}
	if !isLongTerm {
		// store form
		err = store.Chatbot.FormSet(ctx.FormId, model.Form{Values: values, State: model.FormStateSubmitSuccess})
		if err != nil {
			return nil, err
		}

		// store page state
		err = store.Chatbot.PageSet(ctx.FormId, model.Page{State: model.PageStateProcessedSuccess})
		if err != nil {
			return nil, err
		}
	}

	return payload, nil
}

func RunAction(actionRules []action.Rule, ctx types.Context, option string) (types.MsgPayload, error) {
	// check action
	exAction, err := store.Chatbot.ActionGet(ctx.RcptTo, ctx.SeqId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if exAction.ID > 0 && exAction.State > model.ActionStateLongTerm {
		return types.TextMsg{Text: "done"}, nil
	}

	// process action
	rs := action.Ruleset(actionRules)
	payload, err := rs.ProcessAction(ctx, option)
	if err != nil {
		return nil, err
	}

	// is long term
	isLongTerm := false
	for _, rule := range rs {
		if rule.Id == ctx.ActionRuleId {
			isLongTerm = rule.IsLongTerm
		}
	}
	var state model.ActionState
	if !isLongTerm {
		state = model.ActionStateSubmitSuccess
	} else {
		state = model.ActionStateLongTerm
	}
	// store action
	err = store.Chatbot.ActionSet(ctx.RcptTo, ctx.SeqId, model.Action{Uid: ctx.AsUser.UserId(), Value: option, State: state})
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func RunCron(cronRules []cron.Rule, name string, level auth.Level, send func(rcptTo string, uid serverTypes.Uid, out types.MsgPayload)) error {
	ruleset := cron.NewCronRuleset(name, level, cronRules)
	ruleset.Send = send
	ruleset.Daemon()
	return nil
}

func RunCondition(conditionRules []condition.Rule, ctx types.Context, forwarded types.MsgPayload) (types.MsgPayload, error) {
	rs := condition.Ruleset(conditionRules)
	return rs.ProcessCondition(ctx, forwarded)
}

func RunAgent(agentRules []agent.Rule, ctx types.Context, content interface{}) (types.MsgPayload, error) {
	rs := agent.Ruleset(agentRules)
	return rs.ProcessCondition(ctx, content)
}

func RunSession(sessionRules []session.Rule, ctx types.Context, content interface{}) (types.MsgPayload, error) {
	rs := session.Ruleset(sessionRules)
	return rs.ProcessSession(ctx, content)
}

func StoreForm(ctx types.Context, payload types.MsgPayload) types.MsgPayload {
	formId := types.Id().String()
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
		State:  model.PageStateCreated,
	})
	if err != nil {
		logs.Err.Println(err)
		return types.TextMsg{Text: "store form error"}
	}

	return types.LinkMsg{
		Title: fmt.Sprintf("Form [%s]", formId),
		Url:   fmt.Sprintf("%s/extra/page/%s", types.AppUrl(), formId),
	}
}

func StorePage(ctx types.Context, category model.PageType, title string, payload types.MsgPayload) types.MsgPayload {
	pageId := types.Id().String()
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

	// store page
	err = store.Chatbot.PageSet(pageId, model.Page{
		PageId: pageId,
		Uid:    ctx.AsUser.UserId(),
		Topic:  ctx.Original,
		Type:   category,
		Schema: schema,
		State:  model.PageStateCreated,
	})
	if err != nil {
		logs.Err.Println(err)
		return types.TextMsg{Text: "store form error"}
	}

	// fix han compatible styles
	title = fmt.Sprintf("%s %s", category, title)
	if utils.HasHan(title) {
		title = ""
	}

	return types.LinkMsg{
		Title: title,
		Url:   fmt.Sprintf("%s/extra/page/%s", types.AppUrl(), pageId),
	}
}

func SessionStart(ctx types.Context, initValues model.JSON) error {
	sess, err := store.Chatbot.SessionGet(ctx.AsUser, ctx.Original)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if sess.ID > 0 && sess.State == model.SessionStart {
		return errors.New("already a session started")
	}
	var values model.JSON
	values = map[string]interface{}{"val": nil}
	_ = store.Chatbot.SessionCreate(model.Session{
		Uid:    ctx.AsUser.UserId(),
		Topic:  ctx.Original,
		RuleId: ctx.SessionRuleId,
		Init:   initValues,
		Values: values,
		State:  model.SessionStart,
	})
	return nil
}

func SessionDone(ctx types.Context) {
	_ = store.Chatbot.SessionState(ctx.AsUser, ctx.Original, model.SessionDone)
}

func SessionCancel(ctx types.Context) {
	_ = store.Chatbot.SessionState(ctx.AsUser, ctx.Original, model.SessionCancel)
}

func AgentURI(ctx types.Context) types.MsgPayload {
	return types.LinkMsg{
		Title: "Agent",
		Url:   fmt.Sprintf("%s/extra/agent/%d/%d", types.AppUrl(), ctx.AsUser, serverTypes.ParseUserId(ctx.Original)),
	}
}

func CreateShortUrl(text string) (string, error) {
	if utils.IsUrl(text) {
		url, err := store.Chatbot.UrlGetByUrl(text)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", err
		}
		if url.ID > 0 {
			return fmt.Sprintf("%s/u/%s", types.AppUrl(), url.Flag), nil
		}
		flag := strings.ToLower(types.Id().String())
		err = store.Chatbot.UrlCreate(model.Url{
			Flag:  flag,
			Url:   text,
			State: model.UrlStateEnable,
		})
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s/u/%s", types.AppUrl(), flag), nil
	}
	return "", errors.New("error url")
}

const (
	MessageBotIncomingBehavior   = "message_bot_incoming"
	MessageGroupIncomingBehavior = "message_group_incoming"
)

func Behavior(uid serverTypes.Uid, flag string, number int) {
	b, err := store.Chatbot.BehaviorGet(uid, flag)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}
	if b.ID > 0 {
		_ = store.Chatbot.BehaviorIncrease(uid, flag, number)
	} else {
		_ = store.Chatbot.BehaviorSet(model.Behavior{
			Uid:   uid.UserId(),
			Flag:  flag,
			Count: number,
		})
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
func Cron(send func(rcptTo string, uid serverTypes.Uid, out types.MsgPayload)) error {
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
