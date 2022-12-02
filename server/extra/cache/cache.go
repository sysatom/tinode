package cache

import (
	"context"
	"github.com/go-redis/redis/v9"
	"os"
)

var DB *redis.Client

func InitCache() {
	addr := os.Getenv("REDIS_ADDR")
	password := os.Getenv("REDIS_PASSWORD")
	if addr == "" || password == "" {
		panic("redis config error")
	}
	DB = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})
	s := DB.Ping(context.Background())
	_, err := s.Result()
	if err != nil {
		panic("redis server error")
	}
}
