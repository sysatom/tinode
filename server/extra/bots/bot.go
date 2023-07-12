package bots

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/tinode/chat/server/auth"
	"github.com/tinode/chat/server/extra/pkg/parser"
	"github.com/tinode/chat/server/extra/ruleset/action"
	"github.com/tinode/chat/server/extra/ruleset/agent"
	"github.com/tinode/chat/server/extra/ruleset/command"
	"github.com/tinode/chat/server/extra/ruleset/condition"
	"github.com/tinode/chat/server/extra/ruleset/cron"
	"github.com/tinode/chat/server/extra/ruleset/event"
	"github.com/tinode/chat/server/extra/ruleset/form"
	"github.com/tinode/chat/server/extra/ruleset/instruct"
	"github.com/tinode/chat/server/extra/ruleset/session"
	"github.com/tinode/chat/server/extra/ruleset/setting"
	"github.com/tinode/chat/server/extra/ruleset/workflow"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/store/model"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/extra/utils"
	"github.com/tinode/chat/server/logs"
	serverTypes "github.com/tinode/chat/server/store/types"
	"gorm.io/gorm"
	"strings"
	"time"
)

const BotNameSuffix = "_bot"

type Handler interface {
	// Init initializes the bot.
	Init(jsonconf json.RawMessage) error

	// IsReady сhecks if the bot is initialized.
	IsReady() bool

	// Bootstrap Lifecycle hook
	Bootstrap() error

	// WebService restful web service
	WebService() *restful.WebService

	// AuthLevel authorizations
	AuthLevel() auth.Level

	// Help return bot help
	Help() (map[string][]string, error)

	// Rules return bot ruleset
	Rules() []interface{}

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
	Cron(send types.SendFunc) error

	// Condition run conditional process
	Condition(ctx types.Context, forwarded types.MsgPayload) (types.MsgPayload, error)

	// Group return group result
	Group(ctx types.Context, head map[string]interface{}, content interface{}) (types.MsgPayload, error)

	// Workflow return workflow result
	Workflow(ctx types.Context, head map[string]interface{}, content interface{}, operate types.WorkflowOperate) (types.MsgPayload, string, int, error)

	// Agent return group result
	Agent(ctx types.Context, content interface{}) (types.MsgPayload, error)

	// Instruct return instruct list
	Instruct() (instruct.Ruleset, error)
}

type Base struct{}

func (Base) Bootstrap() error {
	return nil
}

func (Base) WebService() *restful.WebService {
	return nil
}

func (Base) AuthLevel() auth.Level {
	return auth.LevelAuth
}

func (b Base) Help() (map[string][]string, error) {
	return Help(b.Rules())
}

func (Base) Rules() []interface{} {
	return nil
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

func (Base) Cron(_ types.SendFunc) error {
	return nil
}

func (Base) Condition(_ types.Context, _ types.MsgPayload) (types.MsgPayload, error) {
	return nil, nil
}

func (Base) Group(_ types.Context, _ map[string]interface{}, _ interface{}) (types.MsgPayload, error) {
	return nil, nil
}

func (Base) Workflow(_ types.Context, _ map[string]interface{}, _ interface{}, _ types.WorkflowOperate) (types.MsgPayload, string, int, error) {
	return nil, "", 0, nil
}

func (Base) Agent(_ types.Context, _ interface{}) (types.MsgPayload, error) {
	return nil, nil
}

func (Base) Instruct() (instruct.Ruleset, error) {
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

func Help(rules []interface{}) (map[string][]string, error) {
	result := make(map[string][]string)

	for _, rule := range rules {
		switch v := rule.(type) {
		case []command.Rule:
			// command
			rs := command.Ruleset(v)
			var rows []string
			for _, rule := range rs {
				rows = append(rows, fmt.Sprintf("%s : %s", rule.Define, rule.Help))
			}
			if len(rows) > 0 {
				result["command"] = rows
			}
		case []agent.Rule:
			// agent
			rs := agent.Ruleset(v)
			var rows []string
			for _, rule := range rs {
				rows = append(rows, fmt.Sprintf("%s : %s", rule.Id, rule.Help))
			}
			if len(rows) > 0 {
				result["agent"] = rows
			}
		case []cron.Rule:
			// cron
			rs := v
			var rows []string
			for _, rule := range rs {
				rows = append(rows, fmt.Sprintf("%s : %s", rule.Name, rule.Help))
			}
			if len(rows) > 0 {
				result["cron"] = rows
			}
		}
	}

	return result, nil
}

func RunGroup(eventRules []event.Rule, ctx types.Context, head map[string]interface{}, content interface{}) (types.MsgPayload, error) {
	rs := event.Ruleset(eventRules)
	payload, err := rs.ProcessEvent(ctx, head, content)
	if err != nil {
		return nil, err
	}
	// todo
	if len(payload) > 0 {
		return payload[0], nil
	}
	return nil, nil
}

func HelpWorkflow(workflowRules []workflow.Rule, _ types.Context, _ map[string]interface{}, content interface{}) (types.MsgPayload, error) {
	rs := workflow.Ruleset(workflowRules)
	in, ok := content.(string)
	if ok {
		payload, err := rs.Help(in)
		if err != nil {
			return nil, err
		}
		if payload != nil {
			return payload, nil
		}
	}
	return nil, nil
}

func TriggerWorkflow(workflowRules []workflow.Rule, ctx types.Context, _ map[string]interface{}, content interface{}, trigger types.TriggerType) (string, workflow.Rule, error) {
	rs := workflow.Ruleset(workflowRules)
	in, ok := content.(string)
	if ok {
		rule, err := rs.TriggerWorkflow(ctx, trigger, in)
		if err != nil {
			return "", workflow.Rule{}, err
		}

		workflowFlag := ""
		if ctx.WorkflowFlag == "" {
			// store workflow
			flag, err := StoreWorkflow(ctx, rule, 0)
			if err != nil {
				logs.Err.Println(err)
				return "", workflow.Rule{}, err
			}
			workflowFlag = flag
		}

		return workflowFlag, rule, nil
	}
	return "", workflow.Rule{}, errors.New("error trigger")
}

func ProcessWorkflow(ctx types.Context, workflowRule workflow.Rule, index int) (types.MsgPayload, error) {
	if index < 0 || index > len(workflowRule.Step) {
		return nil, errors.New("error workflow step index")
	}
	if index == len(workflowRule.Step) {
		return types.TextMsg{Text: "Workflow Done"}, SetWorkflowState(ctx, ctx.WorkflowFlag, model.WorkflowDone)
	}
	var payload types.MsgPayload
	step := workflowRule.Step[index]
	switch step.Type {
	case types.FormStep:
		payload = FormMsg(ctx, step.Flag)
	case types.ActionStep:
		payload = ActionMsg(ctx, step.Flag)
	case types.CommandStep:
		for name, handler := range List() {
			if step.Bot != types.Bot(name) {
				continue
			}
			for _, item := range handler.Rules() {
				switch v := item.(type) {
				case []command.Rule:
					for _, rule := range v {
						tokens, err := parser.ParseString(strings.Join(step.Args, " "))
						if err != nil {
							return nil, err
						}
						check, err := parser.SyntaxCheck(rule.Define, tokens)
						if err != nil {
							return nil, err
						}
						if !check {
							continue
						}
						payload = rule.Handler(ctx, tokens)
					}
				}
			}
		}
	case types.InstructStep:
		data := make(map[string]interface{}) // fixme
		for i, arg := range step.Args {
			data[fmt.Sprintf("val%d", i+1)] = arg
		}
		payload = InstructMsg(ctx, step.Flag, data)
	case types.SessionStep:
		data := make(map[string]interface{}) // fixme
		for i, arg := range step.Args {
			data[fmt.Sprintf("val%d", i+1)] = arg
		}
		payload = SessionMsg(ctx, step.Flag, data)
	}

	if payload != nil {
		return payload, nil
	}
	return nil, errors.New("error workflow process")
}

func RunWorkflow(workflowRules []workflow.Rule, ctx types.Context, head map[string]interface{}, content interface{}, operate types.WorkflowOperate) (types.MsgPayload, string, int, error) {
	switch operate {
	case types.WorkflowCommandTriggerOperate:
		payload, err := HelpWorkflow(workflowRules, ctx, head, content)
		if err != nil {
			return nil, "", 0, err
		}
		if payload != nil {
			return payload, "", 0, nil
		}
		flag, rule, err := TriggerWorkflow(workflowRules, ctx, head, content, types.TriggerCommandType)
		if err != nil {
			return nil, "", 0, err
		}
		ctx.WorkflowFlag = flag
		ctx.WorkflowVersion = rule.Version
		payload, err = ProcessWorkflow(ctx, rule, 0)
		if err != nil {
			return nil, "", 0, err
		}
		return payload, flag, rule.Version, SetWorkflowStep(ctx, flag, 1)
	case types.WorkflowProcessOperate:
	case types.WorkflowNextOperate:
		for _, rule := range workflowRules {
			if rule.Id == ctx.WorkflowRuleId {
				payload, err := ProcessWorkflow(ctx, rule, ctx.WorkflowStepIndex)
				if err != nil {
					return nil, "", 0, err
				}
				return payload, ctx.WorkflowFlag, ctx.WorkflowVersion, SetWorkflowStep(ctx, ctx.WorkflowFlag, ctx.WorkflowStepIndex+1)
			}
		}
	}
	return nil, "", 0, nil
}

func StoreWorkflow(ctx types.Context, workflowRule workflow.Rule, index int) (string, error) {
	flag := types.Id().String()
	return flag, store.Chatbot.WorkflowCreate(model.Workflow{
		UID:     ctx.AsUser.UserId(),
		Topic:   ctx.Original,
		Flag:    flag,
		RuleID:  workflowRule.Id,
		Version: int32(workflowRule.Version),
		Step:    int32(index),
		State:   model.WorkflowStart,
	})
}

func SetWorkflowState(ctx types.Context, flag string, state model.WorkflowState) error {
	return store.Chatbot.WorkflowState(ctx.AsUser, ctx.Original, model.Workflow{
		Flag:  flag,
		State: state,
	})
}

func SetWorkflowStep(ctx types.Context, flag string, index int) error {
	return store.Chatbot.WorkflowStep(ctx.AsUser, ctx.Original, model.Workflow{
		Flag: flag,
		Step: int32(index),
	})
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
	err = store.Chatbot.ActionSet(ctx.RcptTo, ctx.SeqId, model.Action{UID: ctx.AsUser.UserId(), Value: option, State: state})
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func RunCron(cronRules []cron.Rule, name string, level auth.Level, send types.SendFunc) error {
	ruleset := cron.NewCronRuleset(name, level, cronRules)
	ruleset.Send = send
	ruleset.Daemon()
	return nil
}

func RunCondition(conditionRules []condition.Rule, ctx types.Context, forwarded types.MsgPayload) (types.MsgPayload, error) {
	rs := condition.Ruleset(conditionRules)
	return rs.ProcessCondition(ctx, forwarded)
}

func RunAgent(agentVersion int, agentRules []agent.Rule, ctx types.Context, content interface{}) (types.MsgPayload, error) {
	rs := agent.Ruleset(agentRules)
	return rs.ProcessAgent(agentVersion, ctx, content)
}

func RunSession(sessionRules []session.Rule, ctx types.Context, content interface{}) (types.MsgPayload, error) {
	rs := session.Ruleset(sessionRules)
	return rs.ProcessSession(ctx, content)
}

func FormMsg(ctx types.Context, id string) types.MsgPayload {
	// get form fields
	formMsg := types.FormMsg{ID: id}
	var title string
	var field []types.FormField
	for _, handler := range List() {
		for _, item := range handler.Rules() {
			switch v := item.(type) {
			case []form.Rule:
				for _, rule := range v {
					if rule.Id == id {
						title = rule.Title
						field = rule.Field

						// default value type
						for index, formField := range field {
							if formField.ValueType == "" {
								switch formField.Type {
								case types.FormFieldText, types.FormFieldPassword, types.FormFieldTextarea,
									types.FormFieldEmail, types.FormFieldUrl:
									field[index].ValueType = types.FormFieldValueString
								case types.FormFieldNumber:
									field[index].ValueType = types.FormFieldValueInt64
								}
							}
						}
					}
				}
			}
		}
	}
	if len(field) <= 0 {
		return types.TextMsg{Text: "form field error"}
	}
	formMsg.Title = title
	formMsg.Field = field

	return StoreForm(ctx, formMsg)
}

func StoreForm(ctx types.Context, payload types.MsgPayload) types.MsgPayload {
	formMsg, ok := payload.(types.FormMsg)
	if !ok {
		return types.TextMsg{Text: "form msg error"}
	}

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

	// set extra
	var extra model.JSON = make(map[string]interface{})
	if ctx.WorkflowFlag != "" {
		extra["workflow_flag"] = ctx.WorkflowFlag
		extra["workflow_version"] = ctx.WorkflowVersion
	}

	// store form
	err = store.Chatbot.FormSet(formId, model.Form{
		FormID: formId,
		UID:    ctx.AsUser.UserId(),
		Topic:  ctx.Original,
		Schema: schema,
		Values: values,
		Extra:  extra,
		State:  model.FormStateCreated,
	})
	if err != nil {
		logs.Err.Println(err)
		return types.TextMsg{Text: "store form error"}
	}

	// store page
	err = store.Chatbot.PageSet(formId, model.Page{
		PageID: formId,
		UID:    ctx.AsUser.UserId(),
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
		Title: fmt.Sprintf("%s Form[%s]", formMsg.Title, formId),
		Url:   fmt.Sprintf("%s/extra/page/%s", types.AppUrl(), formId),
	}
}

func StoreParameter(params model.JSON, expiredAt time.Time) (string, error) {
	flag := strings.ToLower(types.Id().String())
	return flag, store.Chatbot.ParameterSet(flag, params, expiredAt)
}

func ActionMsg(_ types.Context, id string) types.MsgPayload {
	var title string
	var option []string
	for _, handler := range List() {
		for _, item := range handler.Rules() {
			switch v := item.(type) {
			case []action.Rule:
				for _, rule := range v {
					if rule.Id == id {
						title = rule.Title
						option = rule.Option
					}
				}
			}
		}
	}
	if len(option) <= 0 {
		return types.TextMsg{Text: "error action rule id"}
	}

	return types.ActionMsg{
		ID:     id,
		Title:  title,
		Option: option,
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
		PageID: pageId,
		UID:    ctx.AsUser.UserId(),
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

func SessionMsg(ctx types.Context, id string, data map[string]interface{}) types.MsgPayload {
	var title string
	for _, handler := range List() {
		for _, item := range handler.Rules() {
			switch v := item.(type) {
			case []session.Rule:
				for _, rule := range v {
					if rule.Id == id {
						title = rule.Title
					}
				}
			}
		}
	}
	if title == "" {
		return types.TextMsg{Text: "error session id"}
	}

	ctx.SessionRuleId = id
	err := SessionStart(ctx, data)
	if err != nil {
		return types.TextMsg{Text: "session error"}
	}

	return types.TextMsg{Text: title}
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
		UID:    ctx.AsUser.UserId(),
		Topic:  ctx.Original,
		RuleID: ctx.SessionRuleId,
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
			URL:   text,
			State: model.UrlStateEnable,
		})
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s/u/%s", types.AppUrl(), flag), nil
	}
	return "", errors.New("error url")
}

func InstructMsg(ctx types.Context, id string, data map[string]interface{}) types.MsgPayload {
	var botName string
	for name, handler := range List() {
		for _, item := range handler.Rules() {
			switch v := item.(type) {
			case []instruct.Rule:
				for _, rule := range v {
					if rule.Id == id {
						botName = name
					}
				}
			}
		}
	}

	return StoreInstruct(ctx, types.InstructMsg{
		No:       types.Id().String(),
		Object:   model.InstructObjectLinkit,
		Bot:      botName,
		Flag:     id,
		Content:  data,
		Priority: model.InstructPriorityDefault,
		State:    model.InstructCreate,
		ExpireAt: time.Now().Add(time.Hour),
	})
}

func StoreInstruct(ctx types.Context, payload types.MsgPayload) types.MsgPayload {
	msg, ok := payload.(types.InstructMsg)
	if !ok {
		return types.TextMsg{Text: "error instruct msg type"}
	}

	_, err := store.Chatbot.CreateInstruct(&model.Instruct{
		UID:      ctx.AsUser.UserId(),
		No:       msg.No,
		Object:   msg.Object,
		Bot:      msg.Bot,
		Flag:     msg.Flag,
		Content:  msg.Content,
		Priority: msg.Priority,
		State:    msg.State,
		ExpireAt: msg.ExpireAt,
	})
	if err != nil {
		return types.TextMsg{Text: "store instruct error"}
	}

	return types.TextMsg{Text: fmt.Sprintf("Instruct[%s:%s]", msg.Flag, msg.No)}
}

func SettingCovertForm(id string, rule setting.Rule) form.Rule {
	var result form.Rule
	result.Id = fmt.Sprintf("%s_setting", id)
	result.Title = fmt.Sprintf("%s Bot Setting", utils.FirstUpper(id))
	result.Field = []types.FormField{}

	for _, row := range rule {
		result.Field = append(result.Field, types.FormField{
			Key:      row.Key,
			Type:     row.Type,
			Required: true,
			Label:    row.Title,
		})
	}

	result.Handler = func(ctx types.Context, values map[string]interface{}) types.MsgPayload {
		for key, value := range values {
			err := store.Chatbot.ConfigSet(ctx.AsUser, ctx.Original, fmt.Sprintf("%s_%s", id, key), map[string]interface{}{
				"value": value,
			})
			if err != nil {
				return types.TextMsg{Text: fmt.Sprintf("setting [%s] %s error", ctx.FormId, key)}
			}
		}
		return types.TextMsg{Text: fmt.Sprintf("ok, setting [%s]", ctx.FormId)}
	}

	return result
}

func SettingGet(ctx types.Context, id string, key string) (model.JSON, error) {
	return store.Chatbot.ConfigGet(ctx.AsUser, ctx.Original, fmt.Sprintf("%s_%s", id, key))
}

func SettingMsg(ctx types.Context, id string) types.MsgPayload {
	return FormMsg(ctx, fmt.Sprintf("%s_setting", id))
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
			UID:    uid.UserId(),
			Flag:   flag,
			Count_: int32(number),
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
			// default config
			configItem = []byte(`{"enabled": true}`)
		}
		if err := bot.Init(configItem); err != nil {
			return err
		}
	}

	return nil
}

func Bootstrap() error {
	for _, bot := range handlers {
		if !bot.IsReady() {
			continue
		}
		if err := bot.Bootstrap(); err != nil {
			return err
		}
	}
	return nil
}

// Cron registered handlers
func Cron(send func(rcptTo string, uid serverTypes.Uid, out types.MsgPayload, option ...interface{})) error {
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
