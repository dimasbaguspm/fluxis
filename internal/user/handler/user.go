package handler

import (
	"net/http"

	"github.com/dimasbaguspm/fluxis/pkg/httpx"
)

// GetCurrentUser godoc
//
//	@Summary		Get current user
//	@Description	Returns the authenticated user's profile
//	@Tags			user
//	@Produce		json
//	@Success		200	{object}	domain.UserModel
//	@Failure		401	{object}	httpx.ErrBlock
//	@Failure		404	{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/users/me [get]
func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := httpx.MustUserID(r.Context())

	user, err := h.svc.GetSingleUserById(r.Context(), userID)
	if err != nil {
		httpx.Handle(w, err)
		return
	}
	httpx.OK(w, user)
}
