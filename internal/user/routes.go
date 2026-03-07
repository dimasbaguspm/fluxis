package user

import (
	"net/http"

	"github.com/dimasbaguspm/fluxis/internal/user/handler"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
)

type Module struct {
	h *handler.Handler
}

func NewModule(h *handler.Handler) *Module {
	return &Module{h}
}

func (m *Module) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /users/me", httpx.RequireAuth(m.h.GetCurrentUser))
}
