package pocket

import (
	"errors"
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/pkg/flog"
	"github.com/tinode/chat/server/extra/pkg/parser"
	"github.com/tinode/chat/server/extra/ruleset/command"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/store/model"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/extra/vendors"
	"github.com/tinode/chat/server/extra/vendors/pocket"
	serverTypes "github.com/tinode/chat/server/store/types"
	"gorm.io/gorm"
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
		Define: "oauth",
		Help:   `OAuth`,
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
			// check oauth token
			oauth, err := store.Chatbot.OAuthGet(ctx.AsUser, ctx.Original, Name)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				flog.Error(err)
			}
			if oauth.Token != "" {
				return types.TextMsg{Text: "App is authorized"}
			}

			redirectURI := vendors.RedirectURI(pocket.ID, ctx.AsUser, serverTypes.ParseUserId(ctx.Original))
			provider := pocket.NewPocket(Config.ConsumerKey, "", redirectURI, "")
			_, err = provider.GetCode("")
			if err != nil {
				return types.TextMsg{Text: "get code error"}
			}
			url, err := bots.CreateShortUrl(provider.AuthorizeURL())
			if err != nil {
				return types.TextMsg{Text: "create url error"}
			}
			return types.LinkMsg{Title: "OAuth", Url: url}
		},
	},
	{
		Define: "list",
		Help:   `newest 10`,
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
			oauth, err := store.Chatbot.OAuthGet(ctx.AsUser, ctx.Original, Name)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				flog.Error(err)
			}
			if oauth.Token == "" {
				return types.TextMsg{Text: "App is unauthorized"}
			}

			provider := pocket.NewPocket(Config.ConsumerKey, "", "", oauth.Token)
			items, err := provider.Retrieve(10)
			if err != nil {
				flog.Error(err)
				return types.TextMsg{Text: "retrieve error"}
			}

			var header []string
			var row [][]interface{}
			if len(items.List) > 0 {
				header = []string{"Id", "GivenUrl", "GivenTitle", "WordCount"}
				for _, v := range items.List {
					row = append(row, []interface{}{v.Id, v.GivenUrl, v.GivenTitle, v.WordCount})
				}
			}

			return bots.StorePage(ctx, model.PageTable, "Newest List", types.TableMsg{
				Title:  "Newest List",
				Header: header,
				Row:    row,
			})
		},
	},
}
