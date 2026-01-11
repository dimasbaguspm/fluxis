package services

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/dimasbaguspm/fluxis/internal/configs"
	"github.com/dimasbaguspm/fluxis/internal/models"
	"github.com/dimasbaguspm/fluxis/internal/repositories"
)

type AuthService struct {
	authRepo repositories.AuthRepository
}

func NewAuthService(authRepo repositories.AuthRepository) AuthService {
	return AuthService{authRepo}
}

func (as *AuthService) Login(data models.AuthLoginInputModel, env configs.Environment) (models.AuthLoginOutputModel, error) {
	isValid := env.Admin.Username != "" && env.Admin.Password != "" && data.Username == env.Admin.Username && data.Password == env.Admin.Password

	if !isValid {
		return models.AuthLoginOutputModel{}, huma.Error401Unauthorized("Invalid credentials")
	}

	accessToken, refreshToken, err := as.authRepo.GenerateFreshTokens(data)

	if err != nil {
		return models.AuthLoginOutputModel{}, err
	}

	return models.AuthLoginOutputModel{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Username:     data.Username,
	}, nil
}

func (as *AuthService) Refresh(data models.AuthRefreshInputModel) (models.AuthRefreshOutputModel, error) {
	newAccessToken, err := as.authRepo.RegenerateAccessToken(data.RefreshToken)

	if err != nil {
		return models.AuthRefreshOutputModel{}, err
	}

	return models.AuthRefreshOutputModel{
		AccessToken: newAccessToken,
	}, nil
}
