package handler

import (
	"net/http"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
)

// ListProjects godoc
//
//	@Summary		List projects with pagination
//	@Description	Returns paginated projects in an organisation with optional filtering
//	@Tags			project
//	@Produce		json
//	@Param			orgId	query	string	true	"Organisation ID"
//	@Param			query	query	domain.ProjectsSearchModel	false	"Search parameters: name, pageNumber, pageSize"
//	@Success		200	{object}	domain.ProjectsPagedModel
//	@Failure		400	{object}	httpx.ErrBlock
//	@Failure		401	{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/projects [get]
func (h *Handler) ListProjects(w http.ResponseWriter, r *http.Request) {
	orgID, err := httpx.QueryUUID(r, "orgId")
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	req := domain.ProjectsSearchModel{
		Name:       httpx.QueryString(r, "name"),
		PageNumber: httpx.QueryNumber(r, "pageNumber"),
		PageSize:   httpx.QueryNumber(r, "pageSize"),
	}

	result, err := h.svc.ListProjectsByOrgPaged(r.Context(), orgID, req)
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	httpx.OK(w, result)
}

// CreateProject godoc
//
//	@Summary		Create a project
//	@Description	Creates a new project in an organisation
//	@Tags			project
//	@Accept			json
//	@Produce		json
//	@Param			orgId	query		string						true	"Organisation ID"
//	@Param			body	body		domain.ProjectCreateModel	true	"Project payload"
//	@Success		201		{object}	domain.ProjectModel
//	@Failure		400		{object}	httpx.ErrBlock
//	@Failure		401		{object}	httpx.ErrBlock
//	@Failure		409		{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/projects [post]
func (h *Handler) CreateProject(w http.ResponseWriter, r *http.Request) {
	orgID, err := httpx.QueryUUID(r, "orgId")
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	var req domain.ProjectCreateModel
	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.Handle(w, httpx.BadRequest(err.Error()))
		return
	}

	project, err := h.svc.CreateProject(r.Context(), orgID, req)
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	httpx.Created(w, project)
}

// GetProject godoc
//
//	@Summary		Get a project
//	@Description	Returns a single project by ID
//	@Tags			project
//	@Produce		json
//	@Param			id	path		string					true	"Project ID"
//	@Success		200	{object}	domain.ProjectModel
//	@Failure		400	{object}	httpx.ErrBlock
//	@Failure		401	{object}	httpx.ErrBlock
//	@Failure		404	{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/projects/{id} [get]
func (h *Handler) GetProject(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathUUID(r, "id")
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	project, err := h.svc.GetProjectById(r.Context(), id)
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	httpx.OK(w, project)
}

// UpdateProject godoc
//
//	@Summary		Update a project
//	@Description	Updates a project's name and description
//	@Tags			project
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Project ID"
//	@Param			body	body		domain.ProjectUpdateModel	true	"Project payload"
//	@Success		200	{object}	domain.ProjectModel
//	@Failure		400	{object}	httpx.ErrBlock
//	@Failure		401	{object}	httpx.ErrBlock
//	@Failure		404	{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/projects/{id} [patch]
func (h *Handler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathUUID(r, "id")
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	var req domain.ProjectUpdateModel
	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.Handle(w, httpx.BadRequest(err.Error()))
		return
	}

	project, err := h.svc.UpdateProject(r.Context(), id, req)
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	httpx.OK(w, project)
}

// UpdateProjectVisibility godoc
//
//	@Summary		Update project visibility
//	@Description	Changes a project's visibility (public/private)
//	@Tags			project
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Project ID"
//	@Param			body	body		domain.ProjectVisibilityModel	true	"Visibility payload"
//	@Success		200	{object}	domain.ProjectModel
//	@Failure		400	{object}	httpx.ErrBlock
//	@Failure		401	{object}	httpx.ErrBlock
//	@Failure		404	{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/projects/{id}/visibility [patch]
func (h *Handler) UpdateProjectVisibility(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathUUID(r, "id")
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	var req domain.ProjectVisibilityModel
	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.Handle(w, httpx.BadRequest(err.Error()))
		return
	}

	project, err := h.svc.UpdateProjectVisibility(r.Context(), id, req)
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	httpx.OK(w, project)
}

// DeleteProject godoc
//
//	@Summary		Delete a project
//	@Description	Soft deletes a project
//	@Tags			project
//	@Param			id	path	string	true	"Project ID"
//	@Success		204
//	@Failure		400	{object}	httpx.ErrBlock
//	@Failure		401	{object}	httpx.ErrBlock
//	@Failure		404	{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/projects/{id} [delete]
func (h *Handler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathUUID(r, "id")
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	if err := h.svc.DeleteProject(r.Context(), id); err != nil {
		httpx.Handle(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
