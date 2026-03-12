package service

import (
	"github.com/dimasbaguspm/fluxis/internal/project/repository"
	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

type Deps struct {
	Repo *repository.Queries
	Org  domain.OrgReader
}

type Service struct {
	Deps
}

var _ domain.ProjectReader = (*Service)(nil)
var _ domain.ProjectWriter = (*Service)(nil)

func New(d Deps) *Service {
	return &Service{d}
}
