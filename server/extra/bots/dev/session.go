package dev

import (
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/ruleset/session"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/store/model"
	"github.com/tinode/chat/server/extra/types"
	"strconv"
)

const (
	guessSessionID = "guess_session"
)

var sessionRules = []session.Rule{
	{
		Id:    guessSessionID,
		Title: "Input a number?",
		Handler: func(ctx types.Context, content interface{}) types.MsgPayload {
			number := int64(0)
			if v, ok := content.(string); ok {
				number, _ = strconv.ParseInt(v, 10, 64)
			} else {
				return types.TextMsg{Text: "input error"}
			}
			if number <= 0 {
				return types.TextMsg{Text: "input > 0 number"}
			}

			v, ok := ctx.SessionInitValues.Float64("number")
			if !ok {
				return types.TextMsg{Text: "init number error"}
			}
			initNumber := int64(v)

			// store current values
			_ = store.Chatbot.SessionSet(ctx.AsUser, ctx.Original, model.Session{Values: map[string]interface{}{"number": number}})

			if number == initNumber {
				bots.SessionDone(ctx)
				return types.TextMsg{Text: "Bingo"}
			} else if number > initNumber {
				return types.TextMsg{Text: "higher"}
			} else {
				return types.TextMsg{Text: "lower"}
			}
		},
	},
}
