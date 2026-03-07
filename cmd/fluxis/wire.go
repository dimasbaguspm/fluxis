package main

import (
	"github.com/dimasbaguspm/fluxis/internal/auth"
	authhandler "github.com/dimasbaguspm/fluxis/internal/auth/handler"
	authservice "github.com/dimasbaguspm/fluxis/internal/auth/service"

	userrepo "github.com/dimasbaguspm/fluxis/internal/user/repository"
	userservice "github.com/dimasbaguspm/fluxis/internal/user/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Auth *auth.Module
}

type Deps struct {
	DB     *pgxpool.Pool
	Config *Config
}

func Wire(d Deps) *App {
	userRepo := userrepo.New(d.DB)

	userSvc := userservice.New(userservice.Deps{
		Repo: userRepo,
	})
	authSvc := authservice.New(authservice.Deps{
		Users:  userSvc,
		Config: &authservice.Config{},
	})

	authH := authhandler.New(authSvc)

	return &App{
		Auth: auth.NewModule(authSvc, authH),
	}

}
