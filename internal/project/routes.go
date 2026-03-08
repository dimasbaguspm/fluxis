package project

import (
	"net/http"

	"github.com/dimasbaguspm/fluxis/internal/project/handler"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
)

type Module struct {
	h *handler.Handler
}

func NewModule(h *handler.Handler) *Module {
	return &Module{h}
}

func (m *Module) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /projects", httpx.RequireAuth(m.h.ListProjects))
	mux.HandleFunc("POST /projects", httpx.RequireAuth(m.h.CreateProject))
	mux.HandleFunc("GET /projects/{id}", httpx.RequireAuth(m.h.GetProject))
	mux.HandleFunc("PATCH /projects/{id}", httpx.RequireAuth(m.h.UpdateProject))
	mux.HandleFunc("PATCH /projects/{id}/visibility", httpx.RequireAuth(m.h.UpdateProjectVisibility))
	mux.HandleFunc("DELETE /projects/{id}", httpx.RequireAuth(m.h.DeleteProject))
}
