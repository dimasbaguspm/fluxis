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
	orghandler "github.com/dimasbaguspm/fluxis/internal/org/handler"
	orgrepo "github.com/dimasbaguspm/fluxis/internal/org/repository"
	orgservice "github.com/dimasbaguspm/fluxis/internal/org/service"

	"github.com/dimasbaguspm/fluxis/internal/user"
	userhandler "github.com/dimasbaguspm/fluxis/internal/user/handler"
	userrepo "github.com/dimasbaguspm/fluxis/internal/user/repository"
	userservice "github.com/dimasbaguspm/fluxis/internal/user/service"
	"github.com/dimasbaguspm/fluxis/pkg/testutil"

	"github.com/dimasbaguspm/fluxis/pkg/httpx"
)

var testServer *httptest.Server

func TestMain(m *testing.M) {
	ctx := context.Background()

	pgContainer, err := testutil.NewPostgresContainer(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start postgres container: %v\n", err)
		os.Exit(1)
	}
	defer pgContainer.Terminate(ctx)

	pool := testutil.MustPool(ctx, pgContainer.DSN)
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

	if err := testutil.RunMigrations(pgContainer.DSN, migrationsPath); err != nil {
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

	userSvc := userservice.New(userservice.Deps{
		Repo: userRepo,
	})
	orgSvc := orgservice.New(orgservice.Deps{
		Repo: orgRepo,
	})
	authSvc := authservice.New(authservice.Deps{
		Users:  userSvc,
		Config: &authCfg,
	})

	authH := authhandler.New(authSvc)
	userH := userhandler.New(userSvc)
	orgH := orghandler.New(orgSvc)

	authModule := auth.NewModule(authSvc, authH)
	userModule := user.NewModule(userH)
	orgModule := org.NewModule(orgH)

	httpx.InitAuth(authModule.Service())

	mux := http.NewServeMux()
	authModule.Routes(mux)
	userModule.Routes(mux)
	orgModule.Routes(mux)

	testServer = httptest.NewServer(mux)
	defer testServer.Close()

	code := m.Run()
	os.Exit(code)
}
