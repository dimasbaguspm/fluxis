package services

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/dimasbaguspm/fluxis/internal/common"
	"github.com/dimasbaguspm/fluxis/internal/models"
	"github.com/dimasbaguspm/fluxis/internal/repositories"
)

type ProjectService struct {
	projectRepo repositories.ProjectRepository
}

func NewProjectService(projectRepo repositories.ProjectRepository) ProjectService {
	return ProjectService{projectRepo}
}

func (ps *ProjectService) GetPaginated(ctx context.Context, q models.ProjectSearchModel) (models.ProjectPaginatedModel, error) {
	for _, id := range q.ID {
		if !common.ValidateUUID(id) {
			return models.ProjectPaginatedModel{}, huma.Error400BadRequest("Must provide UUID format")
		}
	}
	return ps.projectRepo.GetPaginated(ctx, q)
}

func (ps *ProjectService) GetDetail(ctx context.Context, id string) (models.ProjectModel, error) {
	isValidID := common.ValidateUUID(id)

	if !isValidID {
		return models.ProjectModel{}, huma.Error400BadRequest("Must provide UUID format")
	}

	return ps.projectRepo.GetDetail(ctx, id)
}

func (ps *ProjectService) Create(ctx context.Context, p models.ProjectCreateModel) (models.ProjectModel, error) {
	return ps.projectRepo.Create(ctx, p)
}

func (ps *ProjectService) Update(ctx context.Context, id string, p models.ProjectUpdateModel) (models.ProjectModel, error) {
	isValidID := common.ValidateUUID(id)

	if !isValidID {
		return models.ProjectModel{}, huma.Error400BadRequest("Must provide UUID format")
	}

	return ps.projectRepo.Update(ctx, id, p)
}

func (ps *ProjectService) Delete(ctx context.Context, id string) error {
	isValidID := common.ValidateUUID(id)

	if !isValidID {
		return huma.Error400BadRequest("Must provide UUID format")
	}

	return ps.projectRepo.Delete(ctx, id)
}
