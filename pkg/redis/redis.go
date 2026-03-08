package redis

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Addr     string
	Password string
}

func MustConnect(ctx context.Context, cfg Config) *redis.Client {
	slog.Info("[Redis]: Connecting to redis database")
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       0,
	})
	slog.Info("[Redis]: Connection established")

	return rdb

}
