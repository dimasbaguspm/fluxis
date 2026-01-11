package services

import (
	"github.com/dimasbaguspm/fluxis/internal/models"
	"github.com/dimasbaguspm/fluxis/internal/repositories"
)

type AuthService struct {
	authRepo repositories.AuthRepository
}

func NewAuthService(authRepo repositories.AuthRepository) AuthService {
	return AuthService{authRepo}
}

func (as *AuthService) Login(data models.AuthLoginInputModel) (models.AuthLoginOutputModel, error) {
	// return as.authRepo
}

func (as *AuthService) Refresh(data models.AuthRefreshInputModel) (models.AuthRefreshOutputModel, error) {
	//
}
