package service

import "github.com/dimasbaguspm/fluxis/internal/org/repository"

type Deps struct {
	Repo *repository.Queries
}

type Service struct {
	Deps
}

func New(d Deps) *Service {
	return &Service{d}
}
