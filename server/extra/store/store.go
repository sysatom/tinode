package store

import (
	"context"
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

	GetObjectiveByID(ctx context.Context, id int64) (*model.Objective, error)
	GetObjectiveBySequence(ctx context.Context, userId, sequence int64) (*model.Objective, error)
	ListObjectives(ctx context.Context, userId int64) ([]*model.Objective, error)
	CreateObjective(ctx context.Context, objective *model.Objective) (int64, error)
	UpdateObjective(ctx context.Context, objective *model.Objective) error
	DeleteObjective(ctx context.Context, id int64) error
	DeleteObjectiveBySequence(ctx context.Context, userId, sequence int64) error
	GetKeyResultByID(ctx context.Context, id int64) (*model.KeyResult, error)
	GetKeyResultBySequence(ctx context.Context, userId, sequence int64) (*model.KeyResult, error)
	ListKeyResults(ctx context.Context, userId int64) ([]*model.KeyResult, error)
	ListKeyResultsById(ctx context.Context, id []int64) ([]*model.KeyResult, error)
	ListKeyResultsByObjectiveId(ctx context.Context, objectiveId int64) ([]*model.KeyResult, error)
	CreateKeyResult(ctx context.Context, keyResult *model.KeyResult) (int64, error)
	UpdateKeyResult(ctx context.Context, keyResult *model.KeyResult) error
	DeleteKeyResult(ctx context.Context, id int64) error
	DeleteKeyResultBySequence(ctx context.Context, userId, sequence int64) error
	AggregateObjectiveValue(ctx context.Context, id int64) error
	AggregateKeyResultValue(ctx context.Context, id int64) error
	CreateKeyResultValue(ctx context.Context, keyResultValue *model.KeyResultValue) (int64, error)
	GetKeyResultValues(ctx context.Context, keyResultId int64) ([]*model.KeyResultValue, error)
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

func (c chatbotMapper) GetObjectiveByID(ctx context.Context, id int64) (*model.Objective, error) {
	return adp.GetObjectiveByID(ctx, id)
}

func (c chatbotMapper) GetObjectiveBySequence(ctx context.Context, userId, sequence int64) (*model.Objective, error) {
	return adp.GetObjectiveBySequence(ctx, userId, sequence)
}

func (c chatbotMapper) ListObjectives(ctx context.Context, userId int64) ([]*model.Objective, error) {
	return adp.ListObjectives(ctx, userId)
}

func (c chatbotMapper) CreateObjective(ctx context.Context, objective *model.Objective) (int64, error) {
	return adp.CreateObjective(ctx, objective)
}

func (c chatbotMapper) UpdateObjective(ctx context.Context, objective *model.Objective) error {
	return adp.UpdateObjective(ctx, objective)
}

func (c chatbotMapper) DeleteObjective(ctx context.Context, id int64) error {
	return adp.DeleteObjective(ctx, id)
}

func (c chatbotMapper) DeleteObjectiveBySequence(ctx context.Context, userId, sequence int64) error {
	return adp.DeleteObjectiveBySequence(ctx, userId, sequence)
}

func (c chatbotMapper) GetKeyResultByID(ctx context.Context, id int64) (*model.KeyResult, error) {
	return adp.GetKeyResultByID(ctx, id)
}

func (c chatbotMapper) GetKeyResultBySequence(ctx context.Context, userId, sequence int64) (*model.KeyResult, error) {
	return adp.GetKeyResultBySequence(ctx, userId, sequence)
}

func (c chatbotMapper) ListKeyResults(ctx context.Context, userId int64) ([]*model.KeyResult, error) {
	return adp.ListKeyResults(ctx, userId)
}

func (c chatbotMapper) ListKeyResultsById(ctx context.Context, id []int64) ([]*model.KeyResult, error) {
	return adp.ListKeyResultsById(ctx, id)
}

func (c chatbotMapper) ListKeyResultsByObjectiveId(ctx context.Context, objectiveId int64) ([]*model.KeyResult, error) {
	return adp.ListKeyResultsByObjectiveId(ctx, objectiveId)
}

func (c chatbotMapper) CreateKeyResult(ctx context.Context, keyResult *model.KeyResult) (int64, error) {
	return adp.CreateKeyResult(ctx, keyResult)
}

func (c chatbotMapper) UpdateKeyResult(ctx context.Context, keyResult *model.KeyResult) error {
	return adp.UpdateKeyResult(ctx, keyResult)
}

func (c chatbotMapper) DeleteKeyResult(ctx context.Context, id int64) error {
	return adp.DeleteKeyResult(ctx, id)
}

func (c chatbotMapper) DeleteKeyResultBySequence(ctx context.Context, userId, sequence int64) error {
	return adp.DeleteKeyResultBySequence(ctx, userId, sequence)
}

func (c chatbotMapper) AggregateObjectiveValue(ctx context.Context, id int64) error {
	return adp.AggregateObjectiveValue(ctx, id)
}

func (c chatbotMapper) AggregateKeyResultValue(ctx context.Context, id int64) error {
	return adp.AggregateKeyResultValue(ctx, id)
}

func (c chatbotMapper) CreateKeyResultValue(ctx context.Context, keyResultValue *model.KeyResultValue) (int64, error) {
	return adp.CreateKeyResultValue(ctx, keyResultValue)
}

func (c chatbotMapper) GetKeyResultValues(ctx context.Context, keyResultId int64) ([]*model.KeyResultValue, error) {
	return adp.GetKeyResultValues(ctx, keyResultId)
}

func init() {
	Store = storeObj{}
	Chatbot = chatbotMapper{}
}
