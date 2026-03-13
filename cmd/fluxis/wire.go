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

	"github.com/dimasbaguspm/fluxis/internal/board"
	boardhandler "github.com/dimasbaguspm/fluxis/internal/board/handler"
	boardrepo "github.com/dimasbaguspm/fluxis/internal/board/repository"
	boardservice "github.com/dimasbaguspm/fluxis/internal/board/service"

	"github.com/dimasbaguspm/fluxis/internal/ticket"
	tickethandler "github.com/dimasbaguspm/fluxis/internal/ticket/handler"
	ticketrepo "github.com/dimasbaguspm/fluxis/internal/ticket/repository"
	ticketservice "github.com/dimasbaguspm/fluxis/internal/ticket/service"

	"github.com/dimasbaguspm/fluxis/pkg/pubsub"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Auth    *auth.Module
	User    *user.Module
	Org     *org.Module
	Project *project.Module
	Sprint  *sprint.Module
	Board   *board.Module
	Ticket  *ticket.Module
}

type Deps struct {
	DB     *pgxpool.Pool
	RDB    *redis.Client
	Config *Config
	Bus    pubsub.Bus
}

func Wire(d Deps) *App {
	userRepo := userrepo.New(d.DB)
	orgRepo := orgrepo.New(d.DB)
	projectRepo := projectrepo.New(d.DB)
	sprintRepo := sprintrepo.New(d.DB)
	boardRepo := boardrepo.New(d.DB)
	ticketRepo := ticketrepo.New(d.DB)

	userSvc := userservice.New(userservice.Deps{
		Repo: userRepo,
	})
	authSvc := authservice.New(authservice.Deps{
		Users:  userSvc,
		Config: &d.Config.Auth,
	})
	orgSvc := orgservice.New(orgservice.Deps{
		Repo: orgRepo,
		User: userSvc,
		Bus:  d.Bus,
	})
	projectSvc := projectservice.New(projectservice.Deps{
		Repo: projectRepo,
		Org:  orgSvc,
		Bus:  d.Bus,
	})
	sprintSvc := sprintservice.New(sprintservice.Deps{
		Repo:    sprintRepo,
		Project: projectSvc,
		Bus:     d.Bus,
	})
	boardSvc := boardservice.New(boardservice.Deps{
		Repo:   boardRepo,
		Sprint: sprintSvc,
		Bus:    d.Bus,
	})
	ticketSvc := ticketservice.New(ticketservice.Deps{
		Repo:    ticketRepo,
		Project: projectSvc,
		Board:   boardSvc,
		Sprint:  sprintSvc,
		Bus:     d.Bus,
	})

	authH := authhandler.New(authSvc)
	userH := userhandler.New(userSvc)
	orgH := orghandler.New(orgSvc)
	projectH := projecthandler.New(projectSvc)
	sprintH := sprinthandler.New(sprintSvc)
	boardH := boardhandler.New(boardSvc)
	ticketH := tickethandler.New(ticketSvc)

	return &App{
		Auth:    auth.NewModule(authSvc, authH),
		User:    user.NewModule(userH),
		Org:     org.NewModule(orgH),
		Project: project.NewModule(projectH),
		Sprint:  sprint.NewModule(sprintH),
		Board:   board.NewModule(boardH),
		Ticket:  ticket.NewModule(ticketH),
	}

}
