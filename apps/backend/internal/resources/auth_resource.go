package resources

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/dimasbaguspm/fluxis/internal/configs"
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
		Method:      http.MethodPost,
		Path:        "/auth/login",
		Summary:     "Login",
		Description: "Authenticate with username and password to receive access and refresh tokens",
		Tags:        []string{"Authentication"},
	}, ar.login)
	huma.Register(api, huma.Operation{
		OperationID: "refresh",
		Method:      http.MethodPost,
		Path:        "/auth/refresh",
		Summary:     "Refresh Token",
		Description: "Exchange a valid refresh token for a new access token",
		Tags:        []string{"Authentication"},
	}, ar.refresh)
}

func (ar AuthResource) login(ctx context.Context, input *struct{ Body models.AuthLoginInputModel }) (*struct{ Body models.AuthLoginOutputModel }, error) {
	env := configs.NewEnvironment()
	svcResp, err := ar.authService.Login(input.Body, env)

	if err != nil {
		return nil, err
	}

	resp := &struct{ Body models.AuthLoginOutputModel }{
		Body: svcResp,
	}

	return resp, nil
}

func (ar AuthResource) refresh(ctx context.Context, input *struct{ Body models.AuthRefreshInputModel }) (*struct{ Body models.AuthRefreshOutputModel }, error) {
	svcResp, err := ar.authService.Refresh(input.Body)
	if err != nil {
		return nil, err
	}

	resp := &struct{ Body models.AuthRefreshOutputModel }{
		Body: svcResp,
	}

	return resp, err
}
