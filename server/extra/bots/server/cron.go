package server

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/tinode/chat/server/extra/cache"
	"github.com/tinode/chat/server/extra/ruleset/cron"
	"github.com/tinode/chat/server/extra/types"
	"time"
)

var cronRules = []cron.Rule{
	{
		Name: "server_user_online_change",
		When: "* * * * *",
		Action: func(ctx types.Context) []types.MsgPayload {
			ctx_ := context.Background()
			keys, _ := cache.DB.Keys(ctx_, "online:*").Result()

			currentCount := int64(len(keys))
			lastKey := fmt.Sprintf("server:cron:online_count_last:%s", ctx.AsUser.UserId())

			lastCount, _ := cache.DB.Get(ctx_, lastKey).Int64()
			cache.DB.Set(ctx_, lastKey, currentCount, redis.KeepTTL)

			if lastCount != currentCount {
				return []types.MsgPayload{
					types.TextMsg{Text: fmt.Sprintf("online change %d (%d)", currentCount-lastCount, time.Now().Unix())},
				}
			}
			return nil
		},
	},
}
