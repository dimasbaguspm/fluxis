package cache

import (
	"log/slog"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	rdb *redis.Client
	cfg Config
}

func New(cfg Config, rdb *redis.Client) RedisClient {
	slog.Info("[Cache]: initializing redis as caching client")
	return RedisClient{
		rdb: rdb,
		cfg: cfg,
	}
}
