package help

import (
	"crypto/rand"
	"fmt"
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/ruleset/command"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/logs"
	storeTypes "github.com/tinode/chat/server/store/types"
	"math/big"
	"strconv"
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
		Define: "rand [number] [number]",
		Help:   `Generate random numbers`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			min, _ := tokens[1].Value.Int64()
			max, _ := tokens[2].Value.Int64()

			nBing, err := rand.Int(rand.Reader, big.NewInt(max+1-min))
			if err != nil {
				logs.Err.Println("bot command rand [number] [number]", err)
				return nil
			}
			t := nBing.Int64() + min

			return types.TextMsg{Text: strconv.FormatInt(t, 10)}
		},
	},
	{
		Define: "id",
		Help:   `Generate random id`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return types.TextMsg{Text: types.Id()}
		},
	},
	{
		Define: "uid [string]",
		Help:   `Decode UID string`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			str, _ := tokens[1].Value.String()
			var uid storeTypes.Uid
			var result string
			err := uid.UnmarshalText([]byte(str))
			if err != nil {
				result = err.Error()
			} else {
				result = fmt.Sprintf("%d", uid)
			}

			return types.TextMsg{Text: result}
		},
	},
	{
		Define: "messages",
		Help:   `Demo messages`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return types.TextMsg{Text: "msg1"}
		},
	},
	{
		Define: `qr [string]`,
		Help:   `Generate QR code`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return types.TextMsg{Text: "msg1"}
		},
	},
	{
		Define: `pinyin [string]`,
		Help:   "chinese pinyin conversion",
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return types.TextMsg{Text: "msg1"}
		},
	},
	{
		Define: "form",
		Help:   `Demo form`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return bots.StoreForm(ctx, types.FormMsg{
				ID:    helpFormID,
				Title: "Current Value: 1, add/reduce ?",
				Field: []types.FormField{
					{
						Key:         "action",
						Type:        types.FormFieldText,
						ValueType:   types.FormFieldValueString,
						Value:       "add",
						Label:       "Add",
						Placeholder: "Add",
					},
					{
						Key:         "action",
						Type:        types.FormFieldText,
						ValueType:   types.FormFieldValueString,
						Value:       "reduce",
						Label:       "Reduce",
						Placeholder: "Reduce",
					},
				},
			})
		},
	},
}
