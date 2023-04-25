package web

import (
	"github.com/tinode/chat/server/extra/pkg/parser"
	"github.com/tinode/chat/server/extra/ruleset/command"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/extra/vendors"
	"github.com/tinode/chat/server/extra/vendors/oneai"
	"github.com/tinode/chat/server/logs"
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
		Define: "summary [string]",
		Help:   `web page summary`,
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
			url, _ := tokens[1].Value.String()
			// api key
			val, err := vendors.GetConfig(oneai.ID, oneai.ApiKey)
			if err != nil {
				return types.TextMsg{Text: "error api config"}
			}

			api := oneai.NewOneAI(val.String())
			resp, err := api.Summarize(url)
			if err != nil {
				logs.Err.Println(err)
				return types.TextMsg{Text: "error summarize"}
			}

			if len(resp.Output) != 2 || len(resp.Output[1].Contents) == 0 {
				return types.TextMsg{Text: "empty summarize"}
			}

			return types.TextMsg{Text: resp.Output[1].Contents[0].Utterance}
		},
	},
}
