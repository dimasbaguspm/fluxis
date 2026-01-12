package services

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/dimasbaguspm/fluxis/internal/common"
	"github.com/dimasbaguspm/fluxis/internal/models"
	"github.com/dimasbaguspm/fluxis/internal/repositories"
	"golang.org/x/sync/errgroup"
)

type TaskService struct {
	taskRepo    repositories.TaskRepository
	projectRepo repositories.ProjectRepository
	statusRepo  repositories.StatusRepository
}

func NewTaskService(taskRepo repositories.TaskRepository, projectRepo repositories.ProjectRepository, statusRepo repositories.StatusRepository) TaskService {
	return TaskService{taskRepo: taskRepo, projectRepo: projectRepo, statusRepo: statusRepo}
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

	return ts.taskRepo.Create(ctx, payload)
}

func (ts *TaskService) Update(ctx context.Context, id string, payload models.TaskUpdateModel) (models.TaskModel, error) {
	if !common.ValidateUUID(id) || (payload.StatusID != "" && !common.ValidateUUID(payload.StatusID)) {
		return models.TaskModel{}, huma.Error400BadRequest("Must provide UUID format")
	}

	if payload.StatusID != "" {
		type taskRes struct {
			t   models.TaskModel
			err error
		}
		type statusRes struct {
			st  models.StatusModel
			err error
		}

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

	return ts.taskRepo.Update(ctx, id, payload)
}

func (ts *TaskService) Delete(ctx context.Context, id string) error {
	if !common.ValidateUUID(id) {
		return huma.Error400BadRequest("Must provide UUID format")
	}
	return ts.taskRepo.Delete(ctx, id)
}
