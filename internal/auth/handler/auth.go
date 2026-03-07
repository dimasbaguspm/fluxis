package handler

import (
	"net/http"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
)

// Register godoc
//
//	@Summary		Register a new user
//	@Description	Creates a new user account and returns access/refresh tokens
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		domain.AuthRegisterModel	true	"Registration payload"
//	@Success		201		{object}	domain.AuthModel
//	@Failure		400		{object}	httpx.ErrBlock
//	@Router			/auth/register [post]
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

// Login godoc
//
//	@Summary		Login with email and password
//	@Description	Authenticates a user and returns access/refresh tokens
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		domain.AuthLoginModel	true	"Login payload"
//	@Success		200		{object}	domain.AuthModel
//	@Failure		400		{object}	httpx.ErrBlock
//	@Failure		401		{object}	httpx.ErrBlock
//	@Router			/auth/login [post]
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

// Refresh godoc
//
//	@Summary		Rotate access token
//	@Description	Issues a new access token using a valid refresh token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		domain.AuthRefreshModel	true	"Refresh payload"
//	@Success		200		{object}	domain.AuthModel
//	@Failure		400		{object}	httpx.ErrBlock
//	@Failure		401		{object}	httpx.ErrBlock
//	@Router			/auth/refresh [post]
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
