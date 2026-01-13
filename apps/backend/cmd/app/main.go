package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/dimasbaguspm/fluxis/internal"
	"github.com/dimasbaguspm/fluxis/internal/configs"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

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
	internal.RegisterPrivateRoutes(ctx, humaApi, pool)

	slog.Info("All is ready! starting HTTP server", "port", env.AppPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", env.AppPort),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server error", "err", err)
		}
	}()

	// wait for shutdown signal
	<-ctx.Done()
	slog.Info("Shutdown signal received, shutting down HTTP server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Graceful shutdown failed, forcing exit", "err", err)
	} else {
		slog.Info("Server stopped")
	}
}
