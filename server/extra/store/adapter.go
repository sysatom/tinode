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

	// ConfigSet kv set
	ConfigSet(uid types.Uid, topic, key string, value model.JSON) error
	// ConfigGet kv get
	ConfigGet(uid types.Uid, topic, key string) (model.JSON, error)
	// OAuthSet oauth set
	OAuthSet(oauth model.OAuth) error
	// OAuthGet oauth get
	OAuthGet(uid types.Uid, topic, t string) (model.OAuth, error)
	// OAuthGetAvailable oauth get available
	OAuthGetAvailable(t string) ([]model.OAuth, error)
}
