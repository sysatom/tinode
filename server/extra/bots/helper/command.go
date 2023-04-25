package helper

import (
	"errors"
	"fmt"
	"github.com/tinode/chat/server/extra/pkg/parser"
	"github.com/tinode/chat/server/extra/ruleset/command"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/types"
	"gorm.io/gorm"
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
		Define: "access url",
		Help:   `get access url`,
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
			// get id
			value, err := store.Chatbot.ConfigGet(ctx.AsUser, "", fmt.Sprintf("helper:%d", uint64(ctx.AsUser)))
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			idValue, ok := value.String("value")
			if !ok || idValue == "" {
				idValue = types.Id().String()
				// set id
				err = store.Chatbot.ConfigSet(ctx.AsUser, "",
					fmt.Sprintf("helper:%d", uint64(ctx.AsUser)), map[string]interface{}{
						"value": idValue,
					})
				if err != nil {
					return nil
				}
			}

			return types.TextMsg{Text: fmt.Sprintf("%s/extra/helper/%d/%s", types.AppUrl(),
				uint64(ctx.AsUser), idValue)}
		},
	},
}
