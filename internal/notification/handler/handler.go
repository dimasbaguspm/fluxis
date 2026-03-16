package handler

import "github.com/dimasbaguspm/fluxis/internal/notification/service"

type Deps struct {
	Service *service.Service
}

type Handler struct {
	svc *service.Service
}

func New(deps Deps) *Handler {
	return &Handler{
		svc: deps.Service,
	}
}
