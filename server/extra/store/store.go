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

	DataSet(uid types.Uid, topic, key string, value extraTypes.KV) error
	DataGet(uid types.Uid, topic, key string) (extraTypes.KV, error)
	DataList(uid types.Uid, topic string, filter extraTypes.DataFilter) ([]*model.Data, error)
	DataDelete(uid types.Uid, topic, key string) error
	ConfigSet(uid types.Uid, topic, key string, value extraTypes.KV) error
	ConfigGet(uid types.Uid, topic, key string) (extraTypes.KV, error)
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
	PipelineCreate(pipeline model.Pipeline) error
	PipelineState(uid types.Uid, topic string, pipeline model.Pipeline) error
	PipelineStep(uid types.Uid, topic string, pipeline model.Pipeline) error
	PipelineGet(uid types.Uid, topic string, flag string) (model.Pipeline, error)
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
	ParameterSet(flag string, params extraTypes.KV, expiredAt time.Time) error
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
	CreateReview(review *model.Review) (int64, error)
	UpdateReview(review *model.Review)
	ListReviews(uid types.Uid, topic string) ([]*model.Review, error)
	GetReviewByID(id int64) (*model.Review, error)
	CreateReviewEvaluation(evaluation *model.ReviewEvaluation) (int64, error)
	UpdateReviewEvaluation(evaluation *model.ReviewEvaluation)
	ListReviewEvaluations(uid types.Uid, topic string, reviewID int64) ([]*model.ReviewEvaluation, error)
	GetReviewEvaluationByID(id int64) (*model.ReviewEvaluation, error)
	CreateCycle(cycle *model.Cycle) (int64, error)
	UpdateCycle(cycle *model.Cycle)
	ListCycles(uid types.Uid, topic string) ([]*model.Cycle, error)
	GetCycleByID(id int64) (*model.Cycle, error)

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

var Chatbot Adapter

func Init() {
	Store = storeObj{}
	Chatbot = availableAdapters["mysql"] // default use mysql
}
