package handler

import (
	"net/http"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
)

// POST /auth/register
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req domain.AuthRegisterModel
	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.Handle(w, httpx.BadRequest(err.Error()))
		return
	}

	resp, err := h.svc.Register(r.Context(), req)
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	httpx.Created(w, resp)
}

// POST /auth/login
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.AuthLoginModel
	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.Handle(w, httpx.BadRequest(err.Error()))
		return
	}

	resp, err := h.svc.Login(r.Context(), req)
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	httpx.OK(w, resp)
}

// POST /auth/refresh
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req domain.AuthRefreshModel
	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.Handle(w, httpx.BadRequest(err.Error()))
		return
	}

	resp, err := h.svc.RotateAccessToken(r.Context(), req)
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	httpx.OK(w, resp)
}
