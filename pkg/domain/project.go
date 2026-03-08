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

type ProjectReader interface {
	GetProjectById(ctx context.Context, id pgtype.UUID) (ProjectModel, error)
	GetProjectByKey(ctx context.Context, orgId pgtype.UUID, key string) (ProjectModel, error)
	ListProjectsByOrg(ctx context.Context, orgId pgtype.UUID) ([]ProjectModel, error)
}

type ProjectWriter interface {
	CreateProject(ctx context.Context, orgId pgtype.UUID, p ProjectCreateModel) (ProjectModel, error)
	UpdateProject(ctx context.Context, id pgtype.UUID, p ProjectUpdateModel) (ProjectModel, error)
	UpdateProjectVisibility(ctx context.Context, id pgtype.UUID, p ProjectVisibilityModel) (ProjectModel, error)
	DeleteProject(ctx context.Context, id pgtype.UUID) error
}
