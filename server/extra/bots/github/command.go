package github

import (
	"errors"
	"github.com/tinode/chat/server/extra/command"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/extra/vendors"
	"github.com/tinode/chat/server/extra/vendors/github"
	"github.com/tinode/chat/server/logs"
	"gorm.io/gorm"
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
			// check oauth token
			oauth, err := store.Chatbot.OAuthGet(ctx.AsUser, ctx.Original, Name)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				logs.Err.Println(err)
			}
			if oauth.Token != "" {
				return []types.MsgPayload{types.TextMsg{Text: "App is authorized"}}
			}

			redirectURI := vendors.RedirectURI(github.ID)
			provider := github.NewGithub(Config.ID, Config.Secret, redirectURI, "")
			return []types.MsgPayload{types.LinkMsg{Url: provider.AuthorizeURL()}}
		},
	},
}
