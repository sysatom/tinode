package webhook

import (
	"encoding/json"
	"errors"
	"github.com/emicklei/go-restful/v3"
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/pkg/route"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/logs"
)

const Name = "webhook"

var handler bot

func init() {
	bots.Register(Name, &handler)
}

type bot struct {
	initialized bool
	bots.Base
}

type configType struct {
	Enabled bool `json:"enabled"`
}

func (bot) Init(jsonconf json.RawMessage) error {

	// Check if the handler is already initialized
	if handler.initialized {
		return errors.New("already initialized")
	}

	var config configType
	if err := json.Unmarshal(jsonconf, &config); err != nil {
		return errors.New("failed to parse config: " + err.Error())
	}

	if !config.Enabled {
		logs.Info.Printf("bot %s disabled", Name)
		return nil
	}

	handler.initialized = true

	return nil
}

func (bot) IsReady() bool {
	return handler.initialized
}

func (bot) WebService() *restful.WebService {
	return route.WebService(
		Name, serviceVersion,
		route.Route("POST", "/webhook/{flag}", webhook, "trigger webhook", route.WithParam(&route.Param{
			Type:        route.PathParamType,
			Name:        "flag",
			Description: "flag param",
			DataType:    "string",
		})),
	)
}

func (b bot) Command(ctx types.Context, content interface{}) (types.MsgPayload, error) {
	return bots.RunCommand(commandRules, ctx, content)
}
