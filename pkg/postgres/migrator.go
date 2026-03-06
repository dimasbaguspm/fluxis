package postgres

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrator() error {
	slog.Info("[Migrator]: trying to migrate tables into DB")

	exe, err := os.Executable()
	if err != nil {
		slog.Error("[Migrator]: failed to resolve executable path", "error", err)
		return err
	}
	migrationsPath := filepath.Clean(filepath.Join(filepath.Dir(exe), "..", "migrations"))

	m, err := migrate.New("file://"+migrationsPath, dbUrl)

	if err != nil {
		slog.Error("[Migrator]: migration failed something odd while lookup the migrations file", "error", err)
		return err
	}

	err = m.Up()

	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			slog.Info("[Migrator]: migration success without no change!")
			return nil
		}

		slog.Error("[Migrator]: unable to migrate the db", "error", err)
		return err
	}

	slog.Info("[Migrator]: success to migrate the latest version!")
	return nil
}
