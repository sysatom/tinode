package store

import (
	"errors"
	"github.com/tinode/chat/server/extra/store/model"
	extraTypes "github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/store/types"
	"time"
)

var adp Adapter

var availableAdapters = make(map[string]Adapter)

func openAdapter() error {
	if adp == nil {
		if len(availableAdapters) >= 1 {
			// Default to the only entry in availableAdapters.
			for _, v := range availableAdapters {
				adp = v
			}
		} else {
			return errors.New("store: db adapter is not specified. Please set `store_config.use_adapter` in `tinode.conf`")
		}
	}

	if adp.IsOpen() {
		return errors.New("store: connection is already opened")
	}

	return adp.Open()
}

func RegisterAdapter(a Adapter) {
	if a == nil {
		panic("store: Register adapter is nil")
	}

	adapterName := a.GetName()
	if _, ok := availableAdapters[adapterName]; ok {
		panic("store: adapter '" + adapterName + "' is already registered")
	}
	availableAdapters[adapterName] = a
}

// PersistentStorageInterface defines methods used for interation with persistent storage.
type PersistentStorageInterface interface {
	Open() error
	Close() error
	IsOpen() bool
	GetAdapter() Adapter
	DbStats() func() interface{}
}

// Store is the main object for interacting with persistent storage.
var Store PersistentStorageInterface

type storeObj struct{}

func (s storeObj) Open() error {
	if err := openAdapter(); err != nil {
		return err
	}
	return nil
}

func (s storeObj) Close() error {
	if adp.IsOpen() {
		return adp.Close()
	}

	return nil
}

func (s storeObj) GetAdapter() Adapter {
	return adp
}

// IsOpen checks if persistent storage connection has been initialized.
func (storeObj) IsOpen() bool {
	if adp != nil {
		return adp.IsOpen()
	}

	return false
}

func (s storeObj) DbStats() func() interface{} {
	if !s.IsOpen() {
		return nil
	}
	return adp.Stats
}

type ChatbotPersistenceInterface interface {
	GetBotUsers() ([]*model.User, error)
	GetNormalUsers() ([]*model.User, error)
	GetGroupTopics(owner types.Uid) ([]*model.Topic, error)
	SearchMessages(uid types.Uid, searchTopic string, filter string) ([]*model.Message, error)
	GetMessage(topic string, seqId int) (model.Message, error)
	GetCredentials() ([]*model.Credential, error)

	DataSet(uid types.Uid, topic, key string, value model.JSON) error
	DataGet(uid types.Uid, topic, key string) (model.JSON, error)
	DataList(uid types.Uid, topic string, filter extraTypes.DataFilter) ([]*model.Data, error)
	DataDelete(uid types.Uid, topic, key string) error
	ConfigSet(uid types.Uid, topic, key string, value model.JSON) error
	ConfigGet(uid types.Uid, topic, key string) (model.JSON, error)
	OAuthSet(oauth model.OAuth) error
	OAuthGet(uid types.Uid, topic, t string) (model.OAuth, error)
	OAuthGetAvailable(t string) ([]model.OAuth, error)
	FormSet(formId string, form model.Form) error
	FormGet(formId string) (model.Form, error)
	ActionSet(topic string, seqId int, action model.Action) error
	ActionGet(topic string, seqId int) (model.Action, error)
	SessionCreate(session model.Session) error
	SessionSet(uid types.Uid, topic string, session model.Session) error
	SessionState(uid types.Uid, topic string, state model.SessionState) error
	SessionGet(uid types.Uid, topic string) (model.Session, error)
	WorkflowCreate(workflow model.Workflow) error
	WorkflowState(uid types.Uid, topic string, workflow model.Workflow) error
	WorkflowStep(uid types.Uid, topic string, workflow model.Workflow) error
	WorkflowGet(uid types.Uid, topic string, flag string) (model.Workflow, error)
	PageSet(pageId string, page model.Page) error
	PageGet(pageId string) (model.Page, error)
	UrlCreate(url model.Url) error
	UrlGetByFlag(flag string) (model.Url, error)
	UrlGetByUrl(url string) (model.Url, error)
	UrlState(flag string, state model.UrlState) error
	UrlViewIncrease(flag string) error
	BehaviorSet(behavior model.Behavior) error
	BehaviorGet(uid types.Uid, flag string) (model.Behavior, error)
	BehaviorList(uid types.Uid) ([]*model.Behavior, error)
	BehaviorIncrease(uid types.Uid, flag string, number int) error
	ParameterSet(flag string, params model.JSON, expiredAt time.Time) error
	ParameterGet(flag string) (model.Parameter, error)
	ParameterDelete(flag string) error

	GetObjectiveByID(id int64) (*model.Objective, error)
	GetObjectiveBySequence(uid types.Uid, topic string, sequence int64) (*model.Objective, error)
	ListObjectives(uid types.Uid, topic string) ([]*model.Objective, error)
	CreateObjective(objective *model.Objective) (int64, error)
	UpdateObjective(objective *model.Objective) error
	DeleteObjective(id int64) error
	DeleteObjectiveBySequence(uid types.Uid, topic string, sequence int64) error
	GetKeyResultByID(id int64) (*model.KeyResult, error)
	GetKeyResultBySequence(uid types.Uid, topic string, sequence int64) (*model.KeyResult, error)
	ListKeyResults(uid types.Uid, topic string) ([]*model.KeyResult, error)
	ListKeyResultsById(id []int64) ([]*model.KeyResult, error)
	ListKeyResultsByObjectiveId(objectiveId int64) ([]*model.KeyResult, error)
	CreateKeyResult(keyResult *model.KeyResult) (int64, error)
	UpdateKeyResult(keyResult *model.KeyResult) error
	DeleteKeyResult(id int64) error
	DeleteKeyResultBySequence(uid types.Uid, topic string, sequence int64) error
	AggregateObjectiveValue(id int64) error
	AggregateKeyResultValue(id int64) error
	CreateKeyResultValue(keyResultValue *model.KeyResultValue) (int64, error)
	GetKeyResultValues(keyResultId int64) ([]*model.KeyResultValue, error)

	CreateTodo(todo *model.Todo) (int64, error)
	ListTodos(uid types.Uid, topic string) ([]*model.Todo, error)
	ListRemindTodos(uid types.Uid, topic string) ([]*model.Todo, error)
	GetTodo(id int64) (*model.Todo, error)
	GetTodoBySequence(uid types.Uid, topic string, sequence int64) (*model.Todo, error)
	CompleteTodo(id int64) error
	CompleteTodoBySequence(uid types.Uid, topic string, sequence int64) error
	UpdateTodo(todo *model.Todo) error
	DeleteTodo(id int64) error
	DeleteTodoBySequence(uid types.Uid, topic string, sequence int64) error

	CreateCounter(counter *model.Counter) (int64, error)
	IncreaseCounter(id, amount int64) error
	DecreaseCounter(id, amount int64) error
	ListCounter(uid types.Uid, topic string) ([]*model.Counter, error)
	GetCounter(id int64) (model.Counter, error)
	GetCounterByFlag(uid types.Uid, topic string, flag string) (model.Counter, error)

	CreateInstruct(instruct *model.Instruct) (int64, error)
	ListInstruct(uid types.Uid, isExpire bool) ([]*model.Instruct, error)
	UpdateInstruct(instruct *model.Instruct) error
}

var Chatbot ChatbotPersistenceInterface

type chatbotMapper struct{}

func (c chatbotMapper) GetBotUsers() ([]*model.User, error) {
	return adp.GetBotUsers()
}

func (c chatbotMapper) GetNormalUsers() ([]*model.User, error) {
	return adp.GetNormalUsers()
}

func (c chatbotMapper) GetGroupTopics(owner types.Uid) ([]*model.Topic, error) {
	return adp.GetGroupTopics(owner)
}

func (c chatbotMapper) SearchMessages(uid types.Uid, searchTopic string, filter string) ([]*model.Message, error) {
	return adp.SearchMessages(uid, searchTopic, filter)
}

func (c chatbotMapper) GetMessage(topic string, seqId int) (model.Message, error) {
	return adp.GetMessage(topic, seqId)
}

func (c chatbotMapper) GetCredentials() ([]*model.Credential, error) {
	return adp.GetCredentials()
}

func (c chatbotMapper) ConfigSet(uid types.Uid, topic, key string, value model.JSON) error {
	return adp.ConfigSet(uid, topic, key, value)
}

func (c chatbotMapper) ConfigGet(uid types.Uid, topic, key string) (model.JSON, error) {
	return adp.ConfigGet(uid, topic, key)
}

func (c chatbotMapper) OAuthSet(oauth model.OAuth) error {
	return adp.OAuthSet(oauth)
}

func (c chatbotMapper) OAuthGet(uid types.Uid, topic, t string) (model.OAuth, error) {
	return adp.OAuthGet(uid, topic, t)
}

func (c chatbotMapper) OAuthGetAvailable(t string) ([]model.OAuth, error) {
	return adp.OAuthGetAvailable(t)
}

func (c chatbotMapper) FormSet(formId string, form model.Form) error {
	return adp.FormSet(formId, form)
}

func (c chatbotMapper) FormGet(formId string) (model.Form, error) {
	return adp.FormGet(formId)
}

func (c chatbotMapper) ActionSet(topic string, seqId int, action model.Action) error {
	return adp.ActionSet(topic, seqId, action)
}

func (c chatbotMapper) ActionGet(topic string, seqId int) (model.Action, error) {
	return adp.ActionGet(topic, seqId)
}

func (c chatbotMapper) SessionCreate(session model.Session) error {
	return adp.SessionCreate(session)
}

func (c chatbotMapper) SessionSet(uid types.Uid, topic string, session model.Session) error {
	return adp.SessionSet(uid, topic, session)
}

func (c chatbotMapper) SessionState(uid types.Uid, topic string, state model.SessionState) error {
	return adp.SessionState(uid, topic, state)
}

func (c chatbotMapper) SessionGet(uid types.Uid, topic string) (model.Session, error) {
	return adp.SessionGet(uid, topic)
}

func (c chatbotMapper) WorkflowCreate(workflow model.Workflow) error {
	return adp.WorkflowCreate(workflow)
}

func (c chatbotMapper) WorkflowState(uid types.Uid, topic string, workflow model.Workflow) error {
	return adp.WorkflowState(uid, topic, workflow)
}

func (c chatbotMapper) WorkflowStep(uid types.Uid, topic string, workflow model.Workflow) error {
	return adp.WorkflowStep(uid, topic, workflow)
}

func (c chatbotMapper) WorkflowGet(uid types.Uid, topic string, flag string) (model.Workflow, error) {
	return adp.WorkflowGet(uid, topic, flag)
}

func (c chatbotMapper) PageSet(pageId string, page model.Page) error {
	return adp.PageSet(pageId, page)
}

func (c chatbotMapper) PageGet(pageId string) (model.Page, error) {
	return adp.PageGet(pageId)
}

func (c chatbotMapper) DataSet(uid types.Uid, topic, key string, value model.JSON) error {
	return adp.DataSet(uid, topic, key, value)
}

func (c chatbotMapper) DataGet(uid types.Uid, topic, key string) (model.JSON, error) {
	return adp.DataGet(uid, topic, key)
}

func (c chatbotMapper) DataList(uid types.Uid, topic string, filter extraTypes.DataFilter) ([]*model.Data, error) {
	return adp.DataList(uid, topic, filter)
}

func (c chatbotMapper) DataDelete(uid types.Uid, topic, key string) error {
	return adp.DataDelete(uid, topic, key)
}

func (c chatbotMapper) UrlCreate(url model.Url) error {
	return adp.UrlCreate(url)
}

func (c chatbotMapper) UrlGetByFlag(flag string) (model.Url, error) {
	return adp.UrlGetByFlag(flag)
}

func (c chatbotMapper) UrlGetByUrl(url string) (model.Url, error) {
	return adp.UrlGetByUrl(url)
}

func (c chatbotMapper) UrlState(flag string, state model.UrlState) error {
	return adp.UrlState(flag, state)
}

func (c chatbotMapper) UrlViewIncrease(flag string) error {
	return adp.UrlViewIncrease(flag)
}

func (c chatbotMapper) BehaviorSet(behavior model.Behavior) error {
	return adp.BehaviorSet(behavior)
}

func (c chatbotMapper) BehaviorGet(uid types.Uid, flag string) (model.Behavior, error) {
	return adp.BehaviorGet(uid, flag)
}

func (c chatbotMapper) BehaviorList(uid types.Uid) ([]*model.Behavior, error) {
	return adp.BehaviorList(uid)
}

func (c chatbotMapper) BehaviorIncrease(uid types.Uid, flag string, number int) error {
	return adp.BehaviorIncrease(uid, flag, number)
}

func (c chatbotMapper) ParameterSet(flag string, params model.JSON, expiredAt time.Time) error {
	return adp.ParameterSet(flag, params, expiredAt)
}

func (c chatbotMapper) ParameterGet(flag string) (model.Parameter, error) {
	return adp.ParameterGet(flag)
}

func (c chatbotMapper) ParameterDelete(flag string) error {
	return adp.ParameterDelete(flag)
}

func (c chatbotMapper) GetObjectiveByID(id int64) (*model.Objective, error) {
	return adp.GetObjectiveByID(id)
}

func (c chatbotMapper) GetObjectiveBySequence(uid types.Uid, topic string, sequence int64) (*model.Objective, error) {
	return adp.GetObjectiveBySequence(uid, topic, sequence)
}

func (c chatbotMapper) ListObjectives(uid types.Uid, topic string) ([]*model.Objective, error) {
	return adp.ListObjectives(uid, topic)
}

func (c chatbotMapper) CreateObjective(objective *model.Objective) (int64, error) {
	return adp.CreateObjective(objective)
}

func (c chatbotMapper) UpdateObjective(objective *model.Objective) error {
	return adp.UpdateObjective(objective)
}

func (c chatbotMapper) DeleteObjective(id int64) error {
	return adp.DeleteObjective(id)
}

func (c chatbotMapper) DeleteObjectiveBySequence(uid types.Uid, topic string, sequence int64) error {
	return adp.DeleteObjectiveBySequence(uid, topic, sequence)
}

func (c chatbotMapper) GetKeyResultByID(id int64) (*model.KeyResult, error) {
	return adp.GetKeyResultByID(id)
}

func (c chatbotMapper) GetKeyResultBySequence(uid types.Uid, topic string, sequence int64) (*model.KeyResult, error) {
	return adp.GetKeyResultBySequence(uid, topic, sequence)
}

func (c chatbotMapper) ListKeyResults(uid types.Uid, topic string) ([]*model.KeyResult, error) {
	return adp.ListKeyResults(uid, topic)
}

func (c chatbotMapper) ListKeyResultsById(id []int64) ([]*model.KeyResult, error) {
	return adp.ListKeyResultsById(id)
}

func (c chatbotMapper) ListKeyResultsByObjectiveId(objectiveId int64) ([]*model.KeyResult, error) {
	return adp.ListKeyResultsByObjectiveId(objectiveId)
}

func (c chatbotMapper) CreateKeyResult(keyResult *model.KeyResult) (int64, error) {
	return adp.CreateKeyResult(keyResult)
}

func (c chatbotMapper) UpdateKeyResult(keyResult *model.KeyResult) error {
	return adp.UpdateKeyResult(keyResult)
}

func (c chatbotMapper) DeleteKeyResult(id int64) error {
	return adp.DeleteKeyResult(id)
}

func (c chatbotMapper) DeleteKeyResultBySequence(uid types.Uid, topic string, sequence int64) error {
	return adp.DeleteKeyResultBySequence(uid, topic, sequence)
}

func (c chatbotMapper) AggregateObjectiveValue(id int64) error {
	return adp.AggregateObjectiveValue(id)
}

func (c chatbotMapper) AggregateKeyResultValue(id int64) error {
	return adp.AggregateKeyResultValue(id)
}

func (c chatbotMapper) CreateKeyResultValue(keyResultValue *model.KeyResultValue) (int64, error) {
	return adp.CreateKeyResultValue(keyResultValue)
}

func (c chatbotMapper) GetKeyResultValues(keyResultId int64) ([]*model.KeyResultValue, error) {
	return adp.GetKeyResultValues(keyResultId)
}

func (c chatbotMapper) CreateTodo(todo *model.Todo) (int64, error) {
	return adp.CreateTodo(todo)
}

func (c chatbotMapper) ListTodos(uid types.Uid, topic string) ([]*model.Todo, error) {
	return adp.ListTodos(uid, topic)
}

func (c chatbotMapper) ListRemindTodos(uid types.Uid, topic string) ([]*model.Todo, error) {
	return adp.ListRemindTodos(uid, topic)
}

func (c chatbotMapper) GetTodo(id int64) (*model.Todo, error) {
	return adp.GetTodo(id)
}

func (c chatbotMapper) GetTodoBySequence(uid types.Uid, topic string, sequence int64) (*model.Todo, error) {
	return adp.GetTodoBySequence(uid, topic, sequence)
}

func (c chatbotMapper) CompleteTodo(id int64) error {
	return adp.CompleteTodo(id)
}

func (c chatbotMapper) CompleteTodoBySequence(uid types.Uid, topic string, sequence int64) error {
	return adp.CompleteTodoBySequence(uid, topic, sequence)
}

func (c chatbotMapper) UpdateTodo(todo *model.Todo) error {
	return adp.UpdateTodo(todo)
}

func (c chatbotMapper) DeleteTodo(id int64) error {
	return adp.DeleteTodo(id)
}

func (c chatbotMapper) DeleteTodoBySequence(uid types.Uid, topic string, sequence int64) error {
	return adp.DeleteTodoBySequence(uid, topic, sequence)
}

func (c chatbotMapper) CreateCounter(counter *model.Counter) (int64, error) {
	return adp.CreateCounter(counter)
}

func (c chatbotMapper) IncreaseCounter(id, amount int64) error {
	return adp.IncreaseCounter(id, amount)
}

func (c chatbotMapper) DecreaseCounter(id, amount int64) error {
	return adp.DecreaseCounter(id, amount)
}

func (c chatbotMapper) ListCounter(uid types.Uid, topic string) ([]*model.Counter, error) {
	return adp.ListCounter(uid, topic)
}

func (c chatbotMapper) GetCounter(id int64) (model.Counter, error) {
	return adp.GetCounter(id)
}

func (c chatbotMapper) GetCounterByFlag(uid types.Uid, topic string, flag string) (model.Counter, error) {
	return adp.GetCounterByFlag(uid, topic, flag)
}

func (c chatbotMapper) CreateInstruct(instruct *model.Instruct) (int64, error) {
	return adp.CreateInstruct(instruct)
}

func (c chatbotMapper) ListInstruct(uid types.Uid, isExpire bool) ([]*model.Instruct, error) {
	return adp.ListInstruct(uid, isExpire)
}

func (c chatbotMapper) UpdateInstruct(instruct *model.Instruct) error {
	return adp.UpdateInstruct(instruct)
}

func init() {
	Store = storeObj{}
	Chatbot = chatbotMapper{}
}
