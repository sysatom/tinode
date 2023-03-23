package anki

import (
	"github.com/tinode/chat/server/extra/ruleset/agent"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/extra/utils"
)

const (
	AgentVersion  = 1
	StatsAgentID  = "stats_agent"
	ReviewAgentID = "review_agent"
)

var agentRules = []agent.Rule{
	{
		Id:   StatsAgentID,
		Help: "import anki stats",
		Args: []string{"html"},
		Handler: func(ctx types.Context, content interface{}) types.MsgPayload {
			j, err := utils.ConvertJSON(content)
			if err != nil {
				return nil
			}
			html, _ := j.String("html")
			if html == "" {
				return nil
			}
			_ = store.Chatbot.DataSet(ctx.AsUser, ctx.Original, "getCollectionStatsHTML", map[string]interface{}{
				"value": html,
			})
			return nil
		},
	},
	{
		Id:   ReviewAgentID,
		Help: "import anki review count",
		Args: []string{"num"},
		Handler: func(ctx types.Context, content interface{}) types.MsgPayload {
			j, err := utils.ConvertJSON(content)
			if err != nil {
				return nil
			}
			num, _ := j.Int64("num")
			_ = store.Chatbot.DataSet(ctx.AsUser, ctx.Original, "getNumCardsReviewedToday", map[string]interface{}{
				"value": num,
			})
			return nil
		},
	},
}
