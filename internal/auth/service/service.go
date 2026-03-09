package service

import (
	"time"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

type Deps struct {
	Users interface {
		domain.UserRead
		domain.UserWrite
	}
	Config *Config
}

type Service struct {
	Deps
}

type Config struct {
	AccessTokenSecret  string        // H256
	RefreshTokenSecret string        // H256
	AccessTokenExpiry  time.Duration // default 15m
	RefreshTokenExpiry time.Duration // default 7d

	BcryptCost int
}

func New(d Deps) *Service {
	return &Service{d}
}
