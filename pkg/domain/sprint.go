package domain

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type SprintModel struct {
	ID                   pgtype.UUID `json:"id"`
	ProjectID            pgtype.UUID `json:"projectId"`
	Name                 string      `json:"name" validate:"required,min=1"`
	Goal                 string      `json:"goal"`
	Status               string      `json:"status" validate:"required,oneof=planned active completed"`
	PlannedStartedAt     *time.Time  `json:"plannedStartedAt"`
	PlannedCompletedAt   *time.Time  `json:"plannedCompletedAt"`
	StartedAt            *time.Time  `json:"startedAt"`
	CompletedAt          *time.Time  `json:"completedAt"`
	CreatedAt            time.Time   `json:"createdAt"`
	UpdatedAt            time.Time   `json:"updatedAt"`
}

type SprintCreateModel struct {
	Name                *string `json:"name" validate:"required,min=1"`
	Goal                *string `json:"goal"`
	Status              *string `json:"status" validate:"omitempty,oneof=planned active completed"`
	PlannedStartedAt    *string `json:"plannedStartedAt"`
	PlannedCompletedAt  *string `json:"plannedCompletedAt"`
}

type SprintUpdateModel struct {
	Name   *string `json:"name" validate:"omitempty,min=1"`
	Goal   *string `json:"goal"`
	Status *string `json:"status" validate:"omitempty,oneof=planned active completed"`
}

type SprintReader interface {
	GetSprint(ctx context.Context, id pgtype.UUID) (SprintModel, error)
	ListSprintsByProject(ctx context.Context, projectID pgtype.UUID) ([]SprintModel, error)
}

type SprintWriter interface {
	CreateSprint(ctx context.Context, projectID pgtype.UUID, p SprintCreateModel) (SprintModel, error)
	UpdateSprint(ctx context.Context, id pgtype.UUID, p SprintUpdateModel) (SprintModel, error)
	StartSprint(ctx context.Context, id pgtype.UUID) (SprintModel, error)
	CompleteSprint(ctx context.Context, id pgtype.UUID) (SprintModel, error)
}
