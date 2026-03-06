package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, close := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT)
	defer close()

	mut := http.NewServeMux()

	svr := http.Server{
		Addr:    ":8080",
		Handler: mut,
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
