package handler

import (
	"context"
	"net/http"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
)

// CreateSprint godoc
//
//	@Summary		Create a sprint
//	@Description	Creates a new sprint in a project
//	@Tags			sprint
//	@Accept			json
//	@Produce		json
//	@Param			body		body		domain.SprintCreateModel	true	"Sprint payload"
//	@Success		201			{object}	domain.SprintModel
//	@Failure		400			{object}	httpx.ErrBlock
//	@Failure		401			{object}	httpx.ErrBlock
//	@Failure		404			{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/sprints [post]
func (h *Handler) CreateSprint(w http.ResponseWriter, r *http.Request) {
	var req domain.SprintCreateModel

	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.Handle(w, httpx.BadRequest(err.Error()))
		return
	}

	sprint, err := h.svc.CreateSprint(r.Context(), req)
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	httpx.Created(w, sprint)
}

// ListSprints godoc
//
//	@Summary		List sprints with pagination
//	@Description	Returns paginated sprints in a project with optional filtering
//	@Tags			sprint
//	@Produce		json
//	@Param			query		query		domain.SprintsSearchModel	false	"Search parameters: name, pageNumber, pageSize"
//	@Success		200			{object}	domain.SprintsPagedModel
//	@Failure		400			{object}	httpx.ErrBlock
//	@Failure		401			{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/sprints [get]
func (h *Handler) ListSprints(w http.ResponseWriter, r *http.Request) {
	req := domain.SprintsSearchModel{
		ID:         httpx.QueryUUIDs(r, "id"),
		ProjectID:  httpx.QueryUUIDs(r, "projectId"),
		Name:       httpx.QueryString(r, "name"),
		PageNumber: httpx.QueryNumber(r, "pageNumber"),
		PageSize:   httpx.QueryNumber(r, "pageSize"),
	}

	result, err := h.svc.ListSprintsPaged(r.Context(), req)
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	httpx.OK(w, result)
}

// GetSprint godoc
//
//	@Summary		Get a sprint
//	@Description	Returns a single sprint by ID
//	@Tags			sprint
//	@Produce		json
//	@Param			sprintId	path		string	true	"Sprint ID"
//	@Success		200			{object}	domain.SprintModel
//	@Failure		400			{object}	httpx.ErrBlock
//	@Failure		401			{object}	httpx.ErrBlock
//	@Failure		404			{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/sprints/{sprintId} [get]
func (h *Handler) GetSprint(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathUUID(r, "sprintId")
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	sprint, err := h.sprintCache.GetSingleSprint(r.Context(), id, func(ctx context.Context) (domain.SprintModel, error) {
		return h.svc.GetSprint(ctx, id)
	})
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	httpx.OK(w, sprint)
}

// UpdateSprint godoc
//
//	@Summary		Update a sprint
//	@Description	Updates sprint details
//	@Tags			sprint
//	@Accept			json
//	@Produce		json
//	@Param			sprintId	path		string	true	"Sprint ID"
//	@Param			body		body		domain.SprintUpdateModel	true	"Sprint payload"
//	@Success		200			{object}	domain.SprintModel
//	@Failure		400			{object}	httpx.ErrBlock
//	@Failure		401			{object}	httpx.ErrBlock
//	@Failure		404			{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/sprints/{sprintId} [patch]
func (h *Handler) UpdateSprint(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathUUID(r, "sprintId")
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	var req domain.SprintUpdateModel

	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.Handle(w, httpx.BadRequest(err.Error()))
		return
	}

	sprint, err := h.svc.UpdateSprint(r.Context(), id, req)
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	httpx.OK(w, sprint)
}

// StartSprint godoc
//
//	@Summary		Start a sprint
//	@Description	Transitions a sprint to active status
//	@Tags			sprint
//	@Produce		json
//	@Param			sprintId	path		string	true	"Sprint ID"
//	@Success		200			{object}	domain.SprintModel
//	@Failure		400			{object}	httpx.ErrBlock
//	@Failure		401			{object}	httpx.ErrBlock
//	@Failure		404			{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/sprints/{sprintId}/start [post]
func (h *Handler) StartSprint(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathUUID(r, "sprintId")
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	sprint, err := h.svc.StartSprint(r.Context(), id)
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	httpx.OK(w, sprint)
}

// CompleteSprint godoc
//
//	@Summary		Complete a sprint
//	@Description	Transitions a sprint to completed status
//	@Tags			sprint
//	@Produce		json
//	@Param			sprintId	path		string	true	"Sprint ID"
//	@Success		200			{object}	domain.SprintModel
//	@Failure		400			{object}	httpx.ErrBlock
//	@Failure		401			{object}	httpx.ErrBlock
//	@Failure		404			{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/sprints/{sprintId}/completed [post]
func (h *Handler) CompleteSprint(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathUUID(r, "sprintId")
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	sprint, err := h.svc.CompleteSprint(r.Context(), id)
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	httpx.OK(w, sprint)
}
