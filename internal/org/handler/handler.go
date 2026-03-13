package handler

import (
	"github.com/dimasbaguspm/fluxis/internal/org/cache"
	"github.com/dimasbaguspm/fluxis/internal/org/service"
)

type Deps struct {
	Svc     *service.Service
	OrgCache *cache.OrgCache
}

type Handler struct {
	svc     *service.Service
	orgCache *cache.OrgCache
}

func New(deps Deps) *Handler {
	return &Handler{
		svc:     deps.Svc,
		orgCache: deps.OrgCache,
	}
}
