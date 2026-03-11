package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dimasbaguspm/fluxis/internal/sprint/repository"
	"github.com/dimasbaguspm/fluxis/pkg/domain"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
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
	if req.Status != nil {
		sprintStatus = repository.SprintStatus(*req.Status)
	}

	goalText := pgtype.Text{Valid: false}
	if req.Goal != nil {
		goalText = pgtype.Text{String: *req.Goal, Valid: true}
	}

	plannedStart := pgtype.Timestamptz{Valid: false}
	if req.PlannedStartedAt != nil {
		plannedStart = pgtype.Timestamptz{}
		plannedStart.Scan(*req.PlannedStartedAt)
	}

	plannedEnd := pgtype.Timestamptz{Valid: false}
	if req.PlannedCompletedAt != nil {
		plannedEnd = pgtype.Timestamptz{}
		plannedEnd.Scan(*req.PlannedCompletedAt)
	}

	sprint, err := s.Repo.CreateSprint(ctx, repository.CreateSprintParams{
		ProjectID:          project.ID,
		Name:               *req.Name,
		Goal:               goalText,
		Status:             sprintStatus,
		PlannedStartedAt:   plannedStart,
		PlannedCompletedAt: plannedEnd,
	})
	if err != nil {
		return domain.SprintModel{}, fmt.Errorf("create sprint: %w", err)
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

// ListSprintsByProject lists all sprints in a project
func (s *Service) ListSprintsByProject(ctx context.Context, projectID pgtype.UUID) ([]domain.SprintModel, error) {
	sprints, err := s.Repo.ListSprintsByProject(ctx, projectID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []domain.SprintModel{}, nil
		}
		return []domain.SprintModel{}, fmt.Errorf("list sprints: %w", err)
	}

	data := make([]domain.SprintModel, len(sprints))
	for i, sprint := range sprints {
		data[i] = toSprintModel(sprint)
	}

	return data, nil
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
	if req.Name != nil {
		updatedName = *req.Name
	}

	updatedGoal := current.Goal
	if req.Goal != nil {
		updatedGoal = pgtype.Text{String: *req.Goal, Valid: true}
	}

	updatedStatus := current.Status
	if req.Status != nil {
		updatedStatus = repository.SprintStatus(*req.Status)
	}

	sprint, err := s.Repo.UpdateSprint(ctx, repository.UpdateSprintParams{
		ID:                 id,
		Name:               updatedName,
		Goal:               updatedGoal,
		Status:             updatedStatus,
		PlannedStartedAt:   current.PlannedStartedAt,
		PlannedCompletedAt: current.PlannedCompletedAt,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.SprintModel{}, ErrSprintNotFound
		}
		return domain.SprintModel{}, fmt.Errorf("update sprint: %w", err)
	}

	return toSprintModel(sprint), nil
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

	return toSprintModel(sprint), nil
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

	return toSprintModel(sprint), nil
}
