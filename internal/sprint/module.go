package sprint

import (
	"context"
	"log/slog"
	"net/http"

	sprintcache "github.com/dimasbaguspm/fluxis/internal/sprint/cache"
	"github.com/dimasbaguspm/fluxis/internal/sprint/handler"
	"github.com/dimasbaguspm/fluxis/pkg/domain"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
	"github.com/dimasbaguspm/fluxis/pkg/pubsub"
)

type Module struct {
	h            *handler.Handler
	sprintCache  *sprintcache.SprintCache
	bus          pubsub.Bus
}

func NewModule(h *handler.Handler, c *sprintcache.SprintCache, bus pubsub.Bus) *Module {
	return &Module{
		h:           h,
		sprintCache: c,
		bus:         bus,
	}
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
		var sprint domain.SprintModel
		if err := httpx.DecodePayload(e.Payload, &sprint); err != nil {
			return nil
		}

		switch e.Type {
		case pubsub.SprintCreated, pubsub.SprintUpdated, pubsub.SprintStarted, pubsub.SprintCompleted:
			m.sprintCache.InvalidateSingleActiveSprint(ctx, sprint.ProjectID)
			m.sprintCache.InvalidateSingleSprint(ctx, sprint.ID)
			m.sprintCache.InvalidatePagedSprints(ctx)
		}
		return nil
	}

	m.bus.Subscribe(ctx, pubsub.Channel(pubsub.Sprint), handler)
}
