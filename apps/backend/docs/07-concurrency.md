# Concurrency and Goroutines

## Overview

This document explains recommended concurrency patterns for the Fluxis backend and when to use them. It aligns with the repository's RSRM architecture and provides practical guidance for writing safe, testable concurrent Go code.

## Guiding principles

- Use concurrency for operations that are independent and may block (database reads, external HTTP calls, long-running I/O). Prefer sequential code for very short, CPU-bound tasks.
- Prefer higher-level synchronization primitives (`errgroup`, `sync.WaitGroup`, buffered channels) over ad-hoc goroutine management.
- Always propagate `context.Context` into goroutines and downstream calls to enable cancellation.
- Fail fast: cancel other concurrent work when one operation returns an error.
- Avoid concurrent updates to the same database rows; do concurrency for independent reads or perform coordination inside a single transaction.

## Recommended patterns

### 1 Parallel reads with cancellation (`errgroup`)

When multiple independent reads are required before continuing, use `errgroup.WithContext(ctx)` so any failure cancels the remaining work.

Example:

```go
g, ctx := errgroup.WithContext(ctx)
var a A
var b B
g.Go(func() error { v, err := fetchA(ctx); if err != nil { return err }; a = v; return nil })
g.Go(func() error { v, err := fetchB(ctx); if err != nil { return err }; b = v; return nil })
if err := g.Wait(); err != nil {
    return err
}
// use a, b
```

This pattern is used in the codebase (for example, fetching a `Project` and a `Status` concurrently, then validating ownership).

### 2 Buffered channels for single-result handoff

When a goroutine needs to send a single result back to the caller, use a buffered channel of size 1 to avoid blocking the goroutine when the receiver is slow or when select/cancellation is used.

### 3 Keep it simple for small, fast checks

For trivial operations (UUID format checks, small slices), a sequential loop is simpler and faster than goroutines. Example: `TaskService.GetPaginated` flattens UUID slices and validates them in a single loop.

## Practical rules for this repo

- Always pass `ctx` to repository/service calls and prefer `errgroup.WithContext(ctx)` for grouped concurrent reads.
- Perform concurrency before starting a database transaction, or carefully coordinate operations inside a single transaction (do not perform concurrent writes on the same rows).
- Use buffered channels (size 1) for single-result communication and close channels when appropriate.
- Use `sync/atomic` or proper synchronization primitives when updating shared memory from multiple goroutines.

## Testing and diagnostics

- Run the race detector during development: `go test ./... -race`.
- Use `pprof` and benchmarks to identify and validate concurrency improvements.

## Examples in the codebase

- `TaskService.Create` and `TaskService.Update` use `errgroup` to fetch related resources in parallel and cancel on the first error.
- `TaskService.GetPaginated` validates UUIDs sequentially for simplicity and speed.

## References

- `golang.org/x/sync/errgroup`
- `context` package for cancellation
