package github

import (
	"github.com/tinode/chat/server/extra/ruleset/cron"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/extra/vendors/github"
	"github.com/tinode/chat/server/logs"
)

var cronRules = []cron.Rule{
	{
		Name: "github_starred",
		When: "* * * * *",
		Action: func(ctx types.Context) []types.MsgPayload {
			// data
			client := github.NewGithub("", "", "", ctx.Token)
			user, err := client.GetUser()
			if err != nil {
				logs.Err.Println("cron github_starred", err)
				return []types.MsgPayload{}
			}
			if *user.Login == "" {
				return []types.MsgPayload{}
			}

			repos, err := client.GetStarred(*user.Login)
			if err != nil {
				logs.Err.Println("cron github_starred", err)
				return []types.MsgPayload{}
			}
			reposList := *repos
			var r []types.MsgPayload
			for i := range reposList {
				item := reposList[i]
				r = append(r, types.InfoMsg{
					Title: *item.FullName,
					Model: map[string]interface{}{
						"Owner":      *item.Owner.Login,
						"Name":       *item.Name,
						"Url":        *item.HTMLURL,
						"Stargazers": *item.StargazersCount,
						"Forks":      *item.ForksCount,
						"Watchers":   *item.WatchersCount,
					},
				})
			}
			return r
		},
	},
	{
		Name: "github_stargazers",
		When: "* * * * *",
		Action: func(types.Context) []types.MsgPayload {
			return nil
		},
	},
}
