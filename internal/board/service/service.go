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

var _ domain.BoardReader = (*Service)(nil)
var _ domain.BoardWriter = (*Service)(nil)

func New(d Deps) *Service {
	return &Service{d}
}
