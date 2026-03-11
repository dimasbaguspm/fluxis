package service

import (
	"github.com/dimasbaguspm/fluxis/internal/sprint/repository"
	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

type Deps struct {
	Repo    *repository.Queries
	Project domain.ProjectReader
}

type Service struct {
	Deps
}

func New(d Deps) *Service {
	return &Service{d}
}
