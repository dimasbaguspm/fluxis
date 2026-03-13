package cache

import (
	"context"
	"sync"
	"time"
)

type MemoryCache struct {
	mu    sync.RWMutex
	cache map[string]*cacheEntry
	cfg   Config
	done  chan struct{}
}

type cacheEntry struct {
	value  []byte
	expiry time.Time
}

func New(cfg Config) *MemoryCache {
	mc := &MemoryCache{
		cache: make(map[string]*cacheEntry),
		cfg:   cfg,
		done:  make(chan struct{}),
	}
	go mc.cleanupLoop()
	return mc
}

func (m *MemoryCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if ttl == 0 {
		ttl = m.cfg.DefaultTTL
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cache[key] = &cacheEntry{
		value:  value,
		expiry: time.Now().Add(ttl),
	}
	return nil
}

func (m *MemoryCache) Get(ctx context.Context, key string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	entry, ok := m.cache[key]
	if !ok {
		return nil, ErrMiss
	}
	if time.Now().After(entry.expiry) {
		return nil, ErrMiss
	}
	return entry.value, nil
}

func (m *MemoryCache) Delete(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.cache, key)
	return nil
}

func (m *MemoryCache) GetConfig() Config {
	return m.cfg
}

func (m *MemoryCache) cleanupLoop() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.cleanup()
		case <-m.done:
			return
		}
	}
}

func (m *MemoryCache) cleanup() {
	const batchSize = 100

	for {
		m.mu.Lock()
		batch := make([]string, 0, batchSize)
		now := time.Now()

		for key, entry := range m.cache {
			if now.After(entry.expiry) {
				batch = append(batch, key)
				if len(batch) >= batchSize {
					break
				}
			}
		}

		for _, key := range batch {
			delete(m.cache, key)
		}
		m.mu.Unlock()

		if len(batch) < batchSize {
			break
		}
	}
}

func (m *MemoryCache) Close() error {
	close(m.done)
	return nil
}
