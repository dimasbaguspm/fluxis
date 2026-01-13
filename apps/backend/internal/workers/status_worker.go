package workers

import (
	"context"
	"strings"
	"sync"

	"github.com/dimasbaguspm/fluxis/internal/common"
	"github.com/dimasbaguspm/fluxis/internal/models"
	"github.com/dimasbaguspm/fluxis/internal/repositories"
)

type StatusWorker struct {
	*common.Worker

	statusRepo repositories.StatusRepository
	logRepo    repositories.LogRepository

	mu          sync.RWMutex
	statusCache map[string]models.StatusModel
	ctx         context.Context
}

func NewStatusWorker(
	ctx context.Context,
	statusRepo repositories.StatusRepository,
	logRepo repositories.LogRepository,
) *StatusWorker {
	sw := &StatusWorker{
		ctx:         ctx,
		statusRepo:  statusRepo,
		logRepo:     logRepo,
		statusCache: make(map[string]models.StatusModel),
	}

	sw.Worker = common.NewWorker(ctx, sw.handle)

	return sw
}

func (sw *StatusWorker) handle(t common.Trigger) {
	switch t.Action {
	case "created":
		sw.handleCreated(t.ID)
	case "updated":
		sw.handleUpdated(t.ID)
	case "deleted":
		sw.handleDeleted(t.ID)
	}
}

func (sw *StatusWorker) handleCreated(id string) {
	status, err := sw.statusRepo.GetDetail(sw.ctx, id)
	if err != nil {
		return
	}

	sw.mu.Lock()
	sw.statusCache[id] = status
	sw.mu.Unlock()

	_ = sw.logRepo.Insert(sw.ctx, models.LogCreateModel{
		ProjectID: status.ProjectID,
		StatusID:  &status.ID,
		Entry:     "status.created",
	})
}

func (sw *StatusWorker) handleUpdated(id string) {
	current, err := sw.statusRepo.GetDetail(sw.ctx, id)
	if err != nil {
		return
	}

	sw.mu.RLock()
	previous, exists := sw.statusCache[id]
	sw.mu.RUnlock()

	if !exists {
		sw.mu.Lock()
		sw.statusCache[id] = current
		sw.mu.Unlock()
		return
	}

	var changed []string
	if current.Name != previous.Name {
		changed = append(changed, "name")
	}
	if current.Position != previous.Position {
		changed = append(changed, "position")
	}
	if current.IsDefault != previous.IsDefault {
		changed = append(changed, "isDefault")
	}

	sw.mu.Lock()
	sw.statusCache[id] = current
	sw.mu.Unlock()

	if len(changed) > 0 {
		entry := "status.updated:" + strings.Join(changed, ",")
		_ = sw.logRepo.Insert(sw.ctx, models.LogCreateModel{
			ProjectID: current.ProjectID,
			StatusID:  &current.ID,
			Entry:     entry,
		})
	}
}

func (sw *StatusWorker) handleDeleted(id string) {
	sw.mu.RLock()
	status, exists := sw.statusCache[id]
	sw.mu.RUnlock()

	sw.mu.Lock()
	delete(sw.statusCache, id)
	sw.mu.Unlock()

	if exists {
		_ = sw.logRepo.Insert(sw.ctx, models.LogCreateModel{
			ProjectID: status.ProjectID,
			StatusID:  &status.ID,
			Entry:     "status.deleted",
		})
	}
}
