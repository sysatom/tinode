package clipboard

import (
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/ruleset/command"
	"github.com/tinode/chat/server/extra/store/model"
	"github.com/tinode/chat/server/extra/types"
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
		Define: "share [string]",
		Help:   `share clipboard to helper`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			txt, _ := tokens[1].Value.String()
			data := model.JSON{}
			data["txt"] = txt
			return bots.StoreInstruct(ctx, types.InstructMsg{
				No:       types.Id().String(),
				Object:   model.InstructObjectHelper,
				Bot:      Name,
				Flag:     "clipboard_share",
				Content:  data,
				Priority: model.InstructPriorityDefault,
				State:    model.InstructCreate,
				ExpireAt: time.Now().Add(time.Hour),
			})
		},
	},
}
