package notification

import (
	"net/http"

	"github.com/dimasbaguspm/fluxis/internal/notification/handler"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
)

type Module struct {
	h *handler.Handler
}

type Deps struct {
	Handler *handler.Handler
}

func NewModule(h *handler.Handler) *Module {
	return &Module{
		h: h,
	}
}

func (m *Module) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /notifications", httpx.RequireAuth(m.h.Stream))
}
