package cloudflare

import (
	"fmt"
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/pkg/parser"
	"github.com/tinode/chat/server/extra/ruleset/command"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/extra/vendors/cloudflare"
	"github.com/tinode/chat/server/logs"
	"time"
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
			c2, _ := store.Chatbot.ConfigGet(ctx.AsUser, ctx.Original, ZoneIdKey)
			zoneIdValue, _ := c2.String("value")
			c3, _ := store.Chatbot.ConfigGet(ctx.AsUser, ctx.Original, AccountIdKey)
			accountIdValue, _ := c3.String("value")

			return bots.StoreForm(ctx, types.FormMsg{
				ID:    configFormID,
				Title: "Config",
				Field: []types.FormField{
					{
						Type:        types.FormFieldText,
						Key:         "token",
						Value:       tokenValue,
						ValueType:   types.FormFieldValueString,
						Label:       "Token",
						Placeholder: "Input token",
					},
					{
						Type:        types.FormFieldText,
						Key:         "zone_id",
						Value:       zoneIdValue,
						ValueType:   types.FormFieldValueString,
						Label:       "Zone Id",
						Placeholder: "Input zone id",
					},
					{
						Type:        types.FormFieldText,
						Key:         "account_id",
						Value:       accountIdValue,
						ValueType:   types.FormFieldValueString,
						Label:       "Account Id",
						Placeholder: "Input account id",
					},
				},
			})
		},
	},
	{
		Define: "test",
		Help:   "Test",
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
			c1, _ := store.Chatbot.ConfigGet(ctx.AsUser, ctx.Original, TokenKey)
			tokenValue, _ := c1.String("value")
			c2, _ := store.Chatbot.ConfigGet(ctx.AsUser, ctx.Original, ZoneIdKey)
			zoneIdValue, _ := c2.String("value")

			if tokenValue == "" || zoneIdValue == "" {
				return types.TextMsg{Text: "config error"}
			}

			now := time.Now()
			startDate := now.Add(-24 * time.Hour).Format(time.RFC3339)
			endDate := now.Format(time.RFC3339)

			provider := cloudflare.NewCloudflare(tokenValue, zoneIdValue)
			resp, err := provider.GetAnalytics(startDate, endDate)
			if err != nil {
				logs.Err.Println(err)
				return nil
			}
			fmt.Println(resp)

			return nil
		},
	},
}
