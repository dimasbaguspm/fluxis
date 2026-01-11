package repositories

import (
	"errors"
	"time"

	"github.com/dimasbaguspm/fluxis/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	pgx *pgxpool.Pool
}

func NewAuthRepository(pgx *pgxpool.Pool) AuthRepository {
	return AuthRepository{pgx}
}

var (
	AuthErrorInvalidSigningMethod = errors.New("Invalid signin method")
	AuthErrorInvalidToken         = errors.New("Token is invalid")
)

const (
	accessTokenType  = "access"
	refreshTokenType = "refresh"
)

const secretJWT = "some-random-things-that-soon-will-be-replaced"

func (ar AuthRepository) GenerateFreshTokens(m models.AuthLoginInputModel) (accessToken, refreshToken string, err error) {
	accessToken, err = generateToken(accessTokenType)
	if err != nil {
		return "", "", err
	}
	refreshToken, err = generateToken(refreshTokenType)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (ar AuthRepository) RegenerateAccessToken(refreshToken string) (string, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, AuthErrorInvalidSigningMethod
		}
		return []byte(secretJWT), nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", AuthErrorInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["sub"] != refreshTokenType {
		return "", AuthErrorInvalidToken
	}

	return generateToken(accessTokenType)
}

func generateToken(sub string) (string, error) {
	now := time.Now()

	var subject string
	var expiredAt time.Time

	switch sub {
	case accessTokenType:
		subject = accessTokenType
		expiredAt = now.Add(7 * 24 * time.Hour)
	case refreshTokenType:
		subject = refreshTokenType
		expiredAt = now.Add(30 * 24 * time.Hour)
	}

	accessClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiredAt),
		IssuedAt:  jwt.NewNumericDate(now),
		Subject:   subject,
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	return accessToken.SignedString([]byte(secretJWT))
}
