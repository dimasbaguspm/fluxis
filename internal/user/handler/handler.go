package handler

import (
	"github.com/dimasbaguspm/fluxis/internal/user/cache"
	"github.com/dimasbaguspm/fluxis/internal/user/service"
)

type Deps struct {
	Svc       *service.Service
	UserCache *cache.UserCache
}

type Handler struct {
	svc       *service.Service
	userCache *cache.UserCache
}

func New(deps Deps) *Handler {
	return &Handler{
		svc:       deps.Svc,
		userCache: deps.UserCache,
	}
}
