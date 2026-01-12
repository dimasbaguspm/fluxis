package internal

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/dimasbaguspm/fluxis/internal/middlewares"
	"github.com/dimasbaguspm/fluxis/internal/repositories"
	"github.com/dimasbaguspm/fluxis/internal/resources"
	"github.com/dimasbaguspm/fluxis/internal/services"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterPublicRoutes(api huma.API, pgx *pgxpool.Pool) {
	authRepo := repositories.NewAuthRepository(pgx)
	authSrv := services.NewAuthService(authRepo)

	resources.NewAuthResource(authSrv).Routes(api)
}

func RegisterPrivateRoutes(api huma.API, pgx *pgxpool.Pool) {
	api.UseMiddleware(middlewares.SessionMiddleware(api))

	projectRepo := repositories.NewProjectRepository(pgx)
	statusRepo := repositories.NewStatusRepository(pgx)

	projectSrv := services.NewProjectService(projectRepo)
	statusSrv := services.NewStatusService(statusRepo)

	resources.NewProjectResource(projectSrv).Routes(api)
	resources.NewStatusResource(statusSrv).Routes(api)
}
