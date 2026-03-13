package project

import (
	"context"
	"log/slog"
	"net/http"

	projectcache "github.com/dimasbaguspm/fluxis/internal/project/cache"
	"github.com/dimasbaguspm/fluxis/internal/project/handler"
	"github.com/dimasbaguspm/fluxis/pkg/domain"
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
		var project domain.ProjectModel
		if err := httpx.DecodePayload(e.Payload, &project); err != nil {
			return nil
		}

		switch e.Type {
		case pubsub.ProjectCreated, pubsub.ProjectUpdated, pubsub.ProjectDeleted, pubsub.ProjectVisibilityUpdated:
			m.projectCache.InvalidateSingleProject(ctx, project.ID)
			m.projectCache.InvalidateSingleProjectByKey(ctx, project.OrgID, project.Key)
			m.projectCache.InvalidatePagedProjects(ctx)
		}
		return nil
	}

	m.bus.Subscribe(ctx, pubsub.Channel(pubsub.Project), handler)
}
