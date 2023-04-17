package dev

import (
	"github.com/tinode/chat/server/extra/pkg/template"
	"github.com/tinode/chat/server/extra/ruleset/event"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/logs"
)

var eventRules = []event.Rule{
	{
		Event: types.GroupEventJoin,
		Handler: func(ctx types.Context, head map[string]interface{}, content interface{}) types.MsgPayload {
			txt, err := template.Parse(ctx, "Welcome $username")
			if err != nil {
				logs.Err.Println(err)
				return types.TextMsg{Text: "error user"}
			}
			return types.TextMsg{Text: txt}
		},
	},
	{
		Event: types.GroupEventExit,
		Handler: func(ctx types.Context, head map[string]interface{}, content interface{}) types.MsgPayload {
			txt, err := template.Parse(ctx, "Bye $username")
			if err != nil {
				logs.Err.Println(err)
				return types.TextMsg{Text: "error user"}
			}
			return types.TextMsg{Text: txt}
		},
	},
	{
		Event: types.GroupEventReceive,
		Handler: func(ctx types.Context, head map[string]interface{}, content interface{}) types.MsgPayload {
			return types.TextMsg{Text: "receive something"}
		},
	},
}
