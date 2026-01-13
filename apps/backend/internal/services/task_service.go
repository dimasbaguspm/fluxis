package services

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/dimasbaguspm/fluxis/internal/common"
	"github.com/dimasbaguspm/fluxis/internal/models"
	"github.com/dimasbaguspm/fluxis/internal/repositories"
	"github.com/dimasbaguspm/fluxis/internal/workers"
	"golang.org/x/sync/errgroup"
)

type TaskService struct {
	taskRepo    repositories.TaskRepository
	projectRepo repositories.ProjectRepository
	statusRepo  repositories.StatusRepository
	lr          repositories.LogRepository
	lw          *workers.LogWorker
}

func NewTaskService(taskRepo repositories.TaskRepository, projectRepo repositories.ProjectRepository, statusRepo repositories.StatusRepository, lw *workers.LogWorker, lr repositories.LogRepository) TaskService {
	return TaskService{taskRepo: taskRepo, projectRepo: projectRepo, statusRepo: statusRepo, lr: lr, lw: lw}
}

func (ts *TaskService) GetPaginated(ctx context.Context, q models.TaskSearchModel) (models.TaskPaginatedModel, error) {
	all := make([]string, 0, len(q.ID)+len(q.ProjectID)+len(q.StatusID))
	all = append(all, q.ID...)
	all = append(all, q.ProjectID...)
	all = append(all, q.StatusID...)

	for _, id := range all {
		if !common.ValidateUUID(id) {
			return models.TaskPaginatedModel{}, huma.Error400BadRequest("Must provide UUID format")
		}
	}

	return ts.taskRepo.GetPaginated(ctx, q)
}

func (ts *TaskService) GetDetail(ctx context.Context, id string) (models.TaskModel, error) {
	if !common.ValidateUUID(id) {
		return models.TaskModel{}, huma.Error400BadRequest("Must provide UUID format")
	}
	return ts.taskRepo.GetDetail(ctx, id)
}

func (ts *TaskService) Create(ctx context.Context, payload models.TaskCreateModel) (models.TaskModel, error) {
	if !(common.ValidateUUID(payload.ProjectID) && common.ValidateUUID(payload.StatusID)) {
		return models.TaskModel{}, huma.Error400BadRequest("Must provide UUID format")
	}

	var pj models.ProjectModel
	var st models.StatusModel
	g, ctxg := errgroup.WithContext(ctx)

	g.Go(func() error {
		p, err := ts.projectRepo.GetDetail(ctxg, payload.ProjectID)
		if err != nil {
			return err
		}
		pj = p
		return nil
	})
	g.Go(func() error {
		s, err := ts.statusRepo.GetDetail(ctxg, payload.StatusID)
		if err != nil {
			return err
		}
		st = s
		return nil
	})

	if err := g.Wait(); err != nil {
		return models.TaskModel{}, err
	}
	if st.ProjectID != pj.ID {
		return models.TaskModel{}, huma.Error400BadRequest("Status does not belong to the project")
	}

	t, err := ts.taskRepo.Create(ctx, payload)
	if err != nil {
		return t, err
	}
	if ts.lw != nil {
		ts.lw.Enqueue(workers.Trigger{Resource: "task", ID: t.ID, Action: "created"})
	}
	return t, nil
}

func (ts *TaskService) Update(ctx context.Context, id string, payload models.TaskUpdateModel) (models.TaskModel, error) {
	if !common.ValidateUUID(id) || (payload.StatusID != "" && !common.ValidateUUID(payload.StatusID)) {
		return models.TaskModel{}, huma.Error400BadRequest("Must provide UUID format")
	}

	if payload.StatusID != "" {
		var t models.TaskModel
		var st models.StatusModel
		g, ctxg := errgroup.WithContext(ctx)
		g.Go(func() error {
			tt, err := ts.taskRepo.GetDetail(ctxg, id)
			if err != nil {
				return err
			}
			t = tt
			return nil
		})
		g.Go(func() error {
			s, err := ts.statusRepo.GetDetail(ctxg, payload.StatusID)
			if err != nil {
				return err
			}
			st = s
			return nil
		})
		if err := g.Wait(); err != nil {
			return models.TaskModel{}, err
		}
		if st.ProjectID != t.ProjectID {
			return models.TaskModel{}, huma.Error400BadRequest("Status does not belong to the project")
		}
	}

	res, err := ts.taskRepo.Update(ctx, id, payload)
	if err != nil {
		return res, err
	}
	if ts.lw != nil {
		ts.lw.Enqueue(workers.Trigger{Resource: "task", ID: res.ID, Action: "updated"})
	}
	return res, nil
}

func (ts *TaskService) Delete(ctx context.Context, id string) error {
	if !common.ValidateUUID(id) {
		return huma.Error400BadRequest("Must provide UUID format")
	}
	if err := ts.taskRepo.Delete(ctx, id); err != nil {
		return err
	}
	if ts.lw != nil {
		ts.lw.Enqueue(workers.Trigger{Resource: "task", ID: id, Action: "deleted"})
	}
	return nil
}

func (ts *TaskService) GetLogs(ctx context.Context, projectID string, q models.LogSearchModel) (models.LogPaginatedModel, error) {
	if !common.ValidateUUID(projectID) {
		return models.LogPaginatedModel{}, huma.Error400BadRequest("Must provide UUID format")
	}
	return ts.lr.GetPaginated(ctx, projectID, q)
}
