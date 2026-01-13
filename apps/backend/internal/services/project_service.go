package services

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/dimasbaguspm/fluxis/internal/common"
	"github.com/dimasbaguspm/fluxis/internal/models"
	"github.com/dimasbaguspm/fluxis/internal/repositories"
	"github.com/dimasbaguspm/fluxis/internal/workers"
)

type ProjectService struct {
	pr repositories.ProjectRepository
	lw *workers.LogWorker
	lr repositories.LogRepository
}

func NewProjectService(pr repositories.ProjectRepository, lw *workers.LogWorker, lr repositories.LogRepository) ProjectService {
	return ProjectService{pr: pr, lw: lw, lr: lr}
}

func (ps *ProjectService) GetPaginated(ctx context.Context, q models.ProjectSearchModel) (models.ProjectPaginatedModel, error) {
	for _, id := range q.ID {
		if !common.ValidateUUID(id) {
			return models.ProjectPaginatedModel{}, huma.Error400BadRequest("Must provide UUID format")
		}
	}
	return ps.pr.GetPaginated(ctx, q)
}

func (ps *ProjectService) GetDetail(ctx context.Context, id string) (models.ProjectModel, error) {
	isValidID := common.ValidateUUID(id)

	if !isValidID {
		return models.ProjectModel{}, huma.Error400BadRequest("Must provide UUID format")
	}

	return ps.pr.GetDetail(ctx, id)
}

func (ps *ProjectService) Create(ctx context.Context, p models.ProjectCreateModel) (models.ProjectModel, error) {
	proj, err := ps.pr.Create(ctx, p)
	if err != nil {
		return proj, err
	}

	if ps.lw != nil {
		ps.lw.Enqueue(workers.Trigger{Resource: "project", ID: proj.ID, Action: "created"})
	}

	return proj, nil
}

func (ps *ProjectService) Update(ctx context.Context, id string, p models.ProjectUpdateModel) (models.ProjectModel, error) {
	isValidID := common.ValidateUUID(id)

	if !isValidID {
		return models.ProjectModel{}, huma.Error400BadRequest("Must provide UUID format")
	}

	proj, err := ps.pr.Update(ctx, id, p)
	if err != nil {
		return proj, err
	}

	if ps.lw != nil {
		ps.lw.Enqueue(workers.Trigger{Resource: "project", ID: proj.ID, Action: "updated"})
	}

	return proj, nil
}

func (ps *ProjectService) Delete(ctx context.Context, id string) error {
	isValidID := common.ValidateUUID(id)

	if !isValidID {
		return huma.Error400BadRequest("Must provide UUID format")
	}

	if err := ps.pr.Delete(ctx, id); err != nil {
		return err
	}

	if ps.lw != nil {
		ps.lw.Enqueue(workers.Trigger{Resource: "project", ID: id, Action: "deleted"})
	}

	return nil
}

func (ps *ProjectService) GetLogs(ctx context.Context, projectID string, q models.LogSearchModel) (models.LogPaginatedModel, error) {
	if !common.ValidateUUID(projectID) {
		return models.LogPaginatedModel{}, huma.Error400BadRequest("Must provide UUID format")
	}
	return ps.lr.GetPaginated(ctx, projectID, q)
}
