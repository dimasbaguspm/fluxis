package auth

import (
	"net/http"

	"github.com/dimasbaguspm/fluxis/internal/auth/handler"
)

type Module struct {
	h *handler.Handler
}

func NewModule(h *handler.Handler) *Module {
	return &Module{h}
}

func (m *Module) Routes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/register", m.h.Register)
	mux.HandleFunc("POST /auth/login", m.h.Login)
	mux.HandleFunc("POST /auth/refresh", m.h.Refresh)
}
