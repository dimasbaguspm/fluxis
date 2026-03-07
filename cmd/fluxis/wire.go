package main

import (
	"github.com/dimasbaguspm/fluxis/internal/auth"
	authhandler "github.com/dimasbaguspm/fluxis/internal/auth/handler"
	authservice "github.com/dimasbaguspm/fluxis/internal/auth/service"

	"github.com/dimasbaguspm/fluxis/internal/user"
	userhandler "github.com/dimasbaguspm/fluxis/internal/user/handler"
	userrepo "github.com/dimasbaguspm/fluxis/internal/user/repository"
	userservice "github.com/dimasbaguspm/fluxis/internal/user/service"

	"github.com/dimasbaguspm/fluxis/internal/org"
	orghandler "github.com/dimasbaguspm/fluxis/internal/org/handler"
	orgrepo "github.com/dimasbaguspm/fluxis/internal/org/repository"
	orgservice "github.com/dimasbaguspm/fluxis/internal/org/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Auth *auth.Module
	User *user.Module
	Org  *org.Module
}

type Deps struct {
	DB     *pgxpool.Pool
	Config *Config
}

func Wire(d Deps) *App {
	userRepo := userrepo.New(d.DB)
	orgRepo := orgrepo.New(d.DB)

	userSvc := userservice.New(userservice.Deps{
		Repo: userRepo,
	})
	orgSvc := orgservice.New(orgservice.Deps{
		Repo: orgRepo,
	})
	authSvc := authservice.New(authservice.Deps{
		Users:  userSvc,
		Config: &d.Config.Auth,
	})

	authH := authhandler.New(authSvc)
	userH := userhandler.New(userSvc)
	orgH := orghandler.New(orgSvc)

	return &App{
		Auth: auth.NewModule(authSvc, authH),
		User: user.NewModule(userH),
		Org:  org.NewModule(orgH),
	}

}
