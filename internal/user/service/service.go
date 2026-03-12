package service

import (
	"github.com/dimasbaguspm/fluxis/internal/user/repository"
	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

type Deps struct {
	Repo *repository.Queries
}

type Service struct {
	Deps
}

var _ domain.UserRead = (*Service)(nil)
var _ domain.UserWrite = (*Service)(nil)

func New(d Deps) *Service {
	return &Service{d}
}
