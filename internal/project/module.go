package project

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/dimasbaguspm/fluxis/internal/project/handler"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
	"github.com/dimasbaguspm/fluxis/pkg/pubsub"
)

type Module struct {
	h   *handler.Handler
	bus pubsub.Bus
}

func NewModule(h *handler.Handler, bus pubsub.Bus) *Module {
	return &Module{h: h, bus: bus}
}

func (m *Module) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /projects", httpx.RequireAuth(m.h.ListProjects))
	mux.HandleFunc("POST /projects", httpx.RequireAuth(m.h.CreateProject))
	mux.HandleFunc("GET /projects/{id}", httpx.RequireAuth(m.h.GetProject))
	mux.HandleFunc("PATCH /projects/{id}", httpx.RequireAuth(m.h.UpdateProject))
	mux.HandleFunc("PATCH /projects/{id}/visibility", httpx.RequireAuth(m.h.UpdateProjectVisibility))
	mux.HandleFunc("DELETE /projects/{id}", httpx.RequireAuth(m.h.DeleteProject))
}

func (m *Module) StartSubscriber(ctx context.Context) {
	slog.Info("[ProjectModule]: starting bus subscriber")
	handler := func(ctx context.Context, e pubsub.Event) error {
		slog.Info("[ProjectModule]: received event", "type", string(e.Type), "payload", e.Payload)
		return nil
	}

	m.bus.Subscribe(ctx, pubsub.Channel(pubsub.Project), handler)
}
