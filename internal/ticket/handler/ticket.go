package handler

import (
	"net/http"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
)

// ListTickets godoc
//
//	@Summary		List tickets with pagination
//	@Description	Returns paginated tickets for a project, optionally filtered by sprint or board
//	@Tags			ticket
//	@Produce		json
//	@Param			query	query	domain.TicketSearchModel	false	"Search parameters: projectId (required), sprintId (optional), boardId (optional), pageNumber, pageSize"
//	@Success		200	{object}	domain.TicketsPagedModel
//	@Failure		400	{object}	httpx.ErrBlock
//	@Failure		401	{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/tickets [get]
func (h *Handler) ListTickets(w http.ResponseWriter, r *http.Request) {
	req := domain.TicketSearchModel{
		ID:         httpx.QueryUUIDs(r, "id"),
		ProjectID:  httpx.QueryUUIDs(r, "projectId"),
		SprintID:   httpx.QueryUUIDs(r, "sprintId"),
		BoardID:    httpx.QueryUUIDs(r, "boardId"),
		PageNumber: httpx.QueryNumber(r, "pageNumber"),
		PageSize:   httpx.QueryNumber(r, "pageSize"),
	}

	tickets, err := h.svc.ListTickets(r.Context(), req)
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	httpx.OK(w, tickets)
}

// GetTicket godoc
//
//	@Summary		Get a ticket
//	@Description	Returns a single ticket by ID
//	@Tags			ticket
//	@Produce		json
//	@Param			ticketId	path		string	true	"Ticket ID"
//	@Success		200	{object}	domain.TicketModel
//	@Failure		400	{object}	httpx.ErrBlock
//	@Failure		401	{object}	httpx.ErrBlock
//	@Failure		404	{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/tickets/{ticketId} [get]
func (h *Handler) GetTicket(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathUUID(r, "ticketId")
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	ticket, err := h.svc.GetTicket(r.Context(), id)
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	httpx.OK(w, ticket)
}

// CreateTicket godoc
//
//	@Summary		Create a ticket
//	@Description	Creates a new ticket in a project
//	@Tags			ticket
//	@Accept			json
//	@Produce		json
//	@Param			projectId	query		string						true	"Project ID"
//	@Param			body		body		domain.TicketCreateModel	true	"Ticket payload"
//	@Success		201			{object}	domain.TicketModel
//	@Failure		400			{object}	httpx.ErrBlock
//	@Failure		401			{object}	httpx.ErrBlock
//	@Failure		404			{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/tickets [post]
func (h *Handler) CreateTicket(w http.ResponseWriter, r *http.Request) {
	projectID, err := httpx.QueryUUID(r, "projectId")
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	var req domain.TicketCreateModel
	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.Handle(w, httpx.BadRequest(err.Error()))
		return
	}

	ticket, err := h.svc.CreateTicket(r.Context(), projectID, req)
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	httpx.Created(w, ticket)
}

// UpdateTicket godoc
//
//	@Summary		Update a ticket
//	@Description	Updates ticket details (title, description, priority, type, etc.)
//	@Tags			ticket
//	@Accept			json
//	@Produce		json
//	@Param			ticketId	path		string					true	"Ticket ID"
//	@Param			body		body		domain.TicketUpdateModel	true	"Update payload"
//	@Success		200			{object}	domain.TicketModel
//	@Failure		400			{object}	httpx.ErrBlock
//	@Failure		401			{object}	httpx.ErrBlock
//	@Failure		404			{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/tickets/{ticketId} [patch]
func (h *Handler) UpdateTicket(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathUUID(r, "ticketId")
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	var req domain.TicketUpdateModel
	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.Handle(w, httpx.BadRequest(err.Error()))
		return
	}

	ticket, err := h.svc.UpdateTicket(r.Context(), id, req)
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	httpx.OK(w, ticket)
}

// MoveTicketToBoard godoc
//
//	@Summary		Move ticket to board column
//	@Description	Moves a ticket to a specific board and column
//	@Tags			ticket
//	@Accept			json
//	@Produce		json
//	@Param			ticketId	path		string							true	"Ticket ID"
//	@Param			body		body		domain.TicketBoardMoveModel	true	"Move payload"
//	@Success		200			{object}	domain.TicketModel
//	@Failure		400			{object}	httpx.ErrBlock
//	@Failure		401			{object}	httpx.ErrBlock
//	@Failure		404			{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/tickets/{ticketId}/move-to-board [patch]
func (h *Handler) MoveTicketToBoard(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathUUID(r, "ticketId")
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	var req domain.TicketBoardMoveModel
	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.Handle(w, httpx.BadRequest(err.Error()))
		return
	}

	ticket, err := h.svc.MoveTicketToBoard(r.Context(), id, req)
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	httpx.OK(w, ticket)
}

// MoveTicketToSprint godoc
//
//	@Summary		Move ticket to sprint
//	@Description	Moves a ticket to a specific sprint
//	@Tags			ticket
//	@Accept			json
//	@Produce		json
//	@Param			ticketId	path		string	true	"Ticket ID"
//	@Param			sprintId	query		string	true	"Sprint ID"
//	@Success		200			{object}	domain.TicketModel
//	@Failure		400			{object}	httpx.ErrBlock
//	@Failure		401			{object}	httpx.ErrBlock
//	@Failure		404			{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/tickets/{ticketId}/move-to-sprint [patch]
func (h *Handler) MoveTicketToSprint(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathUUID(r, "ticketId")
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	sprintID, err := httpx.QueryUUID(r, "sprintId")
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	ticket, err := h.svc.MoveTicketToSprint(r.Context(), id, sprintID)
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	httpx.OK(w, ticket)
}

// MoveTicketToBoardColumn godoc
//
//	@Summary		Move ticket to board column
//	@Description	Moves a ticket to a specific board column
//	@Tags			ticket
//	@Accept			json
//	@Produce		json
//	@Param			ticketId	path		string							true	"Ticket ID"
//	@Param			body		body		domain.TicketBoardMoveModel	true	"Move payload"
//	@Success		200			{object}	domain.TicketModel
//	@Failure		400			{object}	httpx.ErrBlock
//	@Failure		401			{object}	httpx.ErrBlock
//	@Failure		404			{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/tickets/{ticketId}/move-board-column [patch]
func (h *Handler) MoveTicketToBoardColumn(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathUUID(r, "ticketId")
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	var req domain.TicketBoardMoveModel
	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.Handle(w, httpx.BadRequest(err.Error()))
		return
	}

	ticket, err := h.svc.MoveTicketToBoardColumn(r.Context(), id, req)
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	httpx.OK(w, ticket)
}

// DeleteTicket godoc
//
//	@Summary		Delete a ticket
//	@Description	Soft-deletes a ticket by ID
//	@Tags			ticket
//	@Param			ticketId	path	string	true	"Ticket ID"
//	@Success		204
//	@Failure		400	{object}	httpx.ErrBlock
//	@Failure		401	{object}	httpx.ErrBlock
//	@Failure		404	{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/tickets/{ticketId} [delete]
func (h *Handler) DeleteTicket(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathUUID(r, "ticketId")
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	if err := h.svc.DeleteTicket(r.Context(), id); err != nil {
		httpx.Handle(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
