package store

import (
	"errors"
	"github.com/tinode/chat/server/extra/store/model"
	"github.com/tinode/chat/server/store/types"
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
	DataSet(uid types.Uid, topic, key string, value model.JSON) error
	DataGet(uid types.Uid, topic, key string) (model.JSON, error)
	ConfigSet(uid types.Uid, topic, key string, value model.JSON) error
	ConfigGet(uid types.Uid, topic, key string) (model.JSON, error)
	OAuthSet(oauth model.OAuth) error
	OAuthGet(uid types.Uid, topic, t string) (model.OAuth, error)
	OAuthGetAvailable(t string) ([]model.OAuth, error)
	FormSet(formId string, form model.Form) error
	FormGet(formId string) (model.Form, error)
	PageSet(pageId string, page model.Page) error
	PageGet(pageId string) (model.Page, error)

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
}

var Chatbot ChatbotPersistenceInterface

type chatbotMapper struct{}

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

func init() {
	Store = storeObj{}
	Chatbot = chatbotMapper{}
}
