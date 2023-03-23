package dev

import (
	"bytes"
	"crypto/rand"
	_ "embed"
	"fmt"
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/pkg/queue"
	"github.com/tinode/chat/server/extra/ruleset/command"
	"github.com/tinode/chat/server/extra/store/model"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/logs"
	serverTypes "github.com/tinode/chat/server/store/types"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
	"math/big"
	"strconv"
	"strings"
	"time"
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
		Define: "ts [number]",
		Help:   `timestamp format`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			num, _ := tokens[1].Value.Int64()
			t := time.Unix(num, 0)
			return types.TextMsg{Text: t.Format(time.RFC3339)}
		},
	},
	{
		Define: "messages",
		Help:   `[example] messages`,
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
		Help:   `[example] form`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return bots.StoreForm(ctx, types.FormMsg{
				ID:    devFormID,
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
		Define: "action",
		Help:   "[example] action",
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return types.ActionMsg{
				ID:     devActionID,
				Title:  "Operate ... ?",
				Option: []string{"do1", "do2"},
			}
		},
	},
	{
		Define: "guess",
		Help:   "Guess number game",
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			// rand number
			big, err := rand.Int(rand.Reader, big.NewInt(1000))

			var initValue model.JSON
			initValue = map[string]interface{}{"number": big.Int64()}
			ctx.SessionRuleId = guessSessionID
			err = bots.SessionStart(ctx, initValue)
			if err != nil {
				return types.TextMsg{Text: "session error"}
			}
			return types.TextMsg{Text: "Input a number?"}
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
	{
		Define: "plot",
		Help:   `[example] plot graph`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			p := plot.New()

			p.Title.Text = "Plotutil example"
			p.X.Label.Text = "X"
			p.Y.Label.Text = "Y"

			err := plotutil.AddLinePoints(p,
				"First", randomPoints(15),
				"Second", randomPoints(15),
				"Third", randomPoints(15))
			if err != nil {
				panic(err)
			}

			w := bytes.NewBufferString("")

			c := vgimg.New(vg.Points(500), vg.Points(500))
			dc := draw.New(c)
			p.Draw(dc)

			png := vgimg.PngCanvas{Canvas: c}
			if _, err := png.WriteTo(w); err != nil {
				panic(err)
			}

			return types.ImageConvert(w.Bytes(), "Plot", 500, 500)
		},
	},
	{
		Define: "queue",
		Help:   `[example] publish queue`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			err := queue.AsyncMessage(ctx.RcptTo, ctx.Original, types.TextMsg{Text: time.Now().String()})
			if err != nil {
				return types.TextMsg{Text: err.Error()}
			}
			return types.TextMsg{Text: "ok"}
		},
	},
	{
		Define: "instruct",
		Help:   `[example] create instruct`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			data := model.JSON{}
			data["txt"] = "example"
			return bots.StoreInstruct(ctx, types.InstructMsg{
				No:       types.Id().String(),
				Object:   model.InstructObjectHelper,
				Bot:      Name,
				Flag:     ExampleInstructID,
				Content:  data,
				Priority: model.InstructPriorityDefault,
				State:    model.InstructCreate,
				ExpireAt: time.Now().Add(time.Hour),
			})
		},
	},
	{
		Define: "instruct list",
		Help:   `all bot instruct`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			models := make(map[string]interface{})
			for name, bot := range bots.List() {
				ruleset, _ := bot.Instruct()
				for _, rule := range ruleset {
					models[fmt.Sprintf("(%s) %s", name, rule.Id)] = fmt.Sprintf("[%s]", strings.Join(rule.Args, ","))
				}
			}
			return types.InfoMsg{
				Title: "Instruct",
				Model: models,
			}
		},
	},
}
