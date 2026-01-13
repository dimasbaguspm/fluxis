package workers

import (
	"context"
	"strings"
	"sync"

	"github.com/dimasbaguspm/fluxis/internal/common"
	"github.com/dimasbaguspm/fluxis/internal/models"
	"github.com/dimasbaguspm/fluxis/internal/repositories"
)

type TaskWorker struct {
	*common.Worker

	taskRepo repositories.TaskRepository
	logRepo  repositories.LogRepository

	mu        sync.RWMutex
	taskCache map[string]models.TaskModel
	ctx       context.Context
}

func NewTaskWorker(
	ctx context.Context,
	taskRepo repositories.TaskRepository,
	logRepo repositories.LogRepository,
) *TaskWorker {
	tw := &TaskWorker{
		ctx:       ctx,
		taskRepo:  taskRepo,
		logRepo:   logRepo,
		taskCache: make(map[string]models.TaskModel),
	}

	tw.Worker = common.NewWorker(ctx, tw.handle)

	return tw
}

func (tw *TaskWorker) handle(t common.Trigger) {
	switch t.Action {
	case "created":
		tw.handleCreated(t.ID)
	case "updated":
		tw.handleUpdated(t.ID)
	case "deleted":
		tw.handleDeleted(t.ID)
	case "status_changed":
		tw.handleStatusChanged(t.ID)
	}
}

func (tw *TaskWorker) handleCreated(id string) {
	task, err := tw.taskRepo.GetDetail(tw.ctx, id)
	if err != nil {
		return
	}

	tw.mu.Lock()
	tw.taskCache[id] = task
	tw.mu.Unlock()

	_ = tw.logRepo.Insert(tw.ctx, models.LogCreateModel{
		ProjectID: task.ProjectID,
		TaskID:    &task.ID,
		Entry:     "task.created",
	})
}

func (tw *TaskWorker) handleUpdated(id string) {
	current, err := tw.taskRepo.GetDetail(tw.ctx, id)
	if err != nil {
		return
	}

	tw.mu.RLock()
	previous, exists := tw.taskCache[id]
	tw.mu.RUnlock()

	if !exists {
		tw.mu.Lock()
		tw.taskCache[id] = current
		tw.mu.Unlock()
		return
	}

	var changed []string
	if current.Title != previous.Title {
		changed = append(changed, "title")
	}
	if current.Details != previous.Details {
		changed = append(changed, "details")
	}
	if current.StatusID != previous.StatusID {
		changed = append(changed, "status")
	}
	if current.Priority != previous.Priority {
		changed = append(changed, "priority")
	}
	if (current.DueDate == nil && previous.DueDate != nil) || (current.DueDate != nil && previous.DueDate == nil) {
		changed = append(changed, "dueDate")
	} else if current.DueDate != nil && previous.DueDate != nil && !current.DueDate.Equal(*previous.DueDate) {
		changed = append(changed, "dueDate")
	}

	tw.mu.Lock()
	tw.taskCache[id] = current
	tw.mu.Unlock()

	if len(changed) > 0 {
		entry := "task.updated:" + strings.Join(changed, ",")
		_ = tw.logRepo.Insert(tw.ctx, models.LogCreateModel{
			ProjectID: current.ProjectID,
			TaskID:    &current.ID,
			Entry:     entry,
		})
	}
}

func (tw *TaskWorker) handleDeleted(id string) {
	tw.mu.Lock()
	delete(tw.taskCache, id)
	tw.mu.Unlock()

	_ = tw.logRepo.Insert(tw.ctx, models.LogCreateModel{
		ProjectID: id,
		TaskID:    &id,
		Entry:     "task.deleted",
	})
}

func (tw *TaskWorker) handleStatusChanged(id string) {
	current, err := tw.taskRepo.GetDetail(tw.ctx, id)
	if err != nil {
		return
	}

	tw.mu.Lock()
	tw.taskCache[id] = current
	tw.mu.Unlock()

	_ = tw.logRepo.Insert(tw.ctx, models.LogCreateModel{
		ProjectID: current.ProjectID,
		TaskID:    &current.ID,
		StatusID:  &current.StatusID,
		Entry:     "task.status_changed",
	})
}
