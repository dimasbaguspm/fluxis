package internal

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/dimasbaguspm/fluxis/internal/repositories"
	"github.com/dimasbaguspm/fluxis/internal/resources"
	"github.com/dimasbaguspm/fluxis/internal/services"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterPublicRoutes(api huma.API, pgx *pgxpool.Pool) {
	authRepo := repositories.NewAuthRepository(pgx)
	authService := services.NewAuthService(authRepo)

	resources.NewAuthResource(authService).Routes(api)
}

func RegisterPrivateRoutes(r *http.ServeMux, api huma.API, pgx *pgxpool.Pool) {

}
