package domain

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type AuthModel struct {
	AccessToken  string `json:"accessToken"  example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refreshToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type AuthRegisterModel struct {
	UserCreateModel
}

type AuthLoginModel struct {
	Email    string `json:"email"    validate:"email,required" example:"user@example.com"`
	Password string `json:"password" validate:"required"        example:"s3cr3tP@ssword"`
}

type AuthRefreshModel struct {
	AccessToken  string `json:"accessToken"  validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refreshToken" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
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
