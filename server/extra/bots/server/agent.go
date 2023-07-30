package server

import (
	"github.com/tinode/chat/server/extra/ruleset/agent"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/logs"
)

const (
	AgentVersion = 1
	StatsAgentID = "stats_agent"
)

var agentRules = []agent.Rule{
	{
		Id:   StatsAgentID,
		Help: "upload server status",
		Args: []string{"cpu", "memory", "info"},
		Handler: func(ctx types.Context, content types.KV) types.MsgPayload {
			j := types.KV{}
			err := j.Scan(content)
			if err != nil {
				return nil
			}
			// alert

			// store
			err = store.Chatbot.DataSet(ctx.AsUser, ctx.Original, "stats", j)
			if err != nil {
				logs.Err.Println(err)
				return nil
			}
			return nil
		},
	},
}
