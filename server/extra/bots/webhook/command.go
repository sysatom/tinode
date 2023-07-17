package webhook

import (
	"fmt"
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/pkg/parser"
	"github.com/tinode/chat/server/extra/pkg/route"
	"github.com/tinode/chat/server/extra/ruleset/command"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/logs"
	serverTypes "github.com/tinode/chat/server/store/types"
	"strings"
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
		Define: `list`,
		Help:   `List webhook`,
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
			prefix := "webhook:"
			items, err := store.Chatbot.DataList(ctx.AsUser, ctx.Original, types.DataFilter{Prefix: &prefix})
			if err != nil {
				return nil
			}

			m := make(map[string]interface{})
			for _, item := range items {
				flag := serverTypes.ParseUid(strings.ReplaceAll(item.Key, "webhook:", ""))
				m[item.Key] = route.URL(Name, serviceVersion, fmt.Sprintf("webhook/%s", flag))
			}

			return types.InfoMsg{
				Title: "Webhook list",
				Model: m,
			}
		},
	},
	{
		Define: `create`,
		Help:   `create webhook`,
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
			p := types.KV{}
			p["uid"] = ctx.AsUser.UserId()
			flag, err := bots.StoreParameter(p, time.Now().Add(24*365*time.Hour))
			if err != nil {
				return types.TextMsg{Text: "error parameter"}
			}

			err = store.Chatbot.DataSet(ctx.AsUser, ctx.Original,
				fmt.Sprintf("webhook:%s", flag), map[string]interface{}{
					"value": "",
				})
			if err != nil {
				return types.TextMsg{Text: "error create"}
			}

			return types.TextMsg{Text: fmt.Sprintf("Webhook: %s", route.URL(Name, serviceVersion, fmt.Sprintf("webhook/%s", flag)))}
		},
	},
	{
		Define: `del [string]`,
		Help:   `delete webhook`,
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
			flag, _ := tokens[1].Value.String()

			err := store.Chatbot.ParameterDelete(flag)
			if err != nil {
				logs.Err.Println(err)
				return types.TextMsg{Text: "failed"}
			}

			err = store.Chatbot.DataDelete(ctx.AsUser, ctx.Original, fmt.Sprintf("webhook:%s", flag))
			if err != nil {
				logs.Err.Println(err)
				return types.TextMsg{Text: "failed"}
			}

			return types.TextMsg{Text: "ok"}
		},
	},
}
