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
		Define: "version",
		Help:   `Version`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return types.TextMsg{Text: "V1"}
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
				m[item.Key] = fmt.Sprintf("http://127.0.0.1:6060/extra/webhook/%d/%d/%d",
					uint64(ctx.AsUser), uint64(topicUid), uint64(id)) // todo
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

			return types.TextMsg{Text: fmt.Sprintf("Webhook: http://127.0.0.1:6060/extra/webhook/%d/%d/%d",
				uint64(ctx.AsUser), uint64(topicUid), uint64(id))} // todo
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
