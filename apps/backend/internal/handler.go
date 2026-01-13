package internal

import (
	"context"

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

	pR := repositories.NewProjectRepository(pgx)
	sR := repositories.NewStatusRepository(pgx)
	tR := repositories.NewTaskRepository(pgx)
	lR := repositories.NewLogRepository(pgx)

	pW := workers.NewProjectWorker(ctx, pR, lR)
	sW := workers.NewStatusWorker(ctx, sR, lR)
	tW := workers.NewTaskWorker(ctx, tR, lR)

	pS := services.NewProjectService(pR, pW, lR)
	sS := services.NewStatusService(sR, sW, lR, pR)
	tS := services.NewTaskService(tR, pR, sR, tW, lR)

	resources.NewProjectResource(pS).Routes(api)
	resources.NewStatusResource(sS).Routes(api)
	resources.NewTaskResource(tS).Routes(api)

	go func() {
		<-ctx.Done()
		pW.Stop()
		sW.Stop()
		tW.Stop()
	}()

}
