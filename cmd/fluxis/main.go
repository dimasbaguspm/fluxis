package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dimasbaguspm/fluxis/pkg/cache"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
	"github.com/dimasbaguspm/fluxis/pkg/postgres"
	"github.com/dimasbaguspm/fluxis/pkg/pubsub"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title					Fluxis API
// @version					1.0
// @description				Personal finance management API
//
// @contact.name			Fluxis Support
// @contact.url				https://github.com/dimasbaguspm/fluxis
//
// @license.name			MIT
//
// @host					localhost:8080
// @BasePath				/
// @schemes					http https
//
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description					Bearer token obtained from /auth/login or /auth/refresh
func main() {
	cfg := LoadEnv()

	ctx, close := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT)
	defer close()

	db := postgres.MustConnect(ctx, cfg.DB)
	postgres.RunMigration(cfg.DB)

	bus := pubsub.New()
	cache := cache.New(cfg.Cache)

	defer db.Close()
	defer bus.Close()

	app := Wire(Deps{
		DB:     db,
		Config: cfg,
		Bus:    bus,
		Cache:  cache,
	})

	httpx.InitAuth(app.Auth.Service())

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("GET /swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./api/swagger.json")
	})
	mux.Handle("GET /swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// mount domain routes onto the mux
	// each domain registers its own paths
	app.Auth.Routes(mux)
	app.User.Routes(mux)
	app.Org.Routes(mux)
	app.Project.Routes(mux)
	app.Sprint.Routes(mux)
	app.Board.Routes(mux)
	app.Ticket.Routes(mux)

	// start event subscribers
	go app.Auth.StartSubscriber(ctx)
	go app.User.StartSubscriber(ctx)
	go app.Org.StartSubscriber(ctx)
	go app.Project.StartSubscriber(ctx)
	go app.Sprint.StartSubscriber(ctx)
	go app.Board.StartSubscriber(ctx)
	go app.Ticket.StartSubscriber(ctx)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		httpx.Handle(w, httpx.NotImplemented("endpoint is not implemented"))
	})

	svr := http.Server{
		Addr:         cfg.Server.addr(),
		Handler:      mux,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		slog.Info(fmt.Sprintf("[Core]: Server started in port %s", cfg.Server.Port))
		if err := svr.ListenAndServe(); err != http.ErrServerClosed {
			slog.Error("[Core]: Failed to start the serve", "error", err)
		}
	}()

	<-ctx.Done()

}
