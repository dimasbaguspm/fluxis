package workers

import (
	"context"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dimasbaguspm/fluxis/internal/models"
	"github.com/dimasbaguspm/fluxis/internal/repositories"
)

type Trigger struct {
	Resource string
	ID       string
	Action   string
	Meta     map[string]interface{}
}

type LogWorker struct {
	projectRepo repositories.ProjectRepository
	statusRepo  repositories.StatusRepository
	taskRepo    repositories.TaskRepository
	logRepo     repositories.LogRepository

	ch       chan Trigger
	stop     chan struct{}
	wg       sync.WaitGroup
	interval time.Duration

	mu           sync.Mutex
	projectCache map[string]models.ProjectModel
	statusCache  map[string]models.StatusModel
	taskCache    map[string]models.TaskModel
	// atomic flag: 0 running, 1 stopping
	stopping int32
}

func NewLogWorker(projectRepo repositories.ProjectRepository, statusRepo repositories.StatusRepository, taskRepo repositories.TaskRepository, logRepo repositories.LogRepository, interval time.Duration) *LogWorker {
	if interval <= 0 {
		interval = 10 * time.Second
	}
	lw := &LogWorker{
		projectRepo:  projectRepo,
		statusRepo:   statusRepo,
		taskRepo:     taskRepo,
		logRepo:      logRepo,
		ch:           make(chan Trigger, 1024),
		stop:         make(chan struct{}),
		interval:     interval,
		projectCache: make(map[string]models.ProjectModel),
		statusCache:  make(map[string]models.StatusModel),
		taskCache:    make(map[string]models.TaskModel),
	}
	lw.wg.Add(1)
	go lw.run()
	return lw
}

func (lw *LogWorker) Enqueue(t Trigger) {
	// worker is shutting down; drop trigger
	if atomic.LoadInt32(&lw.stopping) == 1 {
		return
	}
	select {
	case lw.ch <- t:
	default:
		// drop trigger if queue full
	}
}

func (lw *LogWorker) Stop() {
	if !atomic.CompareAndSwapInt32(&lw.stopping, 0, 1) {
		return
	}
	// signal run loop to stop and then wait for it to drain
	close(lw.stop)
	lw.wg.Wait()
}

func (lw *LogWorker) run() {
	defer lw.wg.Done()

	ticker := time.NewTicker(lw.interval)
	defer ticker.Stop()

	pending := make(map[string]Trigger)

	drain := func() {
		if len(pending) == 0 {
			return
		}

		for key, t := range pending {
			switch t.Resource {
			case "project":
				lw.processProject(context.Background(), t.ID, t.Action)
			case "status":
				lw.processStatus(context.Background(), t.ID, t.Action)
			case "task":
				lw.processTask(context.Background(), t.ID, t.Action)
			default:
				_ = key
			}
		}
		pending = make(map[string]Trigger)
	}

	for {
		select {
		case <-lw.stop:
			// stop accepted: drain pending and also drain channel until empty
			// stop accepting new enqueues (Enqueue checks stopping flag)
			for {
				select {
				case t := <-lw.ch:
					key := t.Resource + ":" + t.ID
					pending[key] = t
				default:
					drain()
					return
				}
			}
		case t := <-lw.ch:
			// de-duplicate by resource+id
			key := t.Resource + ":" + t.ID
			pending[key] = t
		case <-ticker.C:
			drain()
		}
	}
}

func (lw *LogWorker) processProject(ctx context.Context, id string, action string) {
	switch action {
	case "deleted":
		lw.mu.Lock()
		delete(lw.projectCache, id)
		lw.mu.Unlock()
		_ = lw.logRepo.Insert(ctx, models.LogCreateModel{ProjectID: id, Entry: "project.deleted"})
		return

	case "created":
		cur, err := lw.projectRepo.GetDetail(ctx, id)
		if err != nil {
			return
		}
		lw.mu.Lock()
		lw.projectCache[id] = cur
		lw.mu.Unlock()
		_ = lw.logRepo.Insert(ctx, models.LogCreateModel{ProjectID: id, Entry: "project.created"})
		return

	case "updated":
		cur, err := lw.projectRepo.GetDetail(ctx, id)
		if err != nil {
			return
		}

		lw.mu.Lock()
		prev, ok := lw.projectCache[id]
		lw.mu.Unlock()

		// nothing to compare
		if !ok {
			return
		}

		var changed []string
		if cur.Name != prev.Name {
			changed = append(changed, "name")
		}
		if cur.Description != prev.Description {
			changed = append(changed, "description")
		}
		if cur.Status != prev.Status {
			changed = append(changed, "status")
		}

		// update cache
		lw.mu.Lock()
		lw.projectCache[id] = cur
		lw.mu.Unlock()

		if len(changed) == 0 {
			return
		}

		entry := "project.updated:" + strings.Join(changed, ",")
		_ = lw.logRepo.Insert(ctx, models.LogCreateModel{ProjectID: id, Entry: entry})
		return
	}
}

func (lw *LogWorker) processStatus(ctx context.Context, id string, action string) {
	switch action {
	case "deleted":
		lw.mu.Lock()
		delete(lw.statusCache, id)
		lw.mu.Unlock()
		_ = lw.logRepo.Insert(ctx, models.LogCreateModel{ProjectID: id, Entry: "status.deleted"})
		return

	case "created":
		cur, err := lw.statusRepo.GetDetail(ctx, id)
		if err != nil {
			return
		}
		lw.mu.Lock()
		lw.statusCache[id] = cur
		lw.mu.Unlock()
		_ = lw.logRepo.Insert(ctx, models.LogCreateModel{ProjectID: cur.ProjectID, StatusID: &cur.ID, Entry: "status.created"})
		return

	case "updated":
		cur, err := lw.statusRepo.GetDetail(ctx, id)
		if err != nil {
			return
		}

		lw.mu.Lock()
		prev, ok := lw.statusCache[id]
		lw.mu.Unlock()

		if !ok {
			// warm cache and skip if no previous
			lw.mu.Lock()
			lw.statusCache[id] = cur
			lw.mu.Unlock()
			return
		}

		var changed []string
		if cur.Name != prev.Name {
			changed = append(changed, "name")
		}
		if cur.Position != prev.Position {
			changed = append(changed, "position")
		}
		if cur.IsDefault != prev.IsDefault {
			changed = append(changed, "isDefault")
		}

		lw.mu.Lock()
		lw.statusCache[id] = cur
		lw.mu.Unlock()

		if len(changed) == 0 {
			return
		}

		entry := "status.updated:" + strings.Join(changed, ",")
		_ = lw.logRepo.Insert(ctx, models.LogCreateModel{ProjectID: cur.ProjectID, StatusID: &cur.ID, Entry: entry})
		return
	}
}

func (lw *LogWorker) processTask(ctx context.Context, id string, action string) {
	switch action {
	case "deleted":
		lw.mu.Lock()
		delete(lw.taskCache, id)
		lw.mu.Unlock()
		_ = lw.logRepo.Insert(ctx, models.LogCreateModel{ProjectID: id, Entry: "task.deleted"})
		return

	case "created":
		cur, err := lw.taskRepo.GetDetail(ctx, id)
		if err != nil {
			return
		}
		lw.mu.Lock()
		lw.taskCache[id] = cur
		lw.mu.Unlock()
		_ = lw.logRepo.Insert(ctx, models.LogCreateModel{ProjectID: cur.ProjectID, TaskID: &cur.ID, Entry: "task.created"})
		return

	case "status_changed":
		cur, err := lw.taskRepo.GetDetail(ctx, id)
		if err != nil {
			return
		}

		lw.mu.Lock()
		lw.taskCache[id] = cur
		lw.mu.Unlock()

		// Log status change
		_ = lw.logRepo.Insert(ctx, models.LogCreateModel{
			ProjectID: cur.ProjectID,
			TaskID:    &cur.ID,
			StatusID:  &cur.StatusID,
			Entry:     "task.status_changed",
		})
		return

	case "updated":
		cur, err := lw.taskRepo.GetDetail(ctx, id)
		if err != nil {
			return
		}

		lw.mu.Lock()
		prev, ok := lw.taskCache[id]
		lw.mu.Unlock()

		if !ok {
			// warm cache and skip if no previous
			lw.mu.Lock()
			lw.taskCache[id] = cur
			lw.mu.Unlock()
			return
		}

		var changed []string
		if cur.Title != prev.Title {
			changed = append(changed, "title")
		}
		if cur.Details != prev.Details {
			changed = append(changed, "details")
		}
		if cur.StatusID != prev.StatusID {
			changed = append(changed, "status")
		}
		if cur.Priority != prev.Priority {
			changed = append(changed, "priority")
		}
		// compare due date
		if (cur.DueDate == nil && prev.DueDate != nil) || (cur.DueDate != nil && prev.DueDate == nil) {
			changed = append(changed, "dueDate")
		} else if cur.DueDate != nil && prev.DueDate != nil && !cur.DueDate.Equal(*prev.DueDate) {
			changed = append(changed, "dueDate")
		}

		lw.mu.Lock()
		lw.taskCache[id] = cur
		lw.mu.Unlock()

		if len(changed) == 0 {
			return
		}

		entry := "task.updated:" + strings.Join(changed, ",")
		_ = lw.logRepo.Insert(ctx, models.LogCreateModel{ProjectID: cur.ProjectID, TaskID: &cur.ID, Entry: entry})
		return
	}
}
