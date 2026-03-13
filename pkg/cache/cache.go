package cache

import (
	"context"
	"time"
)

type Config struct {
	DefaultTTL time.Duration // default 15 minutes
	HMACKey    string        // required for signing cache keys
}

type Cache interface {
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
}
