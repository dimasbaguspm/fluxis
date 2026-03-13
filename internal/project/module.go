package project

import (
	"context"
	"log/slog"
	"net/http"

	projectcache "github.com/dimasbaguspm/fluxis/internal/project/cache"
	"github.com/dimasbaguspm/fluxis/internal/project/handler"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
	"github.com/dimasbaguspm/fluxis/pkg/pubsub"
)

type Module struct {
	h              *handler.Handler
	projectCache   *projectcache.ProjectCache
	bus            pubsub.Bus
}

func NewModule(h *handler.Handler, c *projectcache.ProjectCache, bus pubsub.Bus) *Module {
	return &Module{
		h:            h,
		projectCache: c,
		bus:          bus,
	}
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
		switch e.Type {
		case pubsub.ProjectUpdated, pubsub.ProjectDeleted, pubsub.ProjectVisibilityUpdated:
			if projectID, ok := pubsub.UUIDFromPayload(e, "id"); ok {
				m.projectCache.InvalidateSingleProject(ctx, projectID)
			}
			if orgID, ok := pubsub.UUIDFromPayload(e, "orgId"); ok {
				if key, ok := pubsub.StringFromPayload(e, "key"); ok {
					m.projectCache.InvalidateSingleProjectByKey(ctx, orgID, key)
				}
			}
		}
		return nil
	}

	m.bus.Subscribe(ctx, pubsub.Channel(pubsub.Project), handler)
}
