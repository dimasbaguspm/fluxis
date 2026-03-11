package apitest_test

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type PostgresContainer struct {
	container *postgres.PostgresContainer
	DSN       string
}

func NewPostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	container, err := postgres.RunContainer(ctx,
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, fmt.Errorf("failed to get connection string: %w", err)
	}

	return &PostgresContainer{
		container: container,
		DSN:       dsn,
	}, nil
}

// Terminate stops the postgres container.
func (pc *PostgresContainer) Terminate(ctx context.Context) error {
	return pc.container.Terminate(ctx)
}

// MustPool creates a pgxpool.Pool from a DSN, panicking on error.
// It retries for up to 30 seconds to allow the database to be ready.
func MustPool(ctx context.Context, dsn string) *pgxpool.Pool {
	deadline := time.Now().Add(30 * time.Second)
	var lastErr error
	for time.Now().Before(deadline) {
		pool, err := pgxpool.New(ctx, dsn)
		if err == nil {
			// Test the connection
			pingErr := pool.Ping(ctx)
			if pingErr == nil {
				return pool
			}
			pool.Close()
			lastErr = pingErr
		} else {
			lastErr = err
		}
		time.Sleep(500 * time.Millisecond)
	}
	panic(fmt.Sprintf("failed to create pool after 30 seconds: %v", lastErr))
}

// RunMigrations runs golang-migrate migrations from the given path.
// Tolerates ErrNoChange (no migrations to run).
func RunMigrations(dsn, migrationsPath string) error {
	m, err := migrate.New(fmt.Sprintf("file://%s", migrationsPath), dsn)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
