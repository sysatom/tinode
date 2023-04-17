package dev

import (
	"fmt"
	"github.com/tinode/chat/server/extra/ruleset/event"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/store"
)

var eventRules = []event.Rule{
	{
		Event: types.GroupEventJoin,
		Handler: func(ctx types.Context, head map[string]interface{}, content interface{}) types.MsgPayload {
			user, err := store.Users.Get(ctx.AsUser)
			if err != nil {
				return types.TextMsg{Text: "error user"}
			}
			return types.TextMsg{Text: fmt.Sprintf("Welcome %s", user.Public)}
		},
	},
	{
		Event: types.GroupEventExit,
		Handler: func(ctx types.Context, head map[string]interface{}, content interface{}) types.MsgPayload {
			user, err := store.Users.Get(ctx.AsUser)
			if err != nil {
				return types.TextMsg{Text: "error user"}
			}
			return types.TextMsg{Text: fmt.Sprintf("Byebye %s", user.Public)}
		},
	},
	{
		Event: types.GroupEventReceive,
		Handler: func(ctx types.Context, head map[string]interface{}, content interface{}) types.MsgPayload {
			return types.TextMsg{Text: "receive something"}
		},
	},
}
