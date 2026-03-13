package handler

import (
	"github.com/dimasbaguspm/fluxis/internal/project/cache"
	"github.com/dimasbaguspm/fluxis/internal/project/service"
)

type Deps struct {
	Svc          *service.Service
	ProjectCache *cache.ProjectCache
}

type Handler struct {
	svc          *service.Service
	projectCache *cache.ProjectCache
}

func New(deps Deps) *Handler {
	return &Handler{
		svc:          deps.Svc,
		projectCache: deps.ProjectCache,
	}
}
