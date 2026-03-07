package domain

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type AuthModel struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type AuthRegisterModel struct {
	UserCreateModel
}

type AuthLoginModel struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required"`
}

type AuthRefreshModel struct {
	AccessToken  string `json:"accessToken" validate:"required"`
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type AuthTokenClaimModel struct {
	jwt.RegisteredClaims
	ID pgtype.UUID `json:"id" validate:"uuid4"`
}

type AuthWrite interface {
	Register(ctx context.Context, p AuthRegisterModel) (AuthModel, error)
	Login(ctx context.Context, p AuthLoginModel) (AuthModel, error)
	RotateAccessToken(ctx context.Context, p AuthRefreshModel) (AuthModel, error)
	GenerateTokens(ctx context.Context, p UserModel) (AuthModel, error)
	ValidateAccessToken(ctx context.Context, tokenStr string) (AuthTokenClaimModel, error)
	ValidateRefreshToken(ctx context.Context, tokenStr string) (AuthTokenClaimModel, error)
}
