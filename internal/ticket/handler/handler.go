package handler

import (
	"github.com/dimasbaguspm/fluxis/internal/ticket/cache"
	"github.com/dimasbaguspm/fluxis/internal/ticket/service"
)

type Deps struct {
	Svc        *service.Service
	TicketCache *cache.TicketCache
}

type Handler struct {
	svc        *service.Service
	ticketCache *cache.TicketCache
}

func New(deps Deps) *Handler {
	return &Handler{
		svc:        deps.Svc,
		ticketCache: deps.TicketCache,
	}
}
