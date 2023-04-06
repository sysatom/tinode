package rust

import (
	"github.com/tinode/chat/server/extra/ruleset/command"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/extra/vendors/crates"
	"github.com/tinode/chat/server/logs"
)

var commandRules = []command.Rule{
	{
		Define: "info",
		Help:   `Bot info`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return nil
		},
	},
	{
		Define: "crate [string]",
		Help:   `crate info`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			name, _ := tokens[1].Value.String()

			api := crates.NewCrates()
			resp, err := api.Info(name)
			if err != nil {
				logs.Err.Println("bot command crate[number]", err)
				return types.TextMsg{Text: "error create"}
			}
			if resp == nil || resp.Crate.ID == "" {
				return types.TextMsg{Text: "empty create"}
			}

			return types.CrateMsg{
				ID:            resp.Crate.ID,
				Name:          resp.Crate.Name,
				Description:   resp.Crate.Description,
				Documentation: resp.Crate.Documentation,
				Homepage:      resp.Crate.Homepage,
				Repository:    resp.Crate.Repository,
				NewestVersion: resp.Crate.NewestVersion,
				Downloads:     resp.Crate.Downloads,
			}
		},
	},
}
