package postgres

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func MustConnect(ctx context.Context, dbUrl string) *pgxpool.Pool {
	slog.Info("[Database]: Attempting to connect the database")
	conn, err := pgxpool.New(ctx, dbUrl)

	if err != nil {
		slog.Error("[Database]: Unable to connect with db", "error", err)
		os.Exit(1)
		return nil
	}

	slog.Info("[Database]: Connection established")
	return conn
}
