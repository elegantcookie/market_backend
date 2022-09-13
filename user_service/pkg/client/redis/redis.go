package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

func NewClient() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := rdb.Set(ctx, "key", "val", 0).Err()
	if err != nil {
		return nil, err
	}
	err = rdb.Del(ctx, "key").Err()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}
