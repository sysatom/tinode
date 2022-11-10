package help

import (
	"crypto/rand"
	"fmt"
	"github.com/tinode/chat/server/extra/command"
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
			key, err := generateRandomString(16)
			if err != nil {
				logs.Err.Println("bot command id", err)
				return nil
			}

			uGen := storeTypes.UidGenerator{}
			err = uGen.Init(1, []byte(key))
			if err != nil {
				logs.Err.Println("bot command id", err)
				return nil
			}
			return types.TextMsg{Text: uGen.GetStr()}
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
		Define: "form",
		Help:   `Demo form`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return types.FormMsg{
				ID:    helpFormID,
				Title: "Current Value: 1, add/reduce ?",
				Field: []types.FormField{
					{
						Key:      "action",
						Type:     types.FormFieldButton,
						Required: false,
						Value:    "add",
						Default:  nil,
						Intro:    "Add",
					},
					{
						Key:      "action",
						Type:     types.FormFieldButton,
						Required: false,
						Value:    "reduce",
						Default:  nil,
						Intro:    "Reduce",
					},
				},
			}
		},
	},
}

func generateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}
