package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Primary  string
	MaxConns int
	MinConns int
}

func MustConnect(ctx context.Context, cfg Config) *pgxpool.Pool {
	slog.Info("[Database]: Attempting to connect the database")

	config, err := pgxpool.ParseConfig(cfg.Primary)
	config.MinConns = int32(cfg.MinConns)
	config.MaxConns = int32(cfg.MaxConns)

	conn, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		slog.Error(fmt.Sprintf("[Database]: Unable to connect with db, %v", err))
		os.Exit(1)
		return nil
	}

	slog.Info("[Database]: Connection established")
	return conn
}
