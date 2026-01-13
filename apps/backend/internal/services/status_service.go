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
	sr repositories.StatusRepository
	pr repositories.ProjectRepository
	lr repositories.LogRepository
	sw *workers.StatusWorker
}

func NewStatusService(sr repositories.StatusRepository, sw *workers.StatusWorker, lr repositories.LogRepository, pr repositories.ProjectRepository) StatusService {
	return StatusService{sr: sr, lr: lr, sw: sw, pr: pr}
}

func (ss *StatusService) GetByProject(ctx context.Context, pId string) ([]models.StatusModel, error) {
	return ss.sr.GetByProject(ctx, pId)
}

func (ss *StatusService) Create(ctx context.Context, p models.StatusCreateModel) (models.StatusModel, error) {

	s, err := ss.sr.Create(ctx, p)
	if err != nil {
		return s, err
	}

	ss.sw.Enqueue(common.Trigger{Resource: "status", ID: s.ID, Action: "created"})
	return s, nil
}

func (ss *StatusService) Update(ctx context.Context, id string, p models.StatusUpdateModel) (models.StatusModel, error) {
	s, err := ss.sr.Update(ctx, id, p)
	if err != nil {
		return s, err
	}

	ss.sw.Enqueue(common.Trigger{Resource: "status", ID: s.ID, Action: "updated"})
	return s, nil
}

func (ss *StatusService) Delete(ctx context.Context, id string) error {
	if err := ss.sr.Delete(ctx, id); err != nil {
		return err
	}

	ss.sw.Enqueue(common.Trigger{Resource: "status", ID: id, Action: "deleted"})
	return nil
}

func (ss *StatusService) Reorder(ctx context.Context, mod models.StatusReorderModel) ([]models.StatusModel, error) {
	total, matched, err := ss.sr.ValidateReorderCounts(ctx, mod)
	if err != nil {
		return nil, err
	}
	if len(mod.IDs) != total {
		return nil, huma.Error400BadRequest("Reorder p must include all statuses for the project")
	}
	if matched != len(mod.IDs) {
		return nil, huma.Error400BadRequest("Reorder p contains invalid or out-of-project status ids")
	}

	return ss.sr.Reorder(ctx, mod)
}

func (ss *StatusService) GetDetail(ctx context.Context, id string) (models.StatusModel, error) {
	return ss.sr.GetDetail(ctx, id)
}

func (ss *StatusService) GetLogs(ctx context.Context, sId string, q models.LogSearchModel) (models.LogPaginatedModel, error) {
	s, err := ss.GetDetail(ctx, sId)
	if err != nil {
		return models.LogPaginatedModel{}, err
	}

	q.StatusID = []string{s.ID}
	q.TaskID = []string{}
	return ss.lr.GetPaginated(ctx, s.ProjectID, q)
}
