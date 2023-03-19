package clipboard

import (
	"github.com/tinode/chat/server/extra/ruleset/agent"
	"github.com/tinode/chat/server/extra/types"
)

const (
	AgentVersion  = 1
	UploadAgentID = "clipboard_upload"
)

var agentRules = []agent.Rule{
	{
		Id: UploadAgentID,
		Handler: func(ctx types.Context, content interface{}) types.MsgPayload {
			if txt, ok := content.(string); ok {
				return types.TextMsg{Text: txt}
			}
			return nil
		},
	},
}
