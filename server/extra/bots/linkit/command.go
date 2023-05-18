package linkit

import (
	"errors"
	"fmt"
	"github.com/tinode/chat/server/extra/pkg/parser"
	"github.com/tinode/chat/server/extra/ruleset/command"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/store/model"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/extra/utils"
	"github.com/tinode/chat/server/logs"
	"gorm.io/gorm"
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
		Define: "token",
		Help:   `get access token`,
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
			// get token
			value, err := store.Chatbot.ConfigGet(ctx.AsUser, "", fmt.Sprintf("linkit:%d:token", uint64(ctx.AsUser)))
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			idValue, ok := value.String("value")
			if !ok || idValue == "" {
				idValue, err = utils.GenerateRandomString(25)
				if err != nil {
					return types.TextMsg{Text: "generate error"}
				}
				idValue = strings.ToLower(idValue)
				// set token
				err = store.Chatbot.ConfigSet(ctx.AsUser, "",
					fmt.Sprintf("linkit:%d:token", uint64(ctx.AsUser)), map[string]interface{}{
						"value": idValue,
					})
				if err != nil {
					logs.Err.Println(err)
					return types.TextMsg{Text: "set token error"}
				}
				data := model.JSON{}
				data["uid"] = ctx.AsUser.UserId()
				err = store.Chatbot.ParameterSet(idValue, data, time.Now().AddDate(1, 0, 0))
				if err != nil {
					logs.Err.Println(err)
					return types.TextMsg{Text: "set token error"}
				}
			}

			return types.TextMsg{Text: fmt.Sprintf("[One-year validity token] %s", idValue)}
		},
	},
	{
		Define: "reset token",
		Help:   `reset access token`,
		Handler: func(ctx types.Context, tokens []*parser.Token) types.MsgPayload {
			// get old token
			value, err := store.Chatbot.ConfigGet(ctx.AsUser, "", fmt.Sprintf("linkit:%d:token", uint64(ctx.AsUser)))
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			idValue, ok := value.String("value")
			if !ok || idValue == "" {
				return nil
			}
			// expire old token
			err = store.Chatbot.ParameterSet(idValue, model.JSON{}, time.Now())
			if err != nil {
				return nil
			}

			// new token
			idValue, err = utils.GenerateRandomString(25)
			if err != nil {
				return types.TextMsg{Text: "generate error"}
			}
			idValue = strings.ToLower(idValue)
			// set token
			err = store.Chatbot.ConfigSet(ctx.AsUser, "",
				fmt.Sprintf("linkit:%d:token", uint64(ctx.AsUser)), map[string]interface{}{
					"value": idValue,
				})
			if err != nil {
				logs.Err.Println(err)
				return types.TextMsg{Text: "set token error"}
			}
			data := model.JSON{}
			data["uid"] = ctx.AsUser.UserId()
			err = store.Chatbot.ParameterSet(idValue, data, time.Now().AddDate(1, 0, 0))
			if err != nil {
				logs.Err.Println(err)
				return types.TextMsg{Text: "set token error"}
			}

			return types.TextMsg{Text: fmt.Sprintf("[One-year validity token] %s", idValue)}
		},
	},
}
