package github

import (
	"github.com/tinode/chat/server/extra/command"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/extra/vendors"
	"github.com/tinode/chat/server/extra/vendors/github"
)

var commandRules = []command.Rule{
	{
		Define: "version",
		Help:   `Version`,
		Handler: func(ctx types.Context, tokens []*command.Token) []types.MsgPayload {
			return []types.MsgPayload{types.TextMsg{Text: "V1"}}
		},
	},
	{
		Define: "oauth",
		Help:   `OAuth`,
		Handler: func(ctx types.Context, tokens []*command.Token) []types.MsgPayload {
			redirectURI := vendors.RedirectURI(github.ID)
			provider := github.NewGithub(config.ID, config.Secret, redirectURI, "")
			return []types.MsgPayload{types.LinkMsg{Url: provider.AuthorizeURL()},
				types.TextMsg{Text: ctx.AsUser.UserId()},
				types.TextMsg{Text: ctx.RcptTo},
				types.TextMsg{Text: ctx.Original},
			}
		},
	},
}
