package service

import (
	"github.com/dimasbaguspm/fluxis/internal/board/repository"
)

type Deps struct {
	Repo *repository.Queries
}

type Service struct {
	repo *repository.Queries
}

func New(d Deps) *Service {
	return &Service{
		repo: d.Repo,
	}
}
