package github

import (
	"errors"
	"fmt"
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/ruleset/command"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/store/model"
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
		Define: "config",
		Help:   `Config`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			j, _ := store.Chatbot.ConfigGet(ctx.AsUser, ctx.Original, RepoKey)
			repoValue, _ := j.String("value")

			return bots.StoreForm(ctx, types.FormMsg{
				ID:    configFormID,
				Title: "Config",
				Field: []types.FormField{
					{
						Type:        types.FormFieldText,
						Key:         "repo",
						Value:       repoValue,
						ValueType:   types.FormFieldValueString,
						Label:       "Repo",
						Placeholder: "Input repo",
					},
				},
			})
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
			return types.LinkMsg{Title: "OAuth URL", Url: provider.AuthorizeURL()}
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

			return bots.StorePage(ctx, model.PageTable, "User", table)
		},
	},
	{
		Define: "issue [string]",
		Help:   `create issue`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			text, _ := tokens[1].Value.String()

			oauth, err := store.Chatbot.OAuthGet(ctx.AsUser, ctx.Original, github.ID)
			if err != nil {
				return nil
			}
			if oauth.Token == "" {
				return types.TextMsg{Text: "oauth error"}
			}

			// get user
			client := github.NewGithub("", "", "", oauth.Token)
			user, err := client.GetUser()
			if err != nil {
				return nil
			}
			if *user.Login == "" {
				return nil
			}

			// repo value
			j, err := store.Chatbot.ConfigGet(ctx.AsUser, ctx.Original, RepoKey)
			repo, _ := j.String("value")
			if repo == "" {
				return types.TextMsg{Text: "set repo [string]"}
			}

			// create issue
			issue, err := client.CreateIssue(*user.Login, repo, github.Issue{Title: &text})
			if err != nil {
				logs.Err.Println(err)
				return nil
			}
			if *issue.ID == 0 {
				return nil
			}

			return types.LinkMsg{
				Title: fmt.Sprintf("Issue #%d", *issue.Number),
				Url:   *issue.HTMLURL,
			}
		},
	},
	{
		Define: "card [string]",
		Help:   `create project card`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			text, _ := tokens[1].Value.String()

			oauth, err := store.Chatbot.OAuthGet(ctx.AsUser, ctx.Original, github.ID)
			if err != nil {
				return nil
			}
			if oauth.Token == "" {
				return types.TextMsg{Text: "oauth error"}
			}

			// get user
			client := github.NewGithub("", "", "", oauth.Token)
			user, err := client.GetUser()
			if err != nil {
				return nil
			}
			if *user.Login == "" {
				return nil
			}

			// get projects
			projects, err := client.GetUserProjects(*user.Login)
			if err != nil {
				logs.Err.Println(err)
				return nil
			}
			if len(*projects) == 0 {
				return nil
			}

			// get columns
			columns, err := client.GetProjectColumns(*(*projects)[0].ID)
			if err != nil {
				logs.Err.Println(err)
				return nil
			}
			if len(*columns) == 0 {
				return nil
			}

			// create card
			card, err := client.CreateCard(*(*columns)[0].ID, github.ProjectCard{Note: &text})
			if err != nil {
				logs.Err.Println(err)
				return nil
			}
			if *card.ID == 0 {
				return nil
			}

			return types.TextMsg{Text: fmt.Sprintf("Created Project Card #%d", *card.ID)}
		},
	},
}
