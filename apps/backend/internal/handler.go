package internal

import (
	"context"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/dimasbaguspm/fluxis/internal/middlewares"
	"github.com/dimasbaguspm/fluxis/internal/repositories"
	"github.com/dimasbaguspm/fluxis/internal/resources"
	"github.com/dimasbaguspm/fluxis/internal/services"
	"github.com/dimasbaguspm/fluxis/internal/workers"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterPublicRoutes(api huma.API, pgx *pgxpool.Pool) {
	authRepo := repositories.NewAuthRepository(pgx)
	authSrv := services.NewAuthService(authRepo)

	resources.NewAuthResource(authSrv).Routes(api)
}

func RegisterPrivateRoutes(ctx context.Context, api huma.API, pgx *pgxpool.Pool) {
	api.UseMiddleware(middlewares.SessionMiddleware(api))

	prR := repositories.NewProjectRepository(pgx)
	sR := repositories.NewStatusRepository(pgx)
	tR := repositories.NewTaskRepository(pgx)
	lr := repositories.NewLogRepository(pgx)

	lW := workers.NewLogWorker(prR, sR, tR, lr, 10*time.Second)

	pS := services.NewProjectService(prR, lW, lr)
	sS := services.NewStatusService(sR, lW, lr)
	tS := services.NewTaskService(tR, prR, sR, lW, lr)

	resources.NewProjectResource(pS).Routes(api)
	resources.NewStatusResource(sS).Routes(api)
	resources.NewTaskResource(tS).Routes(api)

	go func() {
		<-ctx.Done()
		lW.Stop()
	}()

}
