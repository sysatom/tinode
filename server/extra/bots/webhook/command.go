package webhook

import (
	"fmt"
	"github.com/tinode/chat/server/extra/ruleset/command"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/logs"
	serverTypes "github.com/tinode/chat/server/store/types"
	"strings"
)

var commandRules = []command.Rule{
	{
		Define: "info",
		Help:   `Bot info`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return nil
		},
	},
	{
		Define: `list`,
		Help:   `List webhook`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			items, err := store.Chatbot.DataList(ctx.AsUser, ctx.Original, "webhook:")
			if err != nil {
				return nil
			}

			topicUid := serverTypes.ParseUserId(ctx.Original)

			m := make(map[string]interface{})
			for _, item := range items {
				id := serverTypes.ParseUid(strings.ReplaceAll(item.Key, "webhook:", ""))
				m[item.Key] = fmt.Sprintf("%s/extra/webhook/%d/%d/%d", types.AppUrl(),
					uint64(ctx.AsUser), uint64(topicUid), uint64(id))
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
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			id := types.Id()
			err := store.Chatbot.DataSet(ctx.AsUser, ctx.Original,
				fmt.Sprintf("webhook:%s", id.String()), map[string]interface{}{
					"value": "",
				})
			if err != nil {
				return nil
			}

			topicUid := serverTypes.ParseUserId(ctx.Original)

			return types.TextMsg{Text: fmt.Sprintf("Webhook: %s/extra/webhook/%d/%d/%d", types.AppUrl(),
				uint64(ctx.AsUser), uint64(topicUid), uint64(id))}
		},
	},
	{
		Define: `del [string]`,
		Help:   `delete webhook`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			flag, _ := tokens[1].Value.String()

			err := store.Chatbot.DataDelete(ctx.AsUser, ctx.Original, fmt.Sprintf("webhook:%s", flag))
			if err != nil {
				logs.Err.Println(err)
				return types.TextMsg{Text: "failed"}
			}

			return types.TextMsg{Text: "ok"}
		},
	},
}
