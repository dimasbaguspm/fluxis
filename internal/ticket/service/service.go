package service

import (
	"github.com/dimasbaguspm/fluxis/internal/ticket/repository"
	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

type Deps struct {
	Repo    *repository.Queries
	Project domain.ProjectReader
	Board   domain.BoardReader
	Sprint  domain.SprintReader
}

type Service struct {
	Deps
}

var _ domain.TicketReader = (*Service)(nil)
var _ domain.TicketWriter = (*Service)(nil)

func New(d Deps) *Service {
	return &Service{d}
}
