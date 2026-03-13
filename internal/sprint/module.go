package sprint

import (
	"net/http"

	"github.com/dimasbaguspm/fluxis/internal/sprint/handler"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
)

type Module struct {
	h *handler.Handler
}

func NewModule(h *handler.Handler) *Module {
	return &Module{h}
}

func (m *Module) Routes(mux *http.ServeMux) {
	mux.HandleFunc("POST /sprints", httpx.RequireAuth(m.h.CreateSprint))
	mux.HandleFunc("GET /sprints", httpx.RequireAuth(m.h.ListSprints))
	mux.HandleFunc("GET /sprints/{sprintId}", httpx.RequireAuth(m.h.GetSprint))
	mux.HandleFunc("PATCH /sprints/{sprintId}", httpx.RequireAuth(m.h.UpdateSprint))
	mux.HandleFunc("POST /sprints/{sprintId}/start", httpx.RequireAuth(m.h.StartSprint))
	mux.HandleFunc("POST /sprints/{sprintId}/completed", httpx.RequireAuth(m.h.CompleteSprint))
}
