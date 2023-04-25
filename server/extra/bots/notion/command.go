package notion

import (
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/pkg/parser"
	"github.com/tinode/chat/server/extra/ruleset/command"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/extra/vendors/notion"
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
		Define: "config",
		Help:   `Config`,
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
			c1, _ := store.Chatbot.ConfigGet(ctx.AsUser, ctx.Original, TokenKey)
			tokenValue, _ := c1.String("value")
			c2, _ := store.Chatbot.ConfigGet(ctx.AsUser, ctx.Original, ImportPageIdKey)
			importBlockIdValue, _ := c2.String("value")

			return bots.StoreForm(ctx, types.FormMsg{
				ID:    configFormID,
				Title: "Config",
				Field: []types.FormField{
					{
						Type:        types.FormFieldText,
						Key:         TokenKey,
						Value:       tokenValue,
						ValueType:   types.FormFieldValueString,
						Label:       "Internal Integration Token",
						Placeholder: "Input token",
					},
					{
						Type:        types.FormFieldText,
						Key:         ImportPageIdKey,
						Value:       importBlockIdValue,
						ValueType:   types.FormFieldValueString,
						Label:       "MindCache page id",
						Placeholder: "Input page id",
					},
				},
			})
		},
	},
	{
		Define: "search [string]",
		Help:   "Searches all original pages, databases, and child pages/databases that are shared with the integration.",
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
			query, _ := tokens[1].Value.String()

			// token value
			j, err := store.Chatbot.ConfigGet(ctx.AsUser, ctx.Original, "token")
			if err != nil {
				return nil
			}
			token, _ := j.String("value")
			if token == "" {
				return types.TextMsg{Text: "set config"}
			}

			provider := notion.NewNotion(token)
			pages, err := provider.Search(query)
			if err != nil {
				return types.TextMsg{Text: "search error"}
			}
			if len(pages) == 0 {
				return types.TextMsg{Text: "Empty"}
			}
			var links types.LinkListMsg
			for _, page := range pages {
				links.Links = append(links.Links, types.LinkMsg{Title: page.Object, Url: page.URL})
			}
			return links
		},
	},
	{
		Define: "import [string]",
		Help:   "Append to MindCache page",
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
			text, _ := tokens[1].Value.String()

			// token value
			j, err := store.Chatbot.ConfigGet(ctx.AsUser, ctx.Original, TokenKey)
			if err != nil {
				return nil
			}
			token, _ := j.String("value")
			if token == "" {
				return types.TextMsg{Text: "set config"}
			}

			// block id
			j2, err := store.Chatbot.ConfigGet(ctx.AsUser, ctx.Original, ImportPageIdKey)
			if err != nil {
				return nil
			}
			pageId, _ := j2.String("value")
			if pageId == "" {
				return types.TextMsg{Text: "set config"}
			}

			provider := notion.NewNotion(token)
			err = provider.AppendBlock(pageId, text)
			if err != nil {
				logs.Err.Println(err)
				return types.TextMsg{Text: "import error"}
			}

			return types.TextMsg{Text: "ok"}
		},
	},
}
