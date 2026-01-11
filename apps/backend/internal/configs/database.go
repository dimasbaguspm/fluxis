package configs

import (
	"context"
	"errors"

	pgx "github.com/jackc/pgx/v5/pgxpool"
)

type db struct {
	env Environment
}

func (db db) Connect(ctx context.Context) (*pgx.Pool, error) {
	pool, err := pgx.New(ctx, db.env.Database.Url)
	if err != nil {
		return nil, errors.New("Unable to create connection pool")
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, errors.New("Database is unreachable")
	}

	return pool, nil
}

func NewDatabase(env Environment) db {
	return db{env}
}
