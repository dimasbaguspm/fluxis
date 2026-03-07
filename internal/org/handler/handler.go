package handler

import "github.com/dimasbaguspm/fluxis/internal/org/service"

type Handler struct {
	svc *service.Service
}

func New(svc *service.Service) *Handler {
	return &Handler{svc}
}
