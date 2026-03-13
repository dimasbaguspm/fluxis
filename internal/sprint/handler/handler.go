package handler

import (
	"github.com/dimasbaguspm/fluxis/internal/sprint/cache"
	"github.com/dimasbaguspm/fluxis/internal/sprint/service"
)

type Deps struct {
	Svc        *service.Service
	SprintCache *cache.SprintCache
}

type Handler struct {
	svc        *service.Service
	sprintCache *cache.SprintCache
}

func New(deps Deps) *Handler {
	return &Handler{
		svc:        deps.Svc,
		sprintCache: deps.SprintCache,
	}
}
