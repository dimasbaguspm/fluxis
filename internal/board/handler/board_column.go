package handler

import (
	"net/http"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
)

// CreateBoardColumn godoc
//
//	@Summary		Create a board column
//	@Description	Creates a new column in a board
//	@Tags			board
//	@Accept			json
//	@Produce		json
//	@Param			boardId	query		string								true	"Board ID"
//	@Param			body		body		domain.BoardColumnCreateModel	true	"Column payload"
//	@Success		201			{object}	domain.BoardColumnModel
//	@Failure		400			{object}	httpx.ErrBlock
//	@Failure		401			{object}	httpx.ErrBlock
//	@Failure		404			{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/boards/{boardId}/columns [post]
func (h *Handler) CreateBoardColumn(w http.ResponseWriter, r *http.Request) {
	boardID, err := httpx.PathUUID(r, "boardId")
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	var req domain.BoardColumnCreateModel
	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.Handle(w, httpx.BadRequest(err.Error()))
		return
	}

	column, err := h.svc.CreateBoardColumn(r.Context(), boardID, req)
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	httpx.Created(w, column)
}
