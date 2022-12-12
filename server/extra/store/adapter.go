package store

import (
	"github.com/tinode/chat/server/extra/store/model"
	"github.com/tinode/chat/server/store/types"
)

type Adapter interface {
	// General

	// Open and configure the adapter
	Open() error
	// Close the adapter
	Close() error
	// IsOpen checks if the adapter is ready for use
	IsOpen() bool
	// GetName returns the name of the adapter
	GetName() string
	// Stats DB connection stats object.
	Stats() interface{}

	GetBotUsers() ([]*model.User, error)
	GetNormalUsers() ([]*model.User, error)
	GetGroupTopics(owner types.Uid) ([]*model.Topic, error)
	SearchMessages(uid types.Uid, searchTopic string, filter string) ([]*model.Message, error)
	GetMessage(topic string, seqId int) (model.Message, error)
	GetCredentials() ([]*model.Credential, error)

	// Chatbot

	// DataSet data set
	DataSet(uid types.Uid, topic, key string, value model.JSON) error
	// DataGet data get
	DataGet(uid types.Uid, topic, key string) (model.JSON, error)
	// DataList data list
	DataList(uid types.Uid, topic, prefix string) ([]*model.Data, error)
	// DataDelete data delete
	DataDelete(uid types.Uid, topic, key string) error
	// ConfigSet config set
	ConfigSet(uid types.Uid, topic, key string, value model.JSON) error
	// ConfigGet config get
	ConfigGet(uid types.Uid, topic, key string) (model.JSON, error)
	// OAuthSet oauth set
	OAuthSet(oauth model.OAuth) error
	// OAuthGet oauth get
	OAuthGet(uid types.Uid, topic, t string) (model.OAuth, error)
	// OAuthGetAvailable oauth get available
	OAuthGetAvailable(t string) ([]model.OAuth, error)
	// FormSet form set
	FormSet(formId string, form model.Form) error
	// FormGet form get
	FormGet(formId string) (model.Form, error)
	// ActionSet action set
	ActionSet(topic string, seqId int, action model.Action) error
	// ActionGet action get
	ActionGet(topic string, seqId int) (model.Action, error)
	// SessionCreate session create
	SessionCreate(session model.Session) error
	// SessionSet session set
	SessionSet(uid types.Uid, topic string, session model.Session) error
	// SessionState session set state
	SessionState(uid types.Uid, topic string, state model.SessionState) error
	// SessionGet session get
	SessionGet(uid types.Uid, topic string) (model.Session, error)
	// PageSet page set
	PageSet(pageId string, page model.Page) error
	// PageGet page get
	PageGet(pageId string) (model.Page, error)
	// UrlCreate url create
	UrlCreate(url model.Url) error
	// UrlGetByFlag url get by flag
	UrlGetByFlag(flag string) (model.Url, error)
	// UrlGetByUrl url get by url
	UrlGetByUrl(url string) (model.Url, error)
	// UrlState update url state
	UrlState(flag string, state model.UrlState) error
	// UrlViewIncrease increase url view count
	UrlViewIncrease(flag string) error

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
}
