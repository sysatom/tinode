package dev

import (
	"github.com/tinode/chat/server/extra/ruleset/cron"
	"github.com/tinode/chat/server/extra/types"
)

var cronRules = []cron.Rule{
	{
		Name: "dev_demo",
		Help: "cron example",
		When: "0 */1 * * *",
		Action: func(types.Context) []types.MsgPayload {
			return nil
		},
	},
}
