package channels

import (
	"encoding/json"
	"errors"
)

const ChannelNameSuffix = "_channel"

type Publisher interface {
	// Init initializes the channel.
	Init(jsonconf string) error

	// IsReady сhecks if the channel is initialized.
	IsReady() bool

	// Id return channel id
	Id() string
}

type configType struct {
	Name   string          `json:"name"`
	Config json.RawMessage `json:"config"`
}

var publishers map[string]Publisher

func Register(name string, channel Publisher) {
	if publishers == nil {
		publishers = make(map[string]Publisher)
	}

	if channel == nil {
		panic("Register: channel is nil")
	}
	if _, dup := publishers[name]; dup {
		panic("Register: called twice for channel " + name)
	}
	publishers[name] = channel
}

// Init initializes registered publishers.
func Init(jsconfig string) error {
	var config []configType

	if err := json.Unmarshal([]byte(jsconfig), &config); err != nil {
		return errors.New("failed to parse config: " + err.Error())
	}

	for _, cc := range config {
		if channel := publishers[cc.Name]; channel != nil {
			if err := channel.Init(string(cc.Config)); err != nil {
				return err
			}
		}
	}

	return nil
}

func List() map[string]Publisher {
	return publishers
}
