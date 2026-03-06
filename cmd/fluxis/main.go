package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dimasbaguspm/fluxis/internal/auth"
	"github.com/dimasbaguspm/fluxis/internal/user"
	"github.com/dimasbaguspm/fluxis/pkg/postgres"
)

func main() {
	ctx, close := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT)
	defer close()

	_, err := postgres.Pool(ctx)
	if err != nil {
		os.Exit(1)
	}

	err = postgres.Migrator()
	if err != nil {
		os.Exit(1)
	}

	mux := http.NewServeMux()
	auth.Routes(mux)
	user.Routes(mux)

	svr := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	slog.Info("Start to serve app into port :8080")

	go func() {
		if err := svr.ListenAndServe(); err != http.ErrServerClosed {
			slog.Error("Failed to start the serve", "error", err)
		}
	}()

	fmt.Println("I am here")
	<-ctx.Done()

}
