package pocket

import (
	"errors"
	"github.com/tinode/chat/server/extra/ruleset/cron"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/extra/vendors/pocket"
	"github.com/tinode/chat/server/logs"
	"gorm.io/gorm"
)

var cronRules = []cron.Rule{
	{
		Name: "pocket_add",
		When: "* * * * *",
		Action: func(ctx types.Context) []types.MsgPayload {
			oauth, err := store.Chatbot.OAuthGet(ctx.AsUser, ctx.Original, Name)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				logs.Err.Println("bot command pocket oauth", err)
			}
			if oauth.Token == "" {
				return nil
			}

			provider := pocket.NewPocket(Config.ConsumerKey, "", "", oauth.Token)
			items, err := provider.Retrieve(10)
			if err != nil {
				logs.Err.Println(err)
				return nil
			}

			var r []types.MsgPayload
			for _, item := range items.List {
				r = append(r, types.LinkMsg{
					Title: item.GivenTitle,
					Url:   item.GivenUrl,
				})
			}
			return r
		},
	},
}
