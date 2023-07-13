package markdown

import (
	"fmt"
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/pkg/parser"
	"github.com/tinode/chat/server/extra/pkg/route"
	"github.com/tinode/chat/server/extra/ruleset/command"
	"github.com/tinode/chat/server/extra/store/model"
	"github.com/tinode/chat/server/extra/types"
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
		Define: "editor",
		Help:   `Bot info`,
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
			p := model.JSON{}
			p["uid"] = ctx.AsUser.UserId()
			flag, err := bots.StoreParameter(p, time.Now().Add(time.Hour))
			if err != nil {
				return types.TextMsg{Text: "error parameter"}
			}
			return types.LinkMsg{
				Title: "Markdown Editor",
				Url:   route.URL(Name, serviceVersion, fmt.Sprintf("editor/%s", flag)),
			}
		},
	},
}
