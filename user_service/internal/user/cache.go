package user

import (
	"context"
	"time"
)

type Cache interface {
	Set(ctx context.Context, key, val string, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
}
