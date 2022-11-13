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

	// Chatbot

	// DataSet data set
	DataSet(uid types.Uid, topic, key string, value model.JSON) error
	// DataGet data get
	DataGet(uid types.Uid, topic, key string) (model.JSON, error)
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
	// PageSet page set
	PageSet(pageId string, page model.Page) error
	// PageGet page get
	PageGet(pageId string) (model.Page, error)

	GetObjectiveByID(id int64) (*model.Objective, error)
	GetObjectiveBySequence(userId, sequence int64) (*model.Objective, error)
	ListObjectives(userId int64) ([]*model.Objective, error)
	CreateObjective(objective *model.Objective) (int64, error)
	UpdateObjective(objective *model.Objective) error
	DeleteObjective(id int64) error
	DeleteObjectiveBySequence(userId, sequence int64) error
	GetKeyResultByID(id int64) (*model.KeyResult, error)
	GetKeyResultBySequence(userId, sequence int64) (*model.KeyResult, error)
	ListKeyResults(userId int64) ([]*model.KeyResult, error)
	ListKeyResultsById(id []int64) ([]*model.KeyResult, error)
	ListKeyResultsByObjectiveId(objectiveId int64) ([]*model.KeyResult, error)
	CreateKeyResult(keyResult *model.KeyResult) (int64, error)
	UpdateKeyResult(keyResult *model.KeyResult) error
	DeleteKeyResult(id int64) error
	DeleteKeyResultBySequence(userId, sequence int64) error
	AggregateObjectiveValue(id int64) error
	AggregateKeyResultValue(id int64) error
	CreateKeyResultValue(keyResultValue *model.KeyResultValue) (int64, error)
	GetKeyResultValues(keyResultId int64) ([]*model.KeyResultValue, error)
}
