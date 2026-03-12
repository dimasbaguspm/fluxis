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

type OrganisationPagedModel struct {
	Items      []OrganisationModel `json:"items"`
	TotalCount int                 `json:"totalCount"`
	TotalPages int                 `json:"totalPages"`
	Page       int                 `json:"page"`
	PageSize   int                 `json:"pageSize"`
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
	UserId string `json:"userId" validate:"required,uuid4"`
	Role   string `json:"role" validate:"required,oneof=admin member viewer"`
}

type OrganisationMemberUpdateModel struct {
	Role string `json:"role" validate:"required,oneof=admin member viewer"`
}

type Organisations struct {
	ID         []pgtype.UUID `json:"id" validate:"dive,uuid4"`
	Name       []string      `json:"name" validate:"dive,min=1"`
	PageNumber int           `json:"pageNumber" validate:"min=1"`
	PageSize   int           `json:"pageSize" validate:"min=1,max=100"`
	SortBy     string        `json:"sortBy" validate:"oneof=name createdAt updatedAt"`
	SortOrder  string        `json:"sortOrder" validate:"oneof=asc desc"`
}

func (o *Organisations) ApplyDefaults() {
	const (
		defaultPageNumber = 1
		defaultPageSize   = 25
		defaultSortBy     = "updatedAt"
		defaultSortOrder  = "desc"
	)

	if o.PageNumber == 0 {
		o.PageNumber = defaultPageNumber
	}
	if o.PageSize == 0 {
		o.PageSize = defaultPageSize
	}
	if o.SortBy == "" {
		o.SortBy = defaultSortBy
	}
	if o.SortOrder == "" {
		o.SortOrder = defaultSortOrder
	}
}

func (o *OrganisationPagedModel) Empty(pageNumber, pageSize int) OrganisationPagedModel {
	return OrganisationPagedModel{
		Items:      []OrganisationModel{},
		TotalCount: 0,
		TotalPages: 0,
		Page:       pageNumber,
		PageSize:   pageSize,
	}
}

type OrgReader interface {
	ListOrgs(ctx context.Context, q OrganisationSearchModel) ([]OrganisationModel, error)
	GetOrgById(ctx context.Context, id pgtype.UUID) (OrganisationModel, error)
	GetOrgBySlug(ctx context.Context, slug string) (OrganisationModel, error)
	ListMembers(ctx context.Context, orgId pgtype.UUID) ([]OrganisationMemberModel, error)
}

type OrganisationWrite interface {
	CreateOrg(ctx context.Context, p OrganisationCreateModel) (OrganisationModel, error)
	UpdateOrg(ctx context.Context, id pgtype.UUID, p OrganisationUpdateModel) (OrganisationModel, error)
	DeleteOrg(ctx context.Context, id pgtype.UUID) error
	AddMember(ctx context.Context, orgId pgtype.UUID, p OrganisationMemberCreateModel) error
	UpdateMemberRole(ctx context.Context, orgId, userId pgtype.UUID, p OrganisationMemberUpdateModel) error
	RemoveMember(ctx context.Context, orgId, userId pgtype.UUID) error
}
