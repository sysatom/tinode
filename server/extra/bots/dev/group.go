package dev

import (
	"github.com/tinode/chat/server/extra/ruleset/event"
	"github.com/tinode/chat/server/extra/types"
)

var eventRules = []event.Rule{
	{
		Event: types.GroupEventJoin,
		Handler: func(ctx types.Context, head map[string]interface{}, content interface{}) types.MsgPayload {
			return types.TextMsg{Text: "Welcome"}
		},
	},
	{
		Event: types.GroupEventReceive,
		Handler: func(ctx types.Context, head map[string]interface{}, content interface{}) types.MsgPayload {
			return types.TextMsg{Text: "receive something"}
		},
	},
}
