package workers

import (
	"context"
	"strings"
	"sync"

	"github.com/dimasbaguspm/fluxis/internal/common"
	"github.com/dimasbaguspm/fluxis/internal/models"
	"github.com/dimasbaguspm/fluxis/internal/repositories"
)

type ProjectWorker struct {
	*common.Worker

	projectRepo repositories.ProjectRepository
	logRepo     repositories.LogRepository

	mu           sync.RWMutex
	projectCache map[string]models.ProjectModel
	ctx          context.Context
}

func NewProjectWorker(
	ctx context.Context,
	projectRepo repositories.ProjectRepository,
	logRepo repositories.LogRepository,
) *ProjectWorker {
	pw := &ProjectWorker{
		ctx:          ctx,
		projectRepo:  projectRepo,
		logRepo:      logRepo,
		projectCache: make(map[string]models.ProjectModel),
	}

	pw.Worker = common.NewWorker(ctx, pw.handle)

	return pw
}

func (pw *ProjectWorker) handle(t common.Trigger) {
	switch t.Action {
	case "created":
		pw.handleCreated(t.ID)
	case "updated":
		pw.handleUpdated(t.ID)
	case "deleted":
		pw.handleDeleted(t.ID)
	}
}

func (pw *ProjectWorker) handleCreated(id string) {
	project, err := pw.projectRepo.GetDetail(pw.ctx, id)
	if err != nil {
		return
	}

	pw.mu.Lock()
	pw.projectCache[id] = project
	pw.mu.Unlock()

	_ = pw.logRepo.Insert(pw.ctx, models.LogCreateModel{
		ProjectID: id,
		Entry:     "project.created",
	})
}

func (pw *ProjectWorker) handleUpdated(id string) {
	current, err := pw.projectRepo.GetDetail(pw.ctx, id)
	if err != nil {
		return
	}

	pw.mu.RLock()
	previous, exists := pw.projectCache[id]
	pw.mu.RUnlock()

	if !exists {
		pw.mu.Lock()
		pw.projectCache[id] = current
		pw.mu.Unlock()
		return
	}

	var changed []string
	if current.Name != previous.Name {
		changed = append(changed, "name")
	}
	if current.Description != previous.Description {
		changed = append(changed, "description")
	}
	if current.Status != previous.Status {
		changed = append(changed, "status")
	}

	pw.mu.Lock()
	pw.projectCache[id] = current
	pw.mu.Unlock()

	if len(changed) > 0 {
		entry := "project.updated:" + strings.Join(changed, ",")
		_ = pw.logRepo.Insert(pw.ctx, models.LogCreateModel{
			ProjectID: id,
			Entry:     entry,
		})
	}
}

func (pw *ProjectWorker) handleDeleted(id string) {
	pw.mu.Lock()
	delete(pw.projectCache, id)
	pw.mu.Unlock()

	_ = pw.logRepo.Insert(pw.ctx, models.LogCreateModel{
		ProjectID: id,
		Entry:     "project.deleted",
	})
}
