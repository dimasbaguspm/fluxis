package main

import (
	"github.com/dimasbaguspm/fluxis/internal/auth"
	authhandler "github.com/dimasbaguspm/fluxis/internal/auth/handler"
	authservice "github.com/dimasbaguspm/fluxis/internal/auth/service"
	"github.com/redis/go-redis/v9"

	"github.com/dimasbaguspm/fluxis/internal/user"
	userhandler "github.com/dimasbaguspm/fluxis/internal/user/handler"
	userrepo "github.com/dimasbaguspm/fluxis/internal/user/repository"
	userservice "github.com/dimasbaguspm/fluxis/internal/user/service"

	"github.com/dimasbaguspm/fluxis/internal/org"
	orghandler "github.com/dimasbaguspm/fluxis/internal/org/handler"
	orgrepo "github.com/dimasbaguspm/fluxis/internal/org/repository"
	orgservice "github.com/dimasbaguspm/fluxis/internal/org/service"

	"github.com/dimasbaguspm/fluxis/internal/project"
	projecthandler "github.com/dimasbaguspm/fluxis/internal/project/handler"
	projectrepo "github.com/dimasbaguspm/fluxis/internal/project/repository"
	projectservice "github.com/dimasbaguspm/fluxis/internal/project/service"

	"github.com/dimasbaguspm/fluxis/internal/sprint"
	sprinthandler "github.com/dimasbaguspm/fluxis/internal/sprint/handler"
	sprintrepo "github.com/dimasbaguspm/fluxis/internal/sprint/repository"
	sprintservice "github.com/dimasbaguspm/fluxis/internal/sprint/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Auth    *auth.Module
	User    *user.Module
	Org     *org.Module
	Project *project.Module
	Sprint  *sprint.Module
}

type Deps struct {
	DB     *pgxpool.Pool
	RDB    *redis.Client
	Config *Config
}

func Wire(d Deps) *App {
	userRepo := userrepo.New(d.DB)
	orgRepo := orgrepo.New(d.DB)
	projectRepo := projectrepo.New(d.DB)
	sprintRepo := sprintrepo.New(d.DB)

	userSvc := userservice.New(userservice.Deps{
		Repo: userRepo,
	})
	orgSvc := orgservice.New(orgservice.Deps{
		Repo: orgRepo,
	})
	projectSvc := projectservice.New(projectservice.Deps{
		Repo: projectRepo,
	})
	sprintSvc := sprintservice.New(sprintservice.Deps{
		Repo: sprintRepo,
	})
	authSvc := authservice.New(authservice.Deps{
		Users:  userSvc,
		Config: &d.Config.Auth,
	})

	authH := authhandler.New(authSvc)
	userH := userhandler.New(userSvc)
	orgH := orghandler.New(orgSvc)
	projectH := projecthandler.New(projectSvc)
	sprintH := sprinthandler.New(sprintSvc)

	return &App{
		Auth:    auth.NewModule(authSvc, authH),
		User:    user.NewModule(userH),
		Org:     org.NewModule(orgH),
		Project: project.NewModule(projectH),
		Sprint:  sprint.NewModule(sprintH),
	}

}
