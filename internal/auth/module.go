package auth

import (
	"net/http"

	"github.com/dimasbaguspm/fluxis/internal/auth/handler"
	"github.com/dimasbaguspm/fluxis/internal/auth/service"
)

type Module struct {
	h   *handler.Handler
	svc *service.Service
}

func NewModule(svc *service.Service, h *handler.Handler) *Module {
	return &Module{svc: svc, h: h}
}

func (m *Module) Routes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/register", m.h.Register)
	mux.HandleFunc("POST /auth/login", m.h.Login)
	mux.HandleFunc("POST /auth/refresh", m.h.Refresh)
}

func (m *Module) Service() *service.Service {
	return m.svc
}
