package service

import (
	"context"
	"fmt"
	"time"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrUnableToSignToken  = httpx.Unauthorized("token unable to be signed")
	ErrUnableToParseToken = httpx.Unauthorized("token unable to be parsed")
)

func (s *Service) generateTokens(_ context.Context, p domain.UserModel) (domain.AuthModel, error) {
	now := time.Now()
	accessExpiry := now.Add(s.Config.AccessTokenExpiry)
	refreshExpiry := now.Add(s.Config.RefreshTokenExpiry)

	accessClaims := domain.AuthTokenClaimModel{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   p.Email,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(accessExpiry),
		},
		ID: p.ID,
	}

	acessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(s.Config.AccessTokenSecret))
	if err != nil {
		return domain.AuthModel{}, ErrUnableToSignToken
	}

	refreshClaims := domain.AuthTokenClaimModel{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   p.Email,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(refreshExpiry),
		},
		ID: p.ID,
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(s.Config.RefreshTokenSecret))
	if err != nil {
		return domain.AuthModel{}, ErrUnableToSignToken
	}

	return domain.AuthModel{
		AccessToken:  acessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) validateRefreshToken(ctx context.Context, tokenstr string) (domain.AuthTokenClaimModel, error) {
	var claims domain.AuthTokenClaimModel
	_, err := jwt.ParseWithClaims(tokenstr, &claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(s.Config.RefreshTokenSecret), nil
	})
	if err != nil {
		return domain.AuthTokenClaimModel{}, ErrUnableToParseToken
	}
	return claims, nil
}
