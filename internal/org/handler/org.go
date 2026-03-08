package handler

import (
	"net/http"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
	"github.com/jackc/pgx/v5/pgtype"
)

// ListOrgs godoc
//
//	@Summary		List organisations
//	@Description	Returns all organisations the authenticated user belongs to
//	@Tags			org
//	@Produce		json
//	@Success		200	{array}		domain.OrganisationModel
//	@Failure		401	{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/orgs [get]
func (h *Handler) ListOrgs(w http.ResponseWriter, r *http.Request) {
	userID := httpx.MustUserID(r.Context())

	req := domain.OrganisationSearchModel{
		UserId: userID,
	}

	orgs, err := h.svc.ListOrgs(r.Context(), req)
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	httpx.OK(w, orgs)
}

// CreateOrg godoc
//
//	@Summary		Create an organisation
//	@Description	Creates a new organisation
//	@Tags			org
//	@Accept			json
//	@Produce		json
//	@Param			body	body		domain.OrganisationCreateModel	true	"Organisation payload"
//	@Success		201		{object}	domain.OrganisationModel
//	@Failure		400		{object}	httpx.ErrBlock
//	@Failure		401		{object}	httpx.ErrBlock
//	@Failure		409		{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/orgs [post]
func (h *Handler) CreateOrg(w http.ResponseWriter, r *http.Request) {
	var req domain.OrganisationCreateModel
	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.Handle(w, httpx.BadRequest(err.Error()))
		return
	}

	org, err := h.svc.CreateOrg(r.Context(), req)
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	httpx.Created(w, org)
}

// GetOrg godoc
//
//	@Summary		Get an organisation
//	@Description	Returns a single organisation by ID
//	@Tags			org
//	@Produce		json
//	@Param			id	path		string	true	"Organisation ID"
//	@Success		200	{object}	domain.OrganisationModel
//	@Failure		400	{object}	httpx.ErrBlock
//	@Failure		401	{object}	httpx.ErrBlock
//	@Failure		404	{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/orgs/{id} [get]
func (h *Handler) GetOrg(w http.ResponseWriter, r *http.Request) {
	var id pgtype.UUID
	if err := id.Scan(r.PathValue("id")); err != nil {
		httpx.Handle(w, httpx.BadRequest("invalid org id"))
		return
	}

	org, err := h.svc.GetOrgById(r.Context(), id)
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	httpx.OK(w, org)
}

// UpdateOrg godoc
//
//	@Summary		Update an organisation
//	@Description	Updates an organisation's name
//	@Tags			org
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string							true	"Organisation ID"
//	@Param			body	body		domain.OrganisationUpdateModel	true	"Update payload"
//	@Success		200		{object}	domain.OrganisationModel
//	@Failure		400		{object}	httpx.ErrBlock
//	@Failure		401		{object}	httpx.ErrBlock
//	@Failure		404		{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/orgs/{id} [patch]
func (h *Handler) UpdateOrg(w http.ResponseWriter, r *http.Request) {
	var id pgtype.UUID
	if err := id.Scan(r.PathValue("id")); err != nil {
		httpx.Handle(w, httpx.BadRequest("invalid org id"))
		return
	}

	var req domain.OrganisationUpdateModel
	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.Handle(w, httpx.BadRequest(err.Error()))
		return
	}

	org, err := h.svc.UpdateOrg(r.Context(), id, req)
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	httpx.OK(w, org)
}

// DeleteOrg godoc
//
//	@Summary		Delete an organisation
//	@Description	Soft-deletes an organisation by ID
//	@Tags			org
//	@Param			id	path	string	true	"Organisation ID"
//	@Success		204
//	@Failure		400	{object}	httpx.ErrBlock
//	@Failure		401	{object}	httpx.ErrBlock
//	@Failure		404	{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/orgs/{id} [delete]
func (h *Handler) DeleteOrg(w http.ResponseWriter, r *http.Request) {
	var id pgtype.UUID
	if err := id.Scan(r.PathValue("id")); err != nil {
		httpx.Handle(w, httpx.BadRequest("invalid org id"))
		return
	}

	if err := h.svc.DeleteOrg(r.Context(), id); err != nil {
		httpx.Handle(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
