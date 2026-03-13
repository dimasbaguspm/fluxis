package org

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/dimasbaguspm/fluxis/internal/org/handler"
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
		slog.Info("[OrgModule]: received event", "type", string(e.Type), "payload", e.Payload)
		return nil
	}

	m.bus.Subscribe(ctx, pubsub.Channel(pubsub.Org), handler)
}
