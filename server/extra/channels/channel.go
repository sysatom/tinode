package channels

import (
	"encoding/json"
	"errors"
	"github.com/tinode/chat/server/extra/channels/crawler"
	"gopkg.in/yaml.v2"
	"io/fs"
	"os"
	"path/filepath"
)

const ChannelNameSuffix = "_channel"

type Publisher *crawler.Rule

type configType struct {
	Path string `json:"path"`
}

var publishers map[string]Publisher

// Init initializes registered publishers.
func Init(jsconfig json.RawMessage) error {
	var config configType

	if err := json.Unmarshal(jsconfig, &config); err != nil {
		return errors.New("failed to parse config: " + err.Error())
	}

	if publishers == nil {
		publishers = make(map[string]Publisher)
	}

	return filepath.Walk(config.Path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if ext := filepath.Ext(path); ext != ".yml" {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		var r *crawler.Rule
		err = yaml.Unmarshal(data, &r)
		if err != nil {
			return err
		}

		publishers[r.Name] = r

		return nil
	})
}

func List() map[string]Publisher {
	return publishers
}
