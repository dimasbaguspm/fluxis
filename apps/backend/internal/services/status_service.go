package services

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/dimasbaguspm/fluxis/internal/common"
	"github.com/dimasbaguspm/fluxis/internal/models"
	"github.com/dimasbaguspm/fluxis/internal/repositories"
)

type StatusService struct {
	statusRepo repositories.StatusRepository
}

func NewStatusService(statusRepo repositories.StatusRepository) StatusService {
	return StatusService{statusRepo}
}

func (ss *StatusService) GetByProject(ctx context.Context, projectId string) ([]models.StatusModel, error) {
	if !common.ValidateUUID(projectId) {
		return nil, huma.Error400BadRequest("Must provide UUID format")
	}
	return ss.statusRepo.GetByProject(ctx, projectId)
}

func (ss *StatusService) Create(ctx context.Context, projectId string, payload models.StatusCreateModel) (models.StatusModel, error) {
	if !common.ValidateUUID(projectId) {
		return models.StatusModel{}, huma.Error400BadRequest("Must provide UUID format")
	}
	return ss.statusRepo.Create(ctx, projectId, payload)
}

func (ss *StatusService) Update(ctx context.Context, id string, payload models.StatusUpdateModel) (models.StatusModel, error) {
	if !common.ValidateUUID(id) {
		return models.StatusModel{}, huma.Error400BadRequest("Must provide UUID format")
	}
	return ss.statusRepo.Update(ctx, id, payload)
}

func (ss *StatusService) Delete(ctx context.Context, id string) error {
	if !common.ValidateUUID(id) {
		return huma.Error400BadRequest("Must provide UUID format")
	}
	return ss.statusRepo.Delete(ctx, id)
}

func (ss *StatusService) Reorder(ctx context.Context, projectId string, ids []string) ([]models.StatusModel, error) {
	if !common.ValidateUUID(projectId) {
		return nil, huma.Error400BadRequest("Must provide UUID format")
	}
	for _, id := range ids {
		if !common.ValidateUUID(id) {
			return nil, huma.Error400BadRequest("Must provide UUID format")
		}
	}

	total, matched, err := ss.statusRepo.ValidateReorderCounts(ctx, projectId, ids)
	if err != nil {
		return nil, err
	}
	if len(ids) != total {
		return nil, huma.Error400BadRequest("Reorder payload must include all statuses for the project")
	}
	if matched != len(ids) {
		return nil, huma.Error400BadRequest("Reorder payload contains invalid or out-of-project status ids")
	}

	return ss.statusRepo.Reorder(ctx, projectId, ids)
}

func (ss *StatusService) GetDetail(ctx context.Context, id string) (models.StatusModel, error) {
	if !common.ValidateUUID(id) {
		return models.StatusModel{}, huma.Error400BadRequest("Must provide UUID format")
	}
	return ss.statusRepo.GetDetail(ctx, id)
}
