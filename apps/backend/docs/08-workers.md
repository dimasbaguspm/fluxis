# Worker Pattern and Event Processing

## Overview

This document explains the worker pattern used in Fluxis for asynchronous event processing and entity change tracking. Workers enable decoupled, non-blocking event handling while maintaining data consistency through intelligent caching and batching.

## Guiding Principles

- **Decouple operations** - Workers process events asynchronously, preventing blocking of critical request paths
- **Batch events** - Multiple triggers for the same entity are deduplicated and batched to reduce database calls
- **Cache intelligently** - Workers maintain entity caches for change detection and efficient updates
- **Fail gracefully** - Worker failures don't crash the application; events are logged for auditing
- **Coordinate shutdown** - Workers drain remaining events before application termination
- **Context awareness** - Workers operate within application context for proper lifecycle management

## Architecture

### Components

**1. Base Worker** (`common.Worker`)

- Generic event queue and batch processor
- Handles trigger enqueueing, deduplication, and periodic draining
- Exposes `Enqueue()` and `Stop()` methods
- Accepts a `Handler` function for custom business logic

**2. Entity Workers**

- `ProjectWorker` - Handles project events (created, updated, deleted)
- `StatusWorker` - Handles status events (created, updated, deleted)
- `TaskWorker` - Handles task events (created, updated, deleted, status_changed)
- Each worker embeds `common.Worker` and implements its own handler

**3. Trigger System**

```go
type Trigger struct {
    Resource string                 // "project", "status", "task"
    ID       string                 // Entity ID
    Action   string                 // "created", "updated", "deleted", etc.
    Meta     map[string]interface{} // Optional metadata
}
```

## Worker Lifecycle

### 1. Initialization

Workers are created during application startup in `handler.go`:

```go
pW := workers.NewProjectWorker(ctx, pR, lR)
sW := workers.NewStatusWorker(ctx, sR, lR)
tW := workers.NewTaskWorker(ctx, tR, lR)
```

Each worker:

- Receives the application context for lifecycle management
- Gets injected with required repositories
- Starts an internal goroutine for event processing
- Runs with a configurable interval (default 10 seconds)

### 2. Event Enqueueing

Services enqueue events after entity changes:

```go
// In ProjectService.Create()
proj, err := ps.pr.Create(ctx, p)
if err != nil {
    return proj, err
}
if ps.pw != nil {
    ps.pw.Enqueue(common.Trigger{
        Resource: "project",
        ID: proj.ID,
        Action: "created",
    })
}
```

The enqueue operation is **non-blocking** - events are dropped if the queue is full (1024 items).

### 3. Batch Processing

Workers batch events based on a timer interval:

1. **Accumulation phase** - Triggers arrive and are deduplicated by `resource:id` key
2. **Timer tick** - After the configured interval, drain all pending triggers
3. **Processing phase** - Execute business logic for each unique trigger
4. **Cache update** - Update internal caches and log entries

Deduplication example:

```
Resource: "task", ID: "123", Action: "updated" (arrives)
Resource: "task", ID: "123", Action: "updated" (arrives again)
→ Single "updated" trigger processed for task 123
```

### 4. Change Detection

Workers compare cached state with current state to identify changes:

```go
// In StatusWorker.handleUpdated()
current, err := sw.statusRepo.GetDetail(ctx, id)
prev, exists := sw.statusCache[id]

var changed []string
if current.Name != prev.Name {
    changed = append(changed, "name")
}
if current.Position != prev.Position {
    changed = append(changed, "position")
}
```

Only log entries are created when actual changes are detected.

### 5. Graceful Shutdown

On context cancellation, workers drain remaining events:

```go
go func() {
    <-ctx.Done()
    pW.Stop()
    sW.Stop()
    tW.Stop()
}()
```

The `Stop()` method:

1. Sets atomic flag to prevent new enqueues
2. Signals the run loop to stop
3. Drains the event channel until empty
4. Processes all pending triggers
5. Returns when complete

## Use Cases

### Logging and Auditing

Workers create audit logs for:

- Entity creation
- Specific field changes (tracked by name)
- Entity deletion

Example log entries:

```
"project.created"
"task.updated:title,priority"
"status.deleted"
"task.status_changed"
```

### Change Tracking

Workers cache entity state to detect what changed:

- Enables targeted logging of modifications
- Reduces unnecessary operations
- Provides context for event processing

### Future Extensions

The worker pattern supports:

- **Notifications** - Trigger notifications on specific changes
- **Search indexing** - Update search indices asynchronously
- **External systems** - Sync changes to third-party services
- **Analytics** - Aggregate usage data without blocking requests

## Implementation Guidelines

### When to Create a New Worker

Create a new worker when you need to:

- Process events for a specific entity type
- Maintain a cache of entity state
- Track changes and log them
- Perform async operations that shouldn't block requests

### Worker Implementation Pattern

```go
type MyWorker struct {
    *common.Worker
    myRepo   repositories.MyRepository
    logRepo  repositories.LogRepository
    mu       sync.RWMutex
    myCache  map[string]models.MyModel
    ctx      context.Context
}

func NewMyWorker(ctx context.Context, myRepo repositories.MyRepository,
    logRepo repositories.LogRepository, interval time.Duration) *MyWorker {
    mw := &MyWorker{
        ctx:     ctx,
        myRepo:  myRepo,
        logRepo: logRepo,
        myCache: make(map[string]models.MyModel),
    }
    mw.Worker = common.NewWorker(ctx, mw.handle, interval)
    return mw
}

func (mw *MyWorker) handle(t common.Trigger) {
    switch t.Action {
    case "created":
        mw.handleCreated(t.ID)
    case "updated":
        mw.handleUpdated(t.ID)
    case "deleted":
        mw.handleDeleted(t.ID)
    }
}
```

### Thread Safety

- Use `sync.RWMutex` for cache access
- Read lock when checking cache: `mw.mu.RLock()`
- Write lock when updating cache: `mw.mu.Lock()`
- Always defer unlock: `defer mw.mu.RUnlock()`

### Error Handling

Workers silently fail on errors (no panic):

```go
status, err := sw.statusRepo.GetDetail(ctx, id)
if err != nil {
    return  // silently drop event
}
```

This prevents worker crashes from blocking the application.

## Performance Considerations

- **Queue size** - Workers use a 1024-item buffered channel. High-throughput systems may need tuning
- **Batch interval** - Default 10 seconds balances latency vs. throughput
- **Cache size** - Caches grow with active entities; add cleanup for long-running applications
- **Lock contention** - RWMutex provides efficient concurrent reads; minimize write lock duration

## Testing

Test workers by:

1. Enqueueing triggers
2. Verifying cache state changes
3. Checking log entries were created
4. Testing graceful shutdown with pending events

Example:

```go
pw := workers.NewProjectWorker(ctx, projectRepo, logRepo, 100*time.Millisecond)
pw.Enqueue(common.Trigger{Resource: "project", ID: "123", Action: "created"})
// Give worker time to process
time.Sleep(200*time.Millisecond)
// Assert cache was populated
assert.NotNil(pw.GetCachedProject("123"))
```

## References

- [07-concurrency.md](07-concurrency.md) - Concurrency patterns and guidelines
- [03-rsrm-pattern.md](03-rsrm-pattern.md) - Service layer responsibilities
- [common.Worker](../internal/common/worker.go) - Base worker implementation
- [workers/](../internal/workers/) - Entity-specific worker implementations
