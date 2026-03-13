package ratelimit

import (
	"context"
	"encoding/binary"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dimasbaguspm/fluxis/pkg/cache"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
)

type Config struct {
	MaxRequests int
	Window      time.Duration
}

type Middleware struct {
	c   cache.Cache
	mu  sync.Mutex
	cfg Config
}

func New(cfg Config) *Middleware {
	c := cache.New(cache.Config{DefaultTTL: cfg.Window * 2})
	return &Middleware{
		c:   c,
		cfg: cfg,
	}
}

func (m *Middleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := m.clientIP(r)
		now := time.Now()
		bucket := now.Truncate(m.cfg.Window).Unix()
		ttl := time.Until(now.Truncate(m.cfg.Window).Add(m.cfg.Window))

		key := m.rateLimitKey(ip, bucket)

		m.mu.Lock()
		count := m.getCount(r.Context(), key)
		count++
		m.setCount(r.Context(), key, count, ttl)
		m.mu.Unlock()

		if count > uint32(m.cfg.MaxRequests) {
			w.Header().Set("Retry-After", strconv.FormatInt(int64(ttl.Seconds()), 10))
			httpx.Handle(w, httpx.TooManyRequests("rate limit exceeded"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) clientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}

	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	if r.RemoteAddr != "" {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err == nil {
			return host
		}
		return r.RemoteAddr
	}

	return "unknown"
}

func (m *Middleware) rateLimitKey(ip string, bucket int64) string {
	return "rate-limit:" + ip + ":" + strconv.FormatInt(bucket, 10)
}

func (m *Middleware) getCount(ctx context.Context, key string) uint32 {
	data, err := m.c.Get(ctx, key)
	if err != nil {
		return 0
	}
	if len(data) < 4 {
		return 0
	}
	return binary.BigEndian.Uint32(data[:4])
}

func (m *Middleware) setCount(ctx context.Context, key string, count uint32, ttl time.Duration) {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, count)
	m.c.Set(ctx, key, data, ttl)
}
