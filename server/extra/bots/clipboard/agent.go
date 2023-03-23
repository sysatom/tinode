package clipboard

import (
	"github.com/tinode/chat/server/extra/ruleset/agent"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/extra/utils"
)

const (
	AgentVersion  = 1
	UploadAgentID = "clipboard_upload"
)

var agentRules = []agent.Rule{
	{
		Id:   UploadAgentID,
		Help: "update clipboard",
		Args: []string{"txt"},
		Handler: func(ctx types.Context, content interface{}) types.MsgPayload {
			j, err := utils.ConvertJSON(content)
			if err != nil {
				return nil
			}
			txt, ok := j.String("txt")
			if !ok {
				return nil
			}
			return types.TextMsg{Text: txt}
		},
	},
}
