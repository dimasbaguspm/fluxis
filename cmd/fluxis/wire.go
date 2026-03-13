package main

import (
	"github.com/dimasbaguspm/fluxis/internal/auth"
	authhandler "github.com/dimasbaguspm/fluxis/internal/auth/handler"
	authservice "github.com/dimasbaguspm/fluxis/internal/auth/service"

	"github.com/dimasbaguspm/fluxis/internal/user"
	usercache "github.com/dimasbaguspm/fluxis/internal/user/cache"
	userhandler "github.com/dimasbaguspm/fluxis/internal/user/handler"
	userrepo "github.com/dimasbaguspm/fluxis/internal/user/repository"
	userservice "github.com/dimasbaguspm/fluxis/internal/user/service"

	"github.com/dimasbaguspm/fluxis/internal/org"
	orgcache "github.com/dimasbaguspm/fluxis/internal/org/cache"
	orghandler "github.com/dimasbaguspm/fluxis/internal/org/handler"
	orgrepo "github.com/dimasbaguspm/fluxis/internal/org/repository"
	orgservice "github.com/dimasbaguspm/fluxis/internal/org/service"

	"github.com/dimasbaguspm/fluxis/internal/project"
	projectcache "github.com/dimasbaguspm/fluxis/internal/project/cache"
	projecthandler "github.com/dimasbaguspm/fluxis/internal/project/handler"
	projectrepo "github.com/dimasbaguspm/fluxis/internal/project/repository"
	projectservice "github.com/dimasbaguspm/fluxis/internal/project/service"

	"github.com/dimasbaguspm/fluxis/internal/sprint"
	sprintcache "github.com/dimasbaguspm/fluxis/internal/sprint/cache"
	sprinthandler "github.com/dimasbaguspm/fluxis/internal/sprint/handler"
	sprintrepo "github.com/dimasbaguspm/fluxis/internal/sprint/repository"
	sprintservice "github.com/dimasbaguspm/fluxis/internal/sprint/service"

	"github.com/dimasbaguspm/fluxis/internal/board"
	boardcache "github.com/dimasbaguspm/fluxis/internal/board/cache"
	boardhandler "github.com/dimasbaguspm/fluxis/internal/board/handler"
	boardrepo "github.com/dimasbaguspm/fluxis/internal/board/repository"
	boardservice "github.com/dimasbaguspm/fluxis/internal/board/service"

	"github.com/dimasbaguspm/fluxis/internal/ticket"
	ticketcache "github.com/dimasbaguspm/fluxis/internal/ticket/cache"
	tickethandler "github.com/dimasbaguspm/fluxis/internal/ticket/handler"
	ticketrepo "github.com/dimasbaguspm/fluxis/internal/ticket/repository"
	ticketservice "github.com/dimasbaguspm/fluxis/internal/ticket/service"

	"github.com/dimasbaguspm/fluxis/pkg/cache"
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
	DB        *pgxpool.Pool
	Config    *Config
	Bus       pubsub.Bus
	DataCache cache.Cache
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

	userC := usercache.New(d.DataCache)
	orgC := orgcache.New(d.DataCache)
	projectC := projectcache.New(d.DataCache)
	sprintC := sprintcache.New(d.DataCache)
	boardC := boardcache.New(d.DataCache)
	ticketC := ticketcache.New(d.DataCache)

	authH := authhandler.New(authSvc)
	userH := userhandler.New(userhandler.Deps{
		Svc:       userSvc,
		UserCache: userC,
	})
	orgH := orghandler.New(orghandler.Deps{
		Svc:      orgSvc,
		OrgCache: orgC,
	})
	projectH := projecthandler.New(projecthandler.Deps{
		Svc:          projectSvc,
		ProjectCache: projectC,
	})
	sprintH := sprinthandler.New(sprinthandler.Deps{
		Svc:         sprintSvc,
		SprintCache: sprintC,
	})
	boardH := boardhandler.New(boardhandler.Deps{
		Svc:        boardSvc,
		BoardCache: boardC,
	})
	ticketH := tickethandler.New(tickethandler.Deps{
		Svc:         ticketSvc,
		TicketCache: ticketC,
	})

	return &App{
		Auth:    auth.NewModule(authSvc, authH, d.Bus),
		User:    user.NewModule(userH, userC, d.Bus),
		Org:     org.NewModule(orgH, orgC, d.Bus),
		Project: project.NewModule(projectH, projectC, d.Bus),
		Sprint:  sprint.NewModule(sprintH, sprintC, d.Bus),
		Board:   board.NewModule(boardH, boardC, d.Bus),
		Ticket:  ticket.NewModule(ticketH, ticketC, d.Bus),
	}

}
