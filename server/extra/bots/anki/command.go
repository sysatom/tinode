package anki

import (
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/pkg/parser"
	"github.com/tinode/chat/server/extra/ruleset/command"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/store/model"
	"github.com/tinode/chat/server/extra/types"
)

var commandRules = []command.Rule{
	{
		Define: "info",
		Help:   `Bot info`,
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
			return nil
		},
	},
	{
		Define: "agent",
		Help:   `agent url`,
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
			return bots.AgentURI(ctx)
		},
	},
	{
		Define: "stats",
		Help:   `Anki collection statistics`,
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
			j, err := store.Chatbot.DataGet(ctx.AsUser, ctx.Original, "getCollectionStatsHTML")
			if err != nil {
				return types.TextMsg{Text: "Empty"}
			}
			html, ok := j.String("value")
			if !ok {
				return types.TextMsg{Text: "Empty"}
			}
			return bots.StorePage(ctx, model.PageHtml, "Anki collection statistics", types.HtmlMsg{Raw: html})
		},
	},
}
