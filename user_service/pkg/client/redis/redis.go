package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

func NewClient() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	return rdb, rdb.Ping(ctx).Err()
}
