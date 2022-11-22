package help

import (
	"crypto/rand"
	"fmt"
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/ruleset/command"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/logs"
	serverTypes "github.com/tinode/chat/server/store/types"
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
			return types.TextMsg{Text: types.Id().String()}
		},
	},
	{
		Define: "uid [string]",
		Help:   `Decode UID string`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			str, _ := tokens[1].Value.String()
			var uid serverTypes.Uid
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
						Key:         "text",
						Type:        types.FormFieldText,
						ValueType:   types.FormFieldValueString,
						Value:       "",
						Label:       "Text",
						Placeholder: "Input text",
					},
					{
						Key:         "password",
						Type:        types.FormFieldPassword,
						ValueType:   types.FormFieldValueString,
						Value:       "",
						Label:       "Password",
						Placeholder: "Input password",
						Required:    true,
					},
					{
						Key:         "number",
						Type:        types.FormFieldNumber,
						ValueType:   types.FormFieldValueInt64,
						Value:       "",
						Label:       "Number",
						Placeholder: "Input number",
					},
					{
						Key:         "bool",
						Type:        types.FormFieldRadio,
						ValueType:   types.FormFieldValueBool,
						Value:       "",
						Label:       "Bool",
						Placeholder: "Switch",
						Option:      []string{"true", "false"},
					},
					{
						Key:         "multi",
						Type:        types.FormFieldCheckbox,
						ValueType:   types.FormFieldValueStringSlice,
						Value:       "",
						Label:       "Multiple",
						Placeholder: "Select multiple",
						Option:      []string{"a", "b", "c"},
					},
					{
						Key:         "textarea",
						Type:        types.FormFieldTextarea,
						ValueType:   types.FormFieldValueString,
						Value:       "",
						Label:       "Textarea",
						Placeholder: "Input textarea",
					},
					{
						Key:         "select",
						Type:        types.FormFieldSelect,
						ValueType:   types.FormFieldValueFloat64,
						Value:       "",
						Label:       "Select",
						Placeholder: "Select float",
						Option:      []string{"1.01", "2.02", "3.03"},
					},
					{
						Key:         "range",
						Type:        types.FormFieldRange,
						ValueType:   types.FormFieldValueInt64,
						Value:       "",
						Label:       "Range",
						Placeholder: "range value",
					},
				},
			})
		},
	},
	{
		Define: "echo [any]",
		Help:   "print",
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			val := tokens[1].Value.Source
			return types.TextMsg{Text: fmt.Sprintf("%v", val)}
		},
	},
	{
		Define: "agent",
		Help:   `agent url`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return bots.AgentURI(ctx)
		},
	},
}
