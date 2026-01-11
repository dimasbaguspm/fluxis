package resources

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/dimasbaguspm/fluxis/internal/models"
	"github.com/dimasbaguspm/fluxis/internal/services"
)

type AuthResource struct {
	authService services.AuthService
}

func NewAuthResource(authService services.AuthService) AuthResource {
	return AuthResource{authService}
}

func (ar AuthResource) Routes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "login",
		Method:      http.MethodGet,
		Path:        "/auth/login",
		Summary:     "Login",
	}, ar.login)
	huma.Register(api, huma.Operation{
		OperationID: "refresh",
		Method:      http.MethodPost,
		Path:        "/auth/refresh",
		Summary:     "Place where to get the valid access token",
	}, ar.refresh)
}

func (ar AuthResource) login(ctx context.Context, input *models.AuthLoginInputModel) (*models.AuthLoginOutputModel, error) {

}

func (ar AuthResource) refresh(ctx context.Context, input *models.AuthRefreshInputModel) (*models.AuthRefreshOutputModel, error) {
	//
}
