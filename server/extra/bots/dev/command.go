package dev

import (
	"bytes"
	"crypto/rand"
	_ "embed"
	"fmt"
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/pkg/parser"
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
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
			return nil
		},
	},
	{
		Define: "rand [number] [number]",
		Help:   `Generate random numbers`,
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
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
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
			return types.TextMsg{Text: types.Id().String()}
		},
	},
	{
		Define: "uid [string]",
		Help:   `Decode UID string`,
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
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
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
			num, _ := tokens[1].Value.Int64()
			t := time.Unix(num, 0)
			return types.TextMsg{Text: t.Format(time.RFC3339)}
		},
	},
	{
		Define: "messages",
		Help:   `[example] messages`,
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
			return types.TextMsg{Text: "msg1"}
		},
	},
	{
		Define: `qr [string]`,
		Help:   `Generate QR code`,
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
			return types.TextMsg{Text: "msg1"}
		},
	},
	{
		Define: `pinyin [string]`,
		Help:   "chinese pinyin conversion",
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
			return types.TextMsg{Text: "msg1"}
		},
	},
	{
		Define: "form",
		Help:   `[example] form`,
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
			return bots.FormMsg(ctx, devFormID)
		},
	},
	{
		Define: "action",
		Help:   "[example] action",
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
			return bots.ActionMsg(ctx, devActionID)
		},
	},
	{
		Define: "guess",
		Help:   "Guess number game",
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
			// rand number
			big, _ := rand.Int(rand.Reader, big.NewInt(1000))

			var initValue model.JSON
			initValue = map[string]interface{}{"number": big.Int64()}
			return bots.SessionMsg(ctx, guessSessionID, initValue)
		},
	},
	{
		Define: "echo [any]",
		Help:   "print",
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
			val := tokens[1].Value.Source
			return types.TextMsg{Text: fmt.Sprintf("%v", val)}
		},
	},
	{
		Define: "agent",
		Help:   `agent url`,
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
			return bots.AgentURI(ctx)
		},
	},
	{
		Define: "plot",
		Help:   `[example] plot graph`,
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
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
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
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
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
			data := model.JSON{}
			data["txt"] = "example"
			return bots.InstructMsg(ctx, ExampleInstructID, data)
		},
	},
	{
		Define: "instruct list",
		Help:   `all bot instruct`,
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
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
