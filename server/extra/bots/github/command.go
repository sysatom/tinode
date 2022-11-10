package github

import (
	"errors"
	"github.com/tinode/chat/server/extra/ruleset/command"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/extra/vendors"
	"github.com/tinode/chat/server/extra/vendors/github"
	"github.com/tinode/chat/server/logs"
	serverTypes "github.com/tinode/chat/server/store/types"
	"gorm.io/gorm"
)

var commandRules = []command.Rule{
	{
		Define: "version",
		Help:   `Version`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return types.TextMsg{Text: "V1"}
		},
	},
	{
		Define: "oauth",
		Help:   `OAuth`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			// check oauth token
			oauth, err := store.Chatbot.OAuthGet(ctx.AsUser, ctx.Original, Name)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				logs.Err.Println("bot command github oauth", err)
			}
			if oauth.Token != "" {
				return types.TextMsg{Text: "App is authorized"}
			}

			redirectURI := vendors.RedirectURI(github.ID, ctx.AsUser, serverTypes.ParseUserId(ctx.Original))
			provider := github.NewGithub(Config.ID, Config.Secret, redirectURI, "")
			return types.LinkMsg{Url: provider.AuthorizeURL()}
		},
	},
	{
		Define: "user",
		Help:   `Get user info`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			// get token
			oauth, err := store.Chatbot.OAuthGet(ctx.AsUser, ctx.Original, Name)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				logs.Err.Println("bot command github user", err)
			}
			if oauth.Token == "" {
				return types.TextMsg{Text: "App is unauthorized"}
			}

			provider := github.NewGithub("", "", "", oauth.Token)

			user, err := provider.GetUser()
			if err != nil {
				return types.TextMsg{Text: err.Error()}
			}
			if user == nil {
				return types.TextMsg{Text: "user error"}
			}
			table := types.TableMsg{}
			table.Title = "User"
			table.Header = []string{
				"Login",
				"Followers",
				"Following",
				"URL",
			}
			table.Row = append(table.Row, []interface{}{
				*user.Login,
				*user.Followers,
				*user.Following,
				*user.HTMLURL,
			})

			return table
		},
	},
}
