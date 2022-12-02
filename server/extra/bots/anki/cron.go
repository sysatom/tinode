package anki

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/tinode/chat/server/extra/cache"
	"github.com/tinode/chat/server/extra/ruleset/cron"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/types"
	"strconv"
	"time"
)

var cronRules = []cron.Rule{
	{
		Name: "anki_review_remind",
		When: "* * * * *",
		Action: func(ctx types.Context) []types.MsgPayload {
			j, err := store.Chatbot.DataGet(ctx.AsUser, ctx.Original, "getNumCardsReviewedToday")
			if err != nil {
				return nil
			}
			v, ok := j.Float64("value")
			if !ok {
				return nil
			}
			num := int64(v)
			if num == 0 {
				return nil
			}
			key := fmt.Sprintf("anki:review_remind:%d", ctx.AsUser)

			sendString, err := cache.DB.Get(context.Background(), key).Result()
			if err != nil && !errors.Is(err, redis.Nil) {
				return nil
			}
			oldSend := int64(0)
			if len(sendString) != 0 {
				oldSend, _ = strconv.ParseInt(sendString, 10, 64)
			}

			if time.Now().Unix()-oldSend > 24*60*60 {
				_ = cache.DB.Set(context.Background(), key, strconv.FormatInt(time.Now().Unix(), 10), redis.KeepTTL)

				return []types.MsgPayload{
					types.TextMsg{Text: fmt.Sprintf("Anki review %d", num)},
				}
			}

			return nil
		},
	},
}
