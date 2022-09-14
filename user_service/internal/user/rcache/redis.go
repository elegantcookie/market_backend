package rcache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
	"user_service/internal/user"
	"user_service/pkg/logging"
)

type cache struct {
	client *redis.Client
	logger *logging.Logger
}

// TODO handle errors

func (c cache) Set(ctx context.Context, key, val string, expiration time.Duration) error {
	err := c.client.Set(ctx, key, val, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set data: %v", err)
	}
	return nil
}
func (c cache) Get(ctx context.Context, key string) (string, error) {
	get := c.client.Get(ctx, key)
	err := get.Err()
	if err != nil {
		return "", fmt.Errorf("failed to get value: %v", err)
	}
	val := get.Val()
	return val, nil
}
func (c cache) Del(ctx context.Context, key string) error {
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete value: %v", err)
	}
	return nil
}

func NewCache(client *redis.Client, logger *logging.Logger) user.Cache {
	return &cache{
		client: client,
		logger: logger,
	}

}
