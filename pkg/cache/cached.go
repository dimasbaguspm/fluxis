package cache

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"
)

// ReadOrWrite reads from cache first. On miss, calls fetch, marshals and stores in cache, then returns.
// Cache errors are logged but never returned to the caller.
func ReadOrWrite[T any](
	ctx context.Context,
	c Cache,
	key string,
	ttl time.Duration,
	fetch func(context.Context) (T, error),
) (T, error) {
	// Try to get from cache
	cached, err := c.Get(ctx, key)
	if err == nil {
		var result T
		if err := json.Unmarshal(cached, &result); err == nil {
			return result, nil
		}
		// If unmarshal fails, fall through to fetch
	}

	// Cache miss or unmarshal error: fetch fresh
	result, err := fetch(ctx)
	if err != nil {
		var zero T
		return zero, err
	}

	// Marshal and store (errors silent)
	if data, err := json.Marshal(result); err == nil {
		if err := c.Set(ctx, key, data, ttl); err != nil {
			slog.Debug("[Cache]: failed to store in cache", "key", key, "error", err)
		}
	}

	return result, nil
}

// ReadThrough is a cache-aside pattern helper: Get from cache, fallback to fetch,
// marshal result and store back in cache, return result.
// Cache errors are logged but never returned to the caller.
func ReadThrough[T any](
	ctx context.Context,
	c Cache,
	key string,
	ttl time.Duration,
	fetch func() (T, error),
) (T, error) {
	// Try to get from cache
	cached, err := c.Get(ctx, key)
	if err == nil {
		var result T
		if err := json.Unmarshal(cached, &result); err == nil {
			return result, nil
		}
		// If unmarshal fails, fall through to fetch
	}

	// Cache miss or unmarshal error: fetch fresh
	result, err := fetch()
	if err != nil {
		var zero T
		return zero, err
	}

	// Marshal and store (errors silent)
	if data, err := json.Marshal(result); err == nil {
		if err := c.Set(ctx, key, data, ttl); err != nil {
			slog.Debug("[Cache]: failed to store in cache", "key", key, "error", err)
		}
	}

	return result, nil
}
