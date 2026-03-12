package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/dimasbaguspm/fluxis/internal/project/repository"
	"github.com/dimasbaguspm/fluxis/pkg/domain"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrProjectNotFound = httpx.NotFound("project not found")
	ErrKeyIsTaken      = httpx.Conflict("project key has been taken")
)

func (s *Service) GetProjectById(ctx context.Context, id pgtype.UUID) (domain.ProjectModel, error) {
	project, err := s.Repo.GetProject(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ProjectModel{}, ErrProjectNotFound
		}
		return domain.ProjectModel{}, fmt.Errorf("get project by id: %w", err)
	}

	return domain.ProjectModel{
		ID:          project.ID,
		OrgID:       project.OrgID,
		Key:         project.Key,
		Name:        project.Name,
		Description: project.Description.String,
		Visibility:  string(project.Visibility),
		CreatedAt:   project.CreatedAt.Time,
		UpdatedAt:   project.UpdatedAt.Time,
	}, nil
}

func (s *Service) GetProjectByKey(ctx context.Context, orgId pgtype.UUID, key string) (domain.ProjectModel, error) {
	project, err := s.Repo.GetProjectByKey(ctx, repository.GetProjectByKeyParams{
		OrgID: orgId,
		Key:   key,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ProjectModel{}, ErrProjectNotFound
		}
		return domain.ProjectModel{}, fmt.Errorf("get project by key: %w", err)
	}

	return domain.ProjectModel{
		ID:          project.ID,
		OrgID:       project.OrgID,
		Key:         project.Key,
		Name:        project.Name,
		Description: project.Description.String,
		Visibility:  string(project.Visibility),
		CreatedAt:   project.CreatedAt.Time,
		UpdatedAt:   project.UpdatedAt.Time,
	}, nil
}

func (s *Service) ListProjectsByOrg(ctx context.Context, orgId pgtype.UUID) ([]domain.ProjectModel, error) {
	projects, err := s.Repo.ListProjectsByOrg(ctx, orgId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []domain.ProjectModel{}, nil
		}
		return []domain.ProjectModel{}, fmt.Errorf("list projects by org: %w", err)
	}

	data := make([]domain.ProjectModel, len(projects))
	for i, project := range projects {
		data[i] = domain.ProjectModel{
			ID:          project.ID,
			OrgID:       project.OrgID,
			Key:         project.Key,
			Name:        project.Name,
			Description: project.Description.String,
			Visibility:  string(project.Visibility),
			CreatedAt:   project.CreatedAt.Time,
			UpdatedAt:   project.UpdatedAt.Time,
		}
	}

	return data, nil
}

func (s *Service) CreateProject(ctx context.Context, orgId pgtype.UUID, p domain.ProjectCreateModel) (domain.ProjectModel, error) {
	org, err := s.Org.GetOrgById(ctx, orgId)
	if err != nil {
		return domain.ProjectModel{}, err
	}

	project, err := s.Repo.CreateProject(ctx, repository.CreateProjectParams{
		OrgID:       org.ID,
		Key:         p.Key,
		Name:        p.Name,
		Description: pgtype.Text{String: p.Description, Valid: p.Description != ""},
		Visibility:  repository.ProjectVisibility(p.Visibility),
	})
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" { // unique constraint violation
				return domain.ProjectModel{}, ErrKeyIsTaken
			}
		}
		return domain.ProjectModel{}, fmt.Errorf("create project: %w", err)
	}

	return domain.ProjectModel{
		ID:          project.ID,
		OrgID:       project.OrgID,
		Key:         project.Key,
		Name:        project.Name,
		Description: project.Description.String,
		Visibility:  string(project.Visibility),
		CreatedAt:   project.CreatedAt.Time,
		UpdatedAt:   project.UpdatedAt.Time,
	}, nil
}

func (s *Service) UpdateProject(ctx context.Context, id pgtype.UUID, p domain.ProjectUpdateModel) (domain.ProjectModel, error) {
	project, err := s.Repo.UpdateProject(ctx, repository.UpdateProjectParams{
		ID:          id,
		Name:        p.Name,
		Description: pgtype.Text{String: p.Description, Valid: p.Description != ""},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ProjectModel{}, ErrProjectNotFound
		}
		return domain.ProjectModel{}, fmt.Errorf("update project: %w", err)
	}

	return domain.ProjectModel{
		ID:          project.ID,
		OrgID:       project.OrgID,
		Key:         project.Key,
		Name:        project.Name,
		Description: project.Description.String,
		Visibility:  string(project.Visibility),
		CreatedAt:   project.CreatedAt.Time,
		UpdatedAt:   project.UpdatedAt.Time,
	}, nil
}

func (s *Service) UpdateProjectVisibility(ctx context.Context, id pgtype.UUID, p domain.ProjectVisibilityModel) (domain.ProjectModel, error) {
	project, err := s.Repo.UpdateProjectVisibility(ctx, repository.UpdateProjectVisibilityParams{
		ID:         id,
		Visibility: repository.ProjectVisibility(p.Visibility),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ProjectModel{}, ErrProjectNotFound
		}
		return domain.ProjectModel{}, fmt.Errorf("update project visibility: %w", err)
	}

	return domain.ProjectModel{
		ID:          project.ID,
		OrgID:       project.OrgID,
		Key:         project.Key,
		Name:        project.Name,
		Description: project.Description.String,
		Visibility:  string(project.Visibility),
		CreatedAt:   project.CreatedAt.Time,
		UpdatedAt:   project.UpdatedAt.Time,
	}, nil
}

func (s *Service) DeleteProject(ctx context.Context, id pgtype.UUID) error {
	_, err := s.Repo.DeleteProject(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrProjectNotFound
		}
		return fmt.Errorf("delete project: %w", err)
	}
	return nil
}
