package demo

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/tinode/chat/server/extra/command"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/logs"
	"math/big"
	"strconv"
)

var commandRules = []command.Rule{
	{
		Define: "version",
		Help:   `Chatbot framework version`,
		Handler: func(ctx context.Context, tokens []*command.Token) []types.MsgPayload {

			err := store.Chatbot.ConfigSet(1, "abc", "k", map[string]interface{}{
				"a": "123",
			})
			if err != nil {
				logs.Err.Println(err)
			}
			err = store.Chatbot.ConfigSet(1, "abc", "k", map[string]interface{}{
				"a": "abc",
			})
			if err != nil {
				logs.Err.Println(err)
			}
			v, err := store.Chatbot.ConfigGet(1, "abc", "k")
			if err != nil {
				logs.Err.Println(err)
			}
			fmt.Println(v)
			fmt.Println(v.String("a"))

			return []types.MsgPayload{types.TextMsg{Text: "V1"}}
		},
	},
	{
		Define: "rand [number] [number]",
		Help:   `Generate random numbers`,
		Handler: func(ctx context.Context, tokens []*command.Token) []types.MsgPayload {
			min, _ := tokens[1].Value.Int64()
			max, _ := tokens[2].Value.Int64()

			nBing, err := rand.Int(rand.Reader, big.NewInt(max+1-min))
			if err != nil {
				logs.Err.Println(err)
				return nil
			}
			t := nBing.Int64() + min

			return []types.MsgPayload{types.TextMsg{Text: strconv.FormatInt(t, 10)}}
		},
	},
	{
		Define: "messages",
		Help:   `Demo messages`,
		Handler: func(ctx context.Context, tokens []*command.Token) []types.MsgPayload {
			return []types.MsgPayload{types.TextMsg{Text: "msg1"}, types.TextMsg{Text: "msg2"}}
		},
	},
}
