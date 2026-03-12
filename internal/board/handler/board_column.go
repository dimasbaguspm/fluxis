package handler

import (
	"net/http"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
)

// ListBoardColumns godoc
//
//	@Summary		List board columns
//	@Description	Returns all columns in a board
//	@Tags			board
//	@Produce		json
//	@Param			boardId	path		string	true	"Board ID"
//	@Success		200		{array}		domain.BoardColumnModel
//	@Failure		400		{object}	httpx.ErrBlock
//	@Failure		401		{object}	httpx.ErrBlock
//	@Failure		404		{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/boards/{boardId}/columns [get]
func (h *Handler) ListBoardColumns(w http.ResponseWriter, r *http.Request) {
	boardID, err := httpx.PathUUID(r, "boardId")
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	columns, err := h.svc.ListBoardColumns(r.Context(), boardID)
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	httpx.OK(w, columns)
}

// CreateBoardColumn godoc
//
//	@Summary		Create a board column
//	@Description	Creates a new column in a board (position is auto-calculated)
//	@Tags			board
//	@Accept			json
//	@Produce		json
//	@Param			boardId	path		string								true	"Board ID"
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

// UpdateBoardColumn godoc
//
//	@Summary		Update a board column
//	@Description	Updates column name (use reorder endpoint for position changes)
//	@Tags			board
//	@Accept			json
//	@Produce		json
//	@Param			boardId			path		string								true	"Board ID"
//	@Param			boardColumnId	path		string								true	"Board Column ID"
//	@Param			body			body		domain.BoardColumnUpdateModel	true	"Column payload"
//	@Success		200				{object}	domain.BoardColumnModel
//	@Failure		400				{object}	httpx.ErrBlock
//	@Failure		401				{object}	httpx.ErrBlock
//	@Failure		404				{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/boards/{boardId}/columns/{boardColumnId} [patch]
func (h *Handler) UpdateBoardColumn(w http.ResponseWriter, r *http.Request) {
	boardID, err := httpx.PathUUID(r, "boardId")
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	columnID, err := httpx.PathUUID(r, "boardColumnId")
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	var req domain.BoardColumnUpdateModel
	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.Handle(w, httpx.BadRequest(err.Error()))
		return
	}

	column, err := h.svc.UpdateBoardColumn(r.Context(), boardID, columnID, req)
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	httpx.OK(w, column)
}

// ReorderBoardColumns godoc
//
//	@Summary		Reorder board columns
//	@Description	Reorder columns within a board
//	@Tags			board
//	@Accept			json
//	@Produce		json
//	@Param			boardId	path		string									true	"Board ID"
//	@Param			body		body		domain.BoardColumnReorderModel	true	"Column reorder payload"
//	@Success		200		{array}		domain.BoardColumnModel
//	@Failure		400		{object}	httpx.ErrBlock
//	@Failure		401		{object}	httpx.ErrBlock
//	@Failure		404		{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/boards/{boardId}/columns/reorder [patch]
func (h *Handler) ReorderBoardColumns(w http.ResponseWriter, r *http.Request) {
	boardID, err := httpx.PathUUID(r, "boardId")
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	var req domain.BoardColumnReorderModel
	if err := httpx.Decode(r, &req); err != nil {
		httpx.Handle(w, httpx.BadRequest(err.Error()))
		return
	}

	if len(req) == 0 {
		httpx.Handle(w, httpx.BadRequest("columns array is required and cannot be empty"))
		return
	}

	columns, err := h.svc.ReorderBoardColumns(r.Context(), boardID, req)
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	httpx.OK(w, columns)
}

// DeleteBoardColumn godoc
//
//	@Summary		Delete a board column
//	@Description	Deletes a column from a board
//	@Tags			board
//	@Produce		json
//	@Param			boardId			path		string	true	"Board ID"
//	@Param			boardColumnId	path		string	true	"Board Column ID"
//	@Success		204
//	@Failure		400				{object}	httpx.ErrBlock
//	@Failure		401				{object}	httpx.ErrBlock
//	@Failure		404				{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/boards/{boardId}/columns/{boardColumnId} [delete]
func (h *Handler) DeleteBoardColumn(w http.ResponseWriter, r *http.Request) {
	boardID, err := httpx.PathUUID(r, "boardId")
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	columnID, err := httpx.PathUUID(r, "boardColumnId")
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	if err := h.svc.DeleteBoardColumn(r.Context(), boardID, columnID); err != nil {
		httpx.Handle(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
