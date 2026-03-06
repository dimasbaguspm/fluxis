package postgres

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var dbUrl = os.Getenv("DATABASE_URL")

func Pool(ctx context.Context) (*pgxpool.Pool, error) {
	slog.Info("[Database]: Attempting to connect the database")
	conn, err := pgxpool.New(ctx, dbUrl)

	if err != nil {
		slog.Error("[Database]: Unable to connect with db", "error", err)
		return nil, err
	}

	slog.Info("[Database]: Connection established")
	return conn, nil
}
