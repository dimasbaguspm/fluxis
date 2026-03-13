package org

import (
	"context"
	"log/slog"
	"net/http"

	orgcache "github.com/dimasbaguspm/fluxis/internal/org/cache"
	"github.com/dimasbaguspm/fluxis/internal/org/handler"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
	"github.com/dimasbaguspm/fluxis/pkg/pubsub"
)

type Module struct {
	h        *handler.Handler
	orgCache *orgcache.OrgCache
	bus      pubsub.Bus
}

func NewModule(h *handler.Handler, c *orgcache.OrgCache, bus pubsub.Bus) *Module {
	return &Module{h: h, orgCache: c, bus: bus}
}

func (m *Module) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /orgs", httpx.RequireAuth(m.h.ListOrgs))
	mux.HandleFunc("POST /orgs", httpx.RequireAuth(m.h.CreateOrg))
	mux.HandleFunc("GET /orgs/{id}", httpx.RequireAuth(m.h.GetOrg))
	mux.HandleFunc("PATCH /orgs/{id}", httpx.RequireAuth(m.h.UpdateOrg))
	mux.HandleFunc("DELETE /orgs/{id}", httpx.RequireAuth(m.h.DeleteOrg))
	mux.HandleFunc("GET /orgs/{id}/members", httpx.RequireAuth(m.h.ListOrgMembers))
	mux.HandleFunc("POST /orgs/{id}/members", httpx.RequireAuth(m.h.AddOrgMember))
	mux.HandleFunc("PATCH /orgs/{id}/members/{userId}", httpx.RequireAuth(m.h.UpdateOrgMember))
	mux.HandleFunc("DELETE /orgs/{id}/members/{userId}", httpx.RequireAuth(m.h.DeleteOrgMember))
}

func (m *Module) StartSubscriber(ctx context.Context) {
	slog.Info("[OrgModule]: starting bus subscriber")
	handler := func(ctx context.Context, e pubsub.Event) error {
		switch e.Type {
		case pubsub.OrgCreated, pubsub.OrgUpdated, pubsub.OrgDeleted:
			if orgID, ok := pubsub.UUIDFromPayload(e, "id"); ok {
				m.orgCache.InvalidateSingleOrg(ctx, orgID)
			}
			m.orgCache.InvalidatePagedOrganizations(ctx)
		}
		return nil
	}

	m.bus.Subscribe(ctx, pubsub.Channel(pubsub.Org), handler)
}
