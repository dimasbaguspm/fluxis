package service

import (
	"context"
	"fmt"
	"time"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = httpx.Unauthorized("invalid email or password")
	ErrAccountLocked      = httpx.TooManyRequests("account temporarily locked, try again later")
	ErrUserAlreadyExists  = httpx.Conflict("email already registered").WithCode("email_taken")
	ErrTokenInvalid       = httpx.Unauthorized("token is invalid or expired")
)

func (s *Service) Register(ctx context.Context, p domain.AuthRegisterModel) (domain.AuthModel, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(p.Password), s.Config.BcryptCost)
	if err != nil {
		return domain.AuthModel{}, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.Users.CreateUser(ctx, domain.UserCreateModel{
		Email:       p.Email,
		Password:    string(hash),
		DisplayName: p.DisplayName,
	})
	if err != nil {
		return domain.AuthModel{}, err
	}

	tokens, err := s.generateTokens(ctx, user)
	if err != nil {
		return domain.AuthModel{}, err
	}
	return tokens, nil
}

func (s *Service) Login(ctx context.Context, p domain.AuthLoginModel) (domain.AuthModel, error) {
	user, err := s.Users.GetSingleUserByEmail(ctx, p.Email)
	if err != nil {
		return domain.AuthModel{}, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(p.Password)); err != nil {
		return domain.AuthModel{}, ErrInvalidCredentials
	}

	tokens, err := s.generateTokens(ctx, user)
	if err != nil {
		return domain.AuthModel{}, err
	}
	return tokens, nil
}

func (s *Service) RotateAccessToken(ctx context.Context, p domain.AuthRefreshModel) (domain.AuthModel, error) {
	now := time.Now()

	refreshClaim, err := s.validateRefreshToken(ctx, p.RefreshToken)
	if err != nil {
		return domain.AuthModel{}, ErrTokenInvalid
	}

	if refreshClaim.ExpiresAt.After(now) {
		return domain.AuthModel{}, ErrTokenInvalid
	}

	user, err := s.Users.GetSingleUserById(ctx, refreshClaim.ID)
	if err != nil {
		return domain.AuthModel{}, ErrInvalidCredentials
	}

	tokens, err := s.generateTokens(ctx, user)
	if err != nil {
		return domain.AuthModel{}, err
	}
	return tokens, nil
}
