package server

import (
	"github.com/tinode/chat/server/extra/ruleset/agent"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/store/model"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/logs"
)

const (
	AgentVersion = 1
	StatsAgentID = "stats_agent"
)

var agentRules = []agent.Rule{
	{
		Id: StatsAgentID,
		Handler: func(ctx types.Context, content interface{}) types.MsgPayload {
			var v model.JSON
			err := v.Scan(content)
			if err != nil {
				logs.Err.Println(err)
				return nil
			}
			// alert

			// store
			err = store.Chatbot.DataSet(ctx.AsUser, ctx.Original, "stats", v)
			if err != nil {
				logs.Err.Println(err)
				return nil
			}
			return nil
		},
	},
}
