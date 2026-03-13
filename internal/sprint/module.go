package sprint

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/dimasbaguspm/fluxis/internal/sprint/handler"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
	"github.com/dimasbaguspm/fluxis/pkg/pubsub"
)

type Module struct {
	h   *handler.Handler
	bus pubsub.Bus
}

func NewModule(h *handler.Handler, bus pubsub.Bus) *Module {
	return &Module{h, bus}
}

func (m *Module) Routes(mux *http.ServeMux) {
	mux.HandleFunc("POST /sprints", httpx.RequireAuth(m.h.CreateSprint))
	mux.HandleFunc("GET /sprints", httpx.RequireAuth(m.h.ListSprints))
	mux.HandleFunc("GET /sprints/{sprintId}", httpx.RequireAuth(m.h.GetSprint))
	mux.HandleFunc("PATCH /sprints/{sprintId}", httpx.RequireAuth(m.h.UpdateSprint))
	mux.HandleFunc("POST /sprints/{sprintId}/start", httpx.RequireAuth(m.h.StartSprint))
	mux.HandleFunc("POST /sprints/{sprintId}/completed", httpx.RequireAuth(m.h.CompleteSprint))
}

func (m *Module) StartSubscriber(ctx context.Context) {
	slog.Info("[SprintModule]: starting bus subscriber")
	handler := func(ctx context.Context, e pubsub.Event) error {
		slog.Info("[SprintModule]: received event", "type", string(e.Type), "payload", e.Payload)
		return nil
	}

	m.bus.Subscribe(ctx, pubsub.Channel(pubsub.Sprint), handler)
}
