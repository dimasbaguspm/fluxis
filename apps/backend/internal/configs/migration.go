package configs

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
)

type migration struct {
	env Environment
}

func Migration(env Environment) migration {
	return migration{env}
}

func (m migration) Up() error {
	return proceedMigration(m, "up")
}

func (m migration) Down() error {
	return proceedMigration(m, "down")
}

func proceedMigration(m migration, option string) error {
	migrator, err := migrate.New(
		fmt.Sprintf("file://%s", filepath.Clean("migrations")), m.env.Database.Url)

	if err != nil {
		return fmt.Errorf("Failed to setup migration: %w", err)
	}
	defer migrator.Close()

	var postErr error

	switch option {
	case "up":
		if err := migrator.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				return nil
			}
			postErr = fmt.Errorf("failed to run migrations: %w", err)
		}
	case "down":
		if err := migrator.Down(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				return nil
			}
			postErr = fmt.Errorf("failed to run migrations: %w", err)
		}
	}

	if postErr != nil {
		return postErr
	}

	return nil
}
