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

var _ domain.SprintReader = (*Service)(nil)
var _ domain.SprintWriter = (*Service)(nil)

func New(d Deps) *Service {
	return &Service{d}
}
