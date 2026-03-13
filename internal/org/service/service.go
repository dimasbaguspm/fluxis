package service

import (
	"github.com/dimasbaguspm/fluxis/internal/org/repository"
	"github.com/dimasbaguspm/fluxis/pkg/domain"
	"github.com/dimasbaguspm/fluxis/pkg/pubsub"
)

type Deps struct {
	Repo   *repository.Queries
	User   domain.UserRead
	Bus    pubsub.Publisher
}

type Service struct {
	Deps
}

var _ domain.OrgReader = (*Service)(nil)
var _ domain.OrganisationWrite = (*Service)(nil)

func New(d Deps) *Service {
	return &Service{d}
}
