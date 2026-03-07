package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dimasbaguspm/fluxis/pkg/httpx"
	"github.com/dimasbaguspm/fluxis/pkg/postgres"
)

// @title           Fluxis API
// @version         1.0
func main() {
	cfg := LoadEnv()

	ctx, close := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT)
	defer close()

	db := postgres.MustConnect(ctx, cfg.DB.Primary)
	postgres.RunMigration(cfg.DB.Primary)

	defer db.Close()

	app := Wire(Deps{
		DB:     db,
		Config: cfg,
	})

	httpx.InitAuth(app.Auth.Service())

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// mount domain routes onto the mux
	// each domain registers its own paths
	app.Auth.Routes(mux)

	svr := http.Server{
		Addr:         cfg.Server.addr(),
		Handler:      mux,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		slog.Info(fmt.Sprintf("Server started in port %s", cfg.Server.Port))
		if err := svr.ListenAndServe(); err != http.ErrServerClosed {
			slog.Error("Failed to start the serve", "error", err)
		}
	}()

	<-ctx.Done()

}
