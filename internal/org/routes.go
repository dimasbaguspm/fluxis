package org

import (
	"net/http"

	"github.com/dimasbaguspm/fluxis/internal/org/handler"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
)

type Module struct {
	h *handler.Handler
}

func NewModule(h *handler.Handler) *Module {
	return &Module{h}
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
