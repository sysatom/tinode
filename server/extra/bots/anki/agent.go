package anki

import (
	"encoding/json"
	"github.com/tinode/chat/server/extra/ruleset/agent"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/types"
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
		Handler: func(ctx types.Context, content interface{}) types.MsgPayload {
			html := ""
			if m, ok := content.(map[string]interface{}); ok {
				if v, ok := m["html"]; ok {
					if s, ok := v.(string); ok {
						var h string
						_ = json.Unmarshal([]byte(s), &h)
						html = h
					}
				}
			}
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
		Handler: func(ctx types.Context, content interface{}) types.MsgPayload {
			num := int64(0)
			if m, ok := content.(map[string]interface{}); ok {
				if v, ok := m["num"]; ok {
					if n, ok := v.(float64); ok {
						num = int64(n)
					}
				}
			}
			_ = store.Chatbot.DataSet(ctx.AsUser, ctx.Original, "getNumCardsReviewedToday", map[string]interface{}{
				"value": num,
			})
			return nil
		},
	},
}
