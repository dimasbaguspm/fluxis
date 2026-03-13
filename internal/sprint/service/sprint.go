package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/dimasbaguspm/fluxis/internal/sprint/repository"
	"github.com/dimasbaguspm/fluxis/pkg/domain"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
	"github.com/dimasbaguspm/fluxis/pkg/pubsub"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrSprintNotFound = httpx.NotFound("sprint not found")
)

func toSprintModel(sprint repository.Sprint) domain.SprintModel {
	var startedAt *time.Time
	if sprint.StartedAt.Valid {
		startedAt = &sprint.StartedAt.Time
	}

	var completedAt *time.Time
	if sprint.CompletedAt.Valid {
		completedAt = &sprint.CompletedAt.Time
	}

	var plannedStartedAt *time.Time
	if sprint.PlannedStartedAt.Valid {
		plannedStartedAt = &sprint.PlannedStartedAt.Time
	}

	var plannedCompletedAt *time.Time
	if sprint.PlannedCompletedAt.Valid {
		plannedCompletedAt = &sprint.PlannedCompletedAt.Time
	}

	goalStr := ""
	if sprint.Goal.Valid {
		goalStr = sprint.Goal.String
	}

	return domain.SprintModel{
		ID:                 sprint.ID,
		ProjectID:          sprint.ProjectID,
		Name:               sprint.Name,
		Goal:               goalStr,
		Status:             string(sprint.Status),
		PlannedStartedAt:   plannedStartedAt,
		PlannedCompletedAt: plannedCompletedAt,
		StartedAt:          startedAt,
		CompletedAt:        completedAt,
		CreatedAt:          sprint.CreatedAt.Time,
		UpdatedAt:          sprint.UpdatedAt.Time,
	}
}

// CreateSprint creates a new sprint
func (s *Service) CreateSprint(ctx context.Context, req domain.SprintCreateModel) (domain.SprintModel, error) {
	project, err := s.Project.GetProjectById(ctx, req.ProjectID)
	if err != nil {
		return domain.SprintModel{}, fmt.Errorf("get project: %w", err)
	}

	sprintStatus := repository.SprintStatusPlanned
	if req.Status != "" {
		sprintStatus = repository.SprintStatus(req.Status)
	}

	goalText := pgtype.Text{Valid: false}
	if req.Goal != "" {
		goalText = pgtype.Text{String: req.Goal, Valid: true}
	}

	plannedStart := pgtype.Timestamptz{Valid: false}
	if req.PlannedStartedAt != "" {
		plannedStart = pgtype.Timestamptz{}
		plannedStart.Scan(req.PlannedStartedAt)
	}

	plannedEnd := pgtype.Timestamptz{Valid: false}
	if req.PlannedCompletedAt != "" {
		plannedEnd = pgtype.Timestamptz{}
		plannedEnd.Scan(req.PlannedCompletedAt)
	}

	sprint, err := s.Repo.CreateSprint(ctx, repository.CreateSprintParams{
		ProjectID:          project.ID,
		Name:               req.Name,
		Goal:               goalText,
		Status:             sprintStatus,
		PlannedStartedAt:   plannedStart,
		PlannedCompletedAt: plannedEnd,
	})
	if err != nil {
		return domain.SprintModel{}, fmt.Errorf("create sprint: %w", err)
	}

	result := toSprintModel(sprint)
	if err := s.Bus.Publish(ctx, pubsub.SprintCreated, httpx.EncodePayload(result)); err != nil {
		slog.Warn("[EventBus]: failed to publish event", "type", string(pubsub.SprintCreated), "error", err)
	}

	return toSprintModel(sprint), nil
}

// GetSprint retrieves a single sprint by ID
func (s *Service) GetSprint(ctx context.Context, id pgtype.UUID) (domain.SprintModel, error) {
	sprint, err := s.Repo.GetSprint(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.SprintModel{}, ErrSprintNotFound
		}
		return domain.SprintModel{}, fmt.Errorf("get sprint: %w", err)
	}

	return toSprintModel(sprint), nil
}

// ListSprintsByProjectPaged lists sprints in a project with pagination
func (s *Service) ListSprintsPaged(ctx context.Context, q domain.SprintsSearchModel) (domain.SprintsPagedModel, error) {
	q.ApplyDefaults()

	sprints, err := s.Repo.ListSprintsPaged(ctx, repository.ListSprintsPagedParams{
		Column1: q.ID,
		Column2: q.ProjectID,
		Column3: q.Name,
		Limit:   int32(q.PageSize),
		Offset:  int32((q.PageNumber - 1) * q.PageSize),
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			emptyResult := domain.SprintsPagedModel{}
			return emptyResult.Empty(q.PageNumber, q.PageSize), nil
		}
		return domain.SprintsPagedModel{}, fmt.Errorf("list sprints by project paged: %w", err)
	}

	var totalCount int64
	data := make([]domain.SprintModel, 0, len(sprints))

	for _, row := range sprints {
		totalCount = row.TotalCount
		data = append(data, toSprintModel(repository.Sprint{
			ID:                 row.ID,
			ProjectID:          row.ProjectID,
			Name:               row.Name,
			Goal:               row.Goal,
			Status:             row.Status,
			PlannedStartedAt:   row.PlannedStartedAt,
			PlannedCompletedAt: row.PlannedCompletedAt,
			StartedAt:          row.StartedAt,
			CompletedAt:        row.CompletedAt,
			CreatedAt:          row.CreatedAt,
			UpdatedAt:          row.UpdatedAt,
			DeletedAt:          row.DeletedAt,
		}))
	}

	totalPages := 0
	if totalCount > 0 {
		totalPages = int((totalCount + int64(q.PageSize) - 1) / int64(q.PageSize))
	}

	return domain.SprintsPagedModel{
		Items:      data,
		TotalCount: int(totalCount),
		TotalPages: totalPages,
		PageNumber: q.PageNumber,
		PageSize:   q.PageSize,
	}, nil
}

// UpdateSprint updates sprint details
func (s *Service) UpdateSprint(ctx context.Context, id pgtype.UUID, req domain.SprintUpdateModel) (domain.SprintModel, error) {
	// Get current sprint to preserve existing values
	current, err := s.Repo.GetSprint(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.SprintModel{}, ErrSprintNotFound
		}
		return domain.SprintModel{}, fmt.Errorf("get sprint: %w", err)
	}

	// Use provided values or keep existing ones
	updatedName := current.Name
	if req.Name != "" {
		updatedName = req.Name
	}

	updatedGoal := current.Goal
	if req.Goal != "" {
		updatedGoal = pgtype.Text{String: req.Goal, Valid: true}
	}

	updatedStatus := current.Status
	if req.Status != "" {
		updatedStatus = repository.SprintStatus(req.Status)
	}

	updatedPlannedStart := current.PlannedStartedAt
	if req.PlannedStartedAt != "" {
		ts := pgtype.Timestamptz{}
		ts.Scan(req.PlannedStartedAt)
		updatedPlannedStart = ts
	}

	updatedPlannedComplete := current.PlannedCompletedAt
	if req.PlannedCompletedAt != "" {
		ts := pgtype.Timestamptz{}
		ts.Scan(req.PlannedCompletedAt)
		updatedPlannedComplete = ts
	}

	sprint, err := s.Repo.UpdateSprint(ctx, repository.UpdateSprintParams{
		ID:                 id,
		Name:               updatedName,
		Goal:               updatedGoal,
		Status:             updatedStatus,
		PlannedStartedAt:   updatedPlannedStart,
		PlannedCompletedAt: updatedPlannedComplete,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.SprintModel{}, ErrSprintNotFound
		}
		return domain.SprintModel{}, fmt.Errorf("update sprint: %w", err)
	}

	result := toSprintModel(sprint)
	if err := s.Bus.Publish(ctx, pubsub.SprintUpdated, httpx.EncodePayload(result)); err != nil {
		slog.Warn("[EventBus]: failed to publish event", "type", string(pubsub.SprintUpdated), "error", err)
	}

	return result, nil
}

// StartSprint transitions a sprint to active status
func (s *Service) StartSprint(ctx context.Context, id pgtype.UUID) (domain.SprintModel, error) {
	sprint, err := s.Repo.StartSprint(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.SprintModel{}, ErrSprintNotFound
		}
		return domain.SprintModel{}, fmt.Errorf("start sprint: %w", err)
	}

	result := toSprintModel(sprint)
	if err := s.Bus.Publish(ctx, pubsub.SprintStarted, httpx.EncodePayload(result)); err != nil {
		slog.Warn("[EventBus]: failed to publish event", "type", string(pubsub.SprintStarted), "error", err)
	}

	return result, nil
}

// CompleteSprint transitions a sprint to completed status
func (s *Service) CompleteSprint(ctx context.Context, id pgtype.UUID) (domain.SprintModel, error) {
	sprint, err := s.Repo.CompleteSprint(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.SprintModel{}, ErrSprintNotFound
		}
		return domain.SprintModel{}, fmt.Errorf("complete sprint: %w", err)
	}

	result := toSprintModel(sprint)
	if err := s.Bus.Publish(ctx, pubsub.SprintCompleted, httpx.EncodePayload(result)); err != nil {
		slog.Warn("[EventBus]: failed to publish event", "type", string(pubsub.SprintCompleted), "error", err)
	}

	return result, nil
}
