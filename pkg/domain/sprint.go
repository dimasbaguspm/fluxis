package domain

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type SprintModel struct {
	ID                 pgtype.UUID `json:"id" validate:"required,uuid4"`
	ProjectID          pgtype.UUID `json:"projectId" validate:"required,uuid4"`
	Name               string      `json:"name" validate:"required,min=1"`
	Goal               string      `json:"goal"`
	Status             string      `json:"status" validate:"required,oneof=planned active completed"`
	PlannedStartedAt   *time.Time  `json:"plannedStartedAt"`
	PlannedCompletedAt *time.Time  `json:"plannedCompletedAt"`
	StartedAt          *time.Time  `json:"startedAt"`
	CompletedAt        *time.Time  `json:"completedAt"`
	CreatedAt          time.Time   `json:"createdAt"`
	UpdatedAt          time.Time   `json:"updatedAt"`
}

type SprintCreateModel struct {
	Name               string      `json:"name" validate:"required,min=1"`
	ProjectID          pgtype.UUID `json:"projectId" validate:"required"`
	Goal               string      `json:"goal,omitempty"`
	PlannedStartedAt   string      `json:"plannedStartedAt,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05.000Z07:00"`
	PlannedCompletedAt string      `json:"plannedCompletedAt,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05.000Z07:00"`
}

type SprintUpdateModel struct {
	Name               string `json:"name,omitempty" validate:"omitempty,min=1"`
	Goal               string `json:"goal,omitempty"`
	PlannedStartedAt   string `json:"plannedStartedAt,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05.000Z07:00"`
	PlannedCompletedAt string `json:"plannedCompletedAt,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05.000Z07:00"`
}

type SprintsSearchModel struct {
	ID         []pgtype.UUID `json:"id" validate:"omitempty,dive,uuid4"`
	ProjectID  []pgtype.UUID `json:"projectId" validate:"omitempty,dive,uuid4"`
	Name       string        `json:"name"`
	PageNumber int           `json:"pageNumber" validate:"omitempty,min=1"`
	PageSize   int           `json:"pageSize" validate:"omitempty,min=1,max=100"`
}

func (s *SprintsSearchModel) ApplyDefaults() {
	const (
		defaultPageNumber = 1
		defaultPageSize   = 25
	)
	if s.PageNumber == 0 {
		s.PageNumber = defaultPageNumber
	}
	if s.PageSize == 0 {
		s.PageSize = defaultPageSize
	}
}

type SprintsPagedModel struct {
	Items      []SprintModel `json:"items"`
	TotalCount int           `json:"totalCount"`
	TotalPages int           `json:"totalPages"`
	PageNumber int           `json:"pageNumber"`
	PageSize   int           `json:"pageSize"`
}

func (m SprintsPagedModel) Empty(pageNumber, pageSize int) SprintsPagedModel {
	return SprintsPagedModel{
		Items:      []SprintModel{},
		TotalCount: 0,
		TotalPages: 0,
		PageNumber: pageNumber,
		PageSize:   pageSize,
	}
}

type SprintReader interface {
	GetSprint(ctx context.Context, id pgtype.UUID) (SprintModel, error)
	ListSprintsPaged(ctx context.Context, q SprintsSearchModel) (SprintsPagedModel, error)
}

type SprintWriter interface {
	CreateSprint(ctx context.Context, p SprintCreateModel) (SprintModel, error)
	UpdateSprint(ctx context.Context, id pgtype.UUID, p SprintUpdateModel) (SprintModel, error)
	StartSprint(ctx context.Context, id pgtype.UUID) (SprintModel, error)
	CompleteSprint(ctx context.Context, id pgtype.UUID) (SprintModel, error)
}
