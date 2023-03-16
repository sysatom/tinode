package help

import (
	"github.com/tinode/chat/server/extra/ruleset/agent"
	"github.com/tinode/chat/server/extra/types"
)

const (
	AgentVersion  = 1
	ImportAgentID = "import_agent"
)

var agentRules = []agent.Rule{
	{
		Id:   ImportAgentID,
		Help: "agent example",
		Handler: func(ctx types.Context, content interface{}) types.MsgPayload {
			return types.TextMsg{Text: "imported"}
		},
	},
}
