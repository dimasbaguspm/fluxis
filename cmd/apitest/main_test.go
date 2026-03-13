package apitest_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dimasbaguspm/fluxis/internal/auth"
	authhandler "github.com/dimasbaguspm/fluxis/internal/auth/handler"
	authservice "github.com/dimasbaguspm/fluxis/internal/auth/service"

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

	"github.com/dimasbaguspm/fluxis/internal/user"
	usercache "github.com/dimasbaguspm/fluxis/internal/user/cache"
	userhandler "github.com/dimasbaguspm/fluxis/internal/user/handler"
	userrepo "github.com/dimasbaguspm/fluxis/internal/user/repository"
	userservice "github.com/dimasbaguspm/fluxis/internal/user/service"

	"github.com/dimasbaguspm/fluxis/pkg/cache"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
	"github.com/dimasbaguspm/fluxis/pkg/pubsub"
)

var testServer *httptest.Server

func TestMain(m *testing.M) {
	ctx := context.Background()

	pgContainer, err := NewPostgresContainer(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start postgres container: %v\n", err)
		os.Exit(1)
	}
	defer pgContainer.Terminate(ctx)

	pool := MustPool(ctx, pgContainer.DSN)
	defer pool.Close()

	var migrationsPath string
	possiblePaths := []string{
		filepath.Join(os.Getenv("PWD"), "..", "..", "migrations"),
		"../../../migrations",
		"migrations",
	}
	for _, p := range possiblePaths {
		absPath, _ := filepath.Abs(p)
		if info, err := os.Stat(absPath); err == nil && info.IsDir() {
			migrationsPath = absPath
			break
		}
	}
	if migrationsPath == "" {
		fmt.Fprintf(os.Stderr, "Failed to find migrations directory\n")
		os.Exit(1)
	}

	if err := RunMigrations(pgContainer.DSN, migrationsPath); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run migrations: %v\n", err)
		os.Exit(1)
	}

	authCfg := authservice.Config{
		AccessTokenSecret:  "test-access-secret-32-chars-long-xxx",
		RefreshTokenSecret: "test-refresh-secret-32-chars-long-xx",
		AccessTokenExpiry:  1 * time.Minute,
		RefreshTokenExpiry: 5 * time.Minute,
		BcryptCost:         4,
	}

	userRepo := userrepo.New(pool)
	orgRepo := orgrepo.New(pool)
	projectRepo := projectrepo.New(pool)
	sprintRepo := sprintrepo.New(pool)
	boardRepo := boardrepo.New(pool)
	ticketRepo := ticketrepo.New(pool)

	bus := pubsub.New()
	defer bus.Close()

	cacheCfg := cache.Config{
		DefaultTTL: 15 * time.Minute,
		HMACKey:    "test-cache-hmac-key-32-chars-long",
	}
	memCache := cache.New(cacheCfg)

	userSvc := userservice.New(userservice.Deps{
		Repo: userRepo,
	})
	orgSvc := orgservice.New(orgservice.Deps{
		Repo: orgRepo,
		User: userSvc,
		Bus:  bus,
	})
	projectSvc := projectservice.New(projectservice.Deps{
		Repo: projectRepo,
		Org:  orgSvc,
		Bus:  bus,
	})
	sprintSvc := sprintservice.New(sprintservice.Deps{
		Repo:    sprintRepo,
		Project: projectSvc,
		Bus:     bus,
	})
	boardSvc := boardservice.New(boardservice.Deps{
		Repo:   boardRepo,
		Sprint: sprintSvc,
		Bus:    bus,
	})
	ticketSvc := ticketservice.New(ticketservice.Deps{
		Repo:    ticketRepo,
		Project: projectSvc,
		Board:   boardSvc,
		Sprint:  sprintSvc,
		Bus:     bus,
	})
	authSvc := authservice.New(authservice.Deps{
		Users:  userSvc,
		Config: &authCfg,
	})

	userC := usercache.New(memCache)
	orgC := orgcache.New(memCache)
	projectC := projectcache.New(memCache)
	sprintC := sprintcache.New(memCache)
	boardC := boardcache.New(memCache)
	ticketC := ticketcache.New(memCache)

	authH := authhandler.New(authSvc)
	userH := userhandler.New(userhandler.Deps{
		Svc:       userSvc,
		UserCache: userC,
	})
	orgH := orghandler.New(orghandler.Deps{
		Svc:     orgSvc,
		OrgCache: orgC,
	})
	projectH := projecthandler.New(projecthandler.Deps{
		Svc:          projectSvc,
		ProjectCache: projectC,
	})
	sprintH := sprinthandler.New(sprinthandler.Deps{
		Svc:        sprintSvc,
		SprintCache: sprintC,
	})
	boardH := boardhandler.New(boardhandler.Deps{
		Svc:        boardSvc,
		BoardCache: boardC,
	})
	ticketH := tickethandler.New(tickethandler.Deps{
		Svc:        ticketSvc,
		TicketCache: ticketC,
	})

	authModule := auth.NewModule(authSvc, authH, bus)
	userModule := user.NewModule(userH, userC, bus)
	orgModule := org.NewModule(orgH, orgC, bus)
	projectModule := project.NewModule(projectH, projectC, bus)
	sprintModule := sprint.NewModule(sprintH, sprintC, bus)
	boardModule := board.NewModule(boardH, boardC, bus)
	ticketModule := ticket.NewModule(ticketH, ticketC, bus)

	httpx.InitAuth(authModule.Service())

	mux := http.NewServeMux()
	authModule.Routes(mux)
	userModule.Routes(mux)
	orgModule.Routes(mux)
	projectModule.Routes(mux)
	sprintModule.Routes(mux)
	boardModule.Routes(mux)
	ticketModule.Routes(mux)

	testServer = httptest.NewServer(mux)
	defer testServer.Close()

	code := m.Run()
	os.Exit(code)
}
