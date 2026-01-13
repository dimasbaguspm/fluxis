package services

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/dimasbaguspm/fluxis/internal/common"
	"github.com/dimasbaguspm/fluxis/internal/models"
	"github.com/dimasbaguspm/fluxis/internal/repositories"
	"github.com/dimasbaguspm/fluxis/internal/workers"
)

type StatusService struct {
	statusRepo repositories.StatusRepository
	lr         repositories.LogRepository
	lw         *workers.LogWorker
}

func NewStatusService(statusRepo repositories.StatusRepository, lw *workers.LogWorker, lr repositories.LogRepository) StatusService {
	return StatusService{statusRepo: statusRepo, lr: lr, lw: lw}
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
	s, err := ss.statusRepo.Create(ctx, projectId, payload)
	if err != nil {
		return s, err
	}
	if ss.lw != nil {
		ss.lw.Enqueue(workers.Trigger{Resource: "status", ID: s.ID, Action: "created"})
	}
	return s, nil
}

func (ss *StatusService) Update(ctx context.Context, id string, payload models.StatusUpdateModel) (models.StatusModel, error) {
	if !common.ValidateUUID(id) {
		return models.StatusModel{}, huma.Error400BadRequest("Must provide UUID format")
	}
	s, err := ss.statusRepo.Update(ctx, id, payload)
	if err != nil {
		return s, err
	}
	if ss.lw != nil {
		ss.lw.Enqueue(workers.Trigger{Resource: "status", ID: s.ID, Action: "updated"})
	}
	return s, nil
}

func (ss *StatusService) Delete(ctx context.Context, id string) error {
	if !common.ValidateUUID(id) {
		return huma.Error400BadRequest("Must provide UUID format")
	}
	if err := ss.statusRepo.Delete(ctx, id); err != nil {
		return err
	}
	if ss.lw != nil {
		ss.lw.Enqueue(workers.Trigger{Resource: "status", ID: id, Action: "deleted"})
	}
	return nil
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

func (ss *StatusService) GetLogs(ctx context.Context, projectID string, q models.LogSearchModel) (models.LogPaginatedModel, error) {
	if !common.ValidateUUID(projectID) {
		return models.LogPaginatedModel{}, huma.Error400BadRequest("Must provide UUID format")
	}
	return ss.lr.GetPaginated(ctx, projectID, q)
}
