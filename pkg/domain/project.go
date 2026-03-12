package domain

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type ProjectModel struct {
	ID          pgtype.UUID `json:"id" validate:"required,uuid4"`
	OrgID       pgtype.UUID `json:"orgId" validate:"required,uuid4"`
	Key         string      `json:"key" validate:"required,min=1"`
	Name        string      `json:"name" validate:"required,min=1"`
	Description string      `json:"description"`
	Visibility  string      `json:"visibility" validate:"required,oneof=public private"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
}

type ProjectCreateModel struct {
	Key         string `json:"key" validate:"required,min=1,max=10"`
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description"`
	Visibility  string `json:"visibility" validate:"required,oneof=public private"`
}

type ProjectUpdateModel struct {
	Name        string `json:"name" validate:"min=1,max=100"`
	Description string `json:"description"`
}

type ProjectVisibilityModel struct {
	Visibility string `json:"visibility" validate:"required,oneof=public private"`
}

type ProjectsSearchModel struct {
	Name       string `json:"name"`
	PageNumber int    `json:"pageNumber" validate:"min=1"`
	PageSize   int    `json:"pageSize" validate:"min=1,max=100"`
}

type ProjectsPagedModel struct {
	Items      []ProjectModel `json:"items"`
	TotalCount int            `json:"totalCount"`
	TotalPages int            `json:"totalPages"`
	PageNumber int            `json:"pageNumber"`
	PageSize   int            `json:"pageSize"`
}

func (m *ProjectsSearchModel) ApplyDefaults() {
	const (
		defaultPageNumber = 1
		defaultPageSize   = 25
	)

	if m.PageNumber == 0 {
		m.PageNumber = defaultPageNumber
	}
	if m.PageSize == 0 {
		m.PageSize = defaultPageSize
	}
}

func (m *ProjectsPagedModel) Empty(pageNumber, pageSize int) ProjectsPagedModel {
	return ProjectsPagedModel{
		Items:      []ProjectModel{},
		TotalCount: 0,
		TotalPages: 0,
		PageNumber: pageNumber,
		PageSize:   pageSize,
	}
}

type ProjectReader interface {
	GetProjectById(ctx context.Context, id pgtype.UUID) (ProjectModel, error)
	GetProjectByKey(ctx context.Context, orgId pgtype.UUID, key string) (ProjectModel, error)
	ListProjectsByOrg(ctx context.Context, orgId pgtype.UUID) ([]ProjectModel, error)
	ListProjectsByOrgPaged(ctx context.Context, orgId pgtype.UUID, q ProjectsSearchModel) (ProjectsPagedModel, error)
}

type ProjectWriter interface {
	CreateProject(ctx context.Context, orgId pgtype.UUID, p ProjectCreateModel) (ProjectModel, error)
	UpdateProject(ctx context.Context, id pgtype.UUID, p ProjectUpdateModel) (ProjectModel, error)
	UpdateProjectVisibility(ctx context.Context, id pgtype.UUID, p ProjectVisibilityModel) (ProjectModel, error)
	DeleteProject(ctx context.Context, id pgtype.UUID) error
}
