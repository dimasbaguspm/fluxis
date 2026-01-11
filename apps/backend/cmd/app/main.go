package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/dimasbaguspm/fluxis/internal"
	"github.com/dimasbaguspm/fluxis/internal/configs"
)

func main() {
	ctx := context.Background()
	r := http.NewServeMux()
	env := configs.NewEnvironment()
	db := configs.NewDatabase(env)

	pool, err := db.Connect(ctx)
	if err != nil {
		slog.Error("Database is unreachable", "err", err)
		panic(err)
	}

	migration := configs.Migration(env)

	slog.Info("Performing migration")
	if err := migration.Up(); err != nil {
		slog.Error("Failed to run migration", "error", err.Error())
		panic(err)
	}
	slog.Info("DB migration completed")

	humaApi := humago.New(r, configs.GetOpenapiConfig(env))

	internal.RegisterPublicRoutes(humaApi, pool)
	internal.RegisterPrivateRoutes(r, humaApi, pool)

	slog.Info("All is ready! serving to port ", "Info", env.AppPort)

	http.ListenAndServe(fmt.Sprintf(":%s", env.AppPort), r)
}
