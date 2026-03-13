package handler

import (
	"github.com/dimasbaguspm/fluxis/internal/board/cache"
	"github.com/dimasbaguspm/fluxis/internal/board/service"
)

type Deps struct {
	Svc       *service.Service
	BoardCache *cache.BoardCache
}

type Handler struct {
	svc       *service.Service
	boardCache *cache.BoardCache
}

func New(deps Deps) *Handler {
	return &Handler{
		svc:        deps.Svc,
		boardCache: deps.BoardCache,
	}
}
