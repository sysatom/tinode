package genshin

import (
	"fmt"
	"github.com/tinode/chat/server/extra/pkg/parser"
	"github.com/tinode/chat/server/extra/ruleset/command"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/types"
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
		Define: "uid",
		Help:   `get uid`,
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
			// get
			v, err := store.Chatbot.ConfigGet(ctx.AsUser, ctx.Original, "uid")
			if err != nil {
				logs.Err.Println("bot command key", err)
			}
			uid, _ := v.Float64("value")

			return types.TextMsg{Text: fmt.Sprintf("uid: %d", int64(uid))}
		},
	},
	{
		Define: "uid [number]",
		Help:   `set uid`,
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
			uid, _ := tokens[1].Value.Int64()

			// get
			v, err := store.Chatbot.ConfigGet(ctx.AsUser, ctx.Original, "uid")
			if err != nil {
				logs.Err.Println("bot command uid [number]", err)
			}
			old, _ := v.Int64("value")

			// set
			err = store.Chatbot.ConfigSet(ctx.AsUser, ctx.Original, "uid", types.KV{
				"value": uid,
			})
			if err != nil {
				logs.Err.Println("bot command uid [number]", err)
			}

			return types.TextMsg{Text: fmt.Sprintf("%d --> %d", old, uid)}
		},
	},
	{
		Define: "profile",
		Help:   `genshin profile`,
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
			// get
			v, err := store.Chatbot.ConfigGet(ctx.AsUser, ctx.Original, "uid")
			if err != nil {
				logs.Err.Println("bot command key", err)
			}
			uid, _ := v.Float64("value")

			return types.LinkMsg{Url: fmt.Sprintf("https://enka.network/u/%d", int64(uid))}
		},
	},
}
