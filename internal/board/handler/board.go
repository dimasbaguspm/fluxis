package handler

import (
	"net/http"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
)

// CreateBoard godoc
//
//	@Summary		Create a board
//	@Description	Creates a new board in a sprint
//	@Tags			board
//	@Accept			json
//	@Produce		json
//	@Param			body		body		domain.BoardCreateModel		true	"Board payload"
//	@Success		201			{object}	domain.BoardModel
//	@Failure		400			{object}	httpx.ErrBlock
//	@Failure		401			{object}	httpx.ErrBlock
//	@Failure		404			{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/boards [post]
func (h *Handler) CreateBoard(w http.ResponseWriter, r *http.Request) {
	var req domain.BoardCreateModel

	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.Handle(w, httpx.BadRequest(err.Error()))
		return
	}

	board, err := h.svc.CreateBoard(r.Context(), req)
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	httpx.Created(w, board)
}

// ListBoards godoc
//
//	@Summary		List boards
//	@Description	Returns all boards in a sprint
//	@Tags			board
//	@Produce		json
//	@Param			sprintId	query		string	true	"Sprint ID"
//	@Success		200			{array}		domain.BoardModel
//	@Failure		400			{object}	httpx.ErrBlock
//	@Failure		401			{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/boards [get]
func (h *Handler) ListBoards(w http.ResponseWriter, r *http.Request) {
	sprintID, err := httpx.QueryUUID(r, "sprintId")
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	boards, err := h.svc.ListBoardsBySprint(r.Context(), sprintID)
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	httpx.OK(w, boards)
}

// GetBoard godoc
//
//	@Summary		Get a board
//	@Description	Returns a single board by ID
//	@Tags			board
//	@Produce		json
//	@Param			boardId	path		string	true	"Board ID"
//	@Success		200		{object}	domain.BoardModel
//	@Failure		400		{object}	httpx.ErrBlock
//	@Failure		401		{object}	httpx.ErrBlock
//	@Failure		404		{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/boards/{boardId} [get]
func (h *Handler) GetBoard(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathUUID(r, "boardId")
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	board, err := h.svc.GetBoard(r.Context(), id)
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	httpx.OK(w, board)
}

// UpdateBoard godoc
//
//	@Summary		Update a board
//	@Description	Updates board details
//	@Tags			board
//	@Accept			json
//	@Produce		json
//	@Param			boardId	path		string						true	"Board ID"
//	@Param			body		body		domain.BoardUpdateModel	true	"Board payload"
//	@Success		200		{object}	domain.BoardModel
//	@Failure		400		{object}	httpx.ErrBlock
//	@Failure		401		{object}	httpx.ErrBlock
//	@Failure		404		{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/boards/{boardId} [patch]
func (h *Handler) UpdateBoard(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathUUID(r, "boardId")
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	var req domain.BoardUpdateModel

	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.Handle(w, httpx.BadRequest(err.Error()))
		return
	}

	board, err := h.svc.UpdateBoard(r.Context(), id, req)
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	httpx.OK(w, board)
}

// ReorderBoard godoc
//
//	@Summary		Reorder boards
//	@Description	Reorder boards within a sprint
//	@Tags			board
//	@Accept			json
//	@Produce		json
//	@Param			body		body		domain.BoardReorderModel	true	"Board reorder payload"
//	@Success		200		{array}		domain.BoardModel
//	@Failure		400		{object}	httpx.ErrBlock
//	@Failure		401		{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/boards/reorder [patch]
func (h *Handler) ReorderBoards(w http.ResponseWriter, r *http.Request) {
	var req domain.BoardReorderModel

	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.Handle(w, httpx.BadRequest(err.Error()))
		return
	}

	boards := make([]domain.BoardModel, 0, len(req.Boards))
	for _, b := range req.Boards {
		board, err := h.svc.ReorderBoard(r.Context(), b.ID, b.Position)
		if err != nil {
			httpx.Handle(w, err)
			return
		}
		boards = append(boards, board)
	}

	httpx.OK(w, boards)
}

// DeleteBoard godoc
//
//	@Summary		Delete a board
//	@Description	Deletes a board
//	@Tags			board
//	@Produce		json
//	@Param			boardId	path	string	true	"Board ID"
//	@Success		204
//	@Failure		400		{object}	httpx.ErrBlock
//	@Failure		401		{object}	httpx.ErrBlock
//	@Failure		404		{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/boards/{boardId} [delete]
func (h *Handler) DeleteBoard(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathUUID(r, "boardId")
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	if err := h.svc.DeleteBoard(r.Context(), id); err != nil {
		httpx.Handle(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
