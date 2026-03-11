package service

import (
	"github.com/dimasbaguspm/fluxis/internal/board/repository"
	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

type Deps struct {
	Repo   *repository.Queries
	Sprint domain.SprintReader
}

type Service struct {
	Deps
}

func New(d Deps) *Service {
	return &Service{d}
}
