package handler

import (
	"net/http"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
	"github.com/jackc/pgx/v5/pgtype"
)

// ListOrgMembers godoc
//
//	@Summary		List organisation members
//	@Description	Returns all members of an organisation
//	@Tags			org
//	@Produce		json
//	@Param			id	path		string	true	"Organisation ID"
//	@Success		200	{array}		domain.OrganisationMemberModel
//	@Failure		400	{object}	httpx.ErrBlock
//	@Failure		401	{object}	httpx.ErrBlock
//	@Failure		404	{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/orgs/{id}/members [get]
func (h *Handler) ListOrgMembers(w http.ResponseWriter, r *http.Request) {
	var orgID pgtype.UUID
	if err := orgID.Scan(r.PathValue("id")); err != nil {
		httpx.Handle(w, httpx.BadRequest("invalid org id"))
		return
	}

	members, err := h.svc.GetListOrganisationMembers(r.Context(), orgID)
	if err != nil {
		httpx.Handle(w, err)
		return
	}

	httpx.OK(w, members)
}

// AddOrgMember godoc
//
//	@Summary		Add a member to an organisation
//	@Description	Adds a user to an organisation with a given role
//	@Tags			org
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string								true	"Organisation ID"
//	@Param			body	body		domain.OrganisationMemberCreateModel	true	"Member payload"
//	@Success		201
//	@Failure		400	{object}	httpx.ErrBlock
//	@Failure		401	{object}	httpx.ErrBlock
//	@Failure		404	{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/orgs/{id}/members [post]
func (h *Handler) AddOrgMember(w http.ResponseWriter, r *http.Request) {
	var orgID pgtype.UUID
	if err := orgID.Scan(r.PathValue("id")); err != nil {
		httpx.Handle(w, httpx.BadRequest("invalid org id"))
		return
	}

	var req domain.OrganisationMemberCreateModel
	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.Handle(w, httpx.BadRequest(err.Error()))
		return
	}

	if err := h.svc.AddOrganisationMember(r.Context(), orgID, req); err != nil {
		httpx.Handle(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// UpdateOrgMember godoc
//
//	@Summary		Update a member's role
//	@Description	Updates the role of a member within an organisation
//	@Tags			org
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string								true	"Organisation ID"
//	@Param			body	body		domain.OrganisationMemberUpdateModel	true	"Update payload"
//	@Success		200
//	@Failure		400	{object}	httpx.ErrBlock
//	@Failure		401	{object}	httpx.ErrBlock
//	@Failure		404	{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/orgs/{id}/members [patch]
func (h *Handler) UpdateOrgMember(w http.ResponseWriter, r *http.Request) {
	var orgID pgtype.UUID
	if err := orgID.Scan(r.PathValue("id")); err != nil {
		httpx.Handle(w, httpx.BadRequest("invalid org id"))
		return
	}

	var req domain.OrganisationMemberUpdateModel
	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.Handle(w, httpx.BadRequest(err.Error()))
		return
	}

	if err := h.svc.UpdateOrganisationMemberRole(r.Context(), orgID, req); err != nil {
		httpx.Handle(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UpdateOrgMember godoc
//
//	@Summary		Delete a member from an organsiation
//	@Description	Delete a user from an organisation
//	@Tags			org
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string								true	"Organisation ID"
//	@Param			body	body		domain.OrganisationMemberRemoveModel	true	"Update payload"
//	@Success		200
//	@Failure		400	{object}	httpx.ErrBlock
//	@Failure		401	{object}	httpx.ErrBlock
//	@Failure		404	{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/orgs/{id}/members [delete]
func (h *Handler) DeleteOrgMember(w http.ResponseWriter, r *http.Request) {
	var orgID pgtype.UUID
	if err := orgID.Scan(r.PathValue("id")); err != nil {
		httpx.Handle(w, httpx.BadRequest("invalid org id"))
		return
	}

	var req domain.OrganisationMemberUpdateModel
	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.Handle(w, httpx.BadRequest(err.Error()))
		return
	}

	if err := h.svc.UpdateOrganisationMemberRole(r.Context(), orgID, req); err != nil {
		httpx.Handle(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
