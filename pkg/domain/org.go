package domain

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type OrganisationSearchModel struct {
	OrgId  pgtype.UUID `json:"orgId" validate:"uuid4"`
	UserId pgtype.UUID `json:"userId" validate:"uuid4"`
}

type OrganisationModel struct {
	ID           pgtype.UUID `json:"id" validate:"required,uuid4"`
	Name         string      `json:"name" validate:"min=1"`
	Slug         string      `json:"slug"`
	TotalMembers int64       `json:"totalMembers"`
	CreatedAt    time.Time   `json:"createdAt"`
	UpdatedAt    time.Time   `json:"updatedAt"`
}

type OrganisationCreateModel struct {
	Name string `json:"name" validate:"required,min=1"`
}

type OrganisationUpdateModel struct {
	Name string `json:"name" validate:"min=1"`
}

type OrganisationMemberModel struct {
	UserID   pgtype.UUID `json:"userId"`
	Name     string      `json:"name"`
	Email    string      `json:"email"`
	Role     string      `json:"role"`
	JoinedAt time.Time   `json:"joinedAt"`
}

type OrganisationMemberCreateModel struct {
	UserID pgtype.UUID `json:"userId" validate:"required,uuid4"`
	Role   string      `json:"role" validate:"required,oneof=admin member viewer"`
}

type OrganisationMemberUpdateModel struct {
	UserID pgtype.UUID `json:"userId" validate:"required,uuid4"`
	Role   string      `json:"role" validate:"required,oneof=admin member viewer"`
}

type OrganisationMemberRemoveModel struct {
	UserID pgtype.UUID `json:"userId" validate:"required,uuid4"`
}

type OrganisationRead interface {
	GetListOrganisations(ctx context.Context, q OrganisationSearchModel) ([]OrganisationModel, error)
	GetSingleOrganisationById(ctx context.Context, id pgtype.UUID) (OrganisationModel, error)
	GetSingleOrganisationBySlug(ctx context.Context, slug string) (OrganisationModel, error)
	GetListOrganisationMembers(ctx context.Context, orgId pgtype.UUID) ([]OrganisationMemberModel, error)
}

type OrganisationWrite interface {
	CreateOrganisation(ctx context.Context, p OrganisationCreateModel) (OrganisationModel, error)
	UpdateOrganisation(ctx context.Context, id pgtype.UUID, p OrganisationUpdateModel) (OrganisationModel, error)
	DeleteOrganisation(ctx context.Context, id pgtype.UUID) error
	AddOrganisationMember(ctx context.Context, orgId pgtype.UUID, p OrganisationMemberCreateModel) error
	UpdateOrganisationMemberRole(ctx context.Context, orgId pgtype.UUID, p OrganisationMemberUpdateModel) error
	RemoveOrganisationMember(ctx context.Context, orgId pgtype.UUID, p OrganisationMemberRemoveModel) error
}
