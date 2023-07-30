package dev

import (
	"github.com/tinode/chat/server/extra/ruleset/agent"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/logs"
)

const (
	AgentVersion  = 1
	ImportAgentID = "import_agent"
)

var agentRules = []agent.Rule{
	{
		Id:   ImportAgentID,
		Help: "agent example",
		Args: []string{},
		Handler: func(ctx types.Context, content types.KV) types.MsgPayload {
			logs.Info.Println(content)
			return nil
		},
	},
}
