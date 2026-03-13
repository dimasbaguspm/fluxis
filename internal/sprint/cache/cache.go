package cache

import (
	"github.com/dimasbaguspm/fluxis/pkg/cache"
)

type SprintCache struct {
	c       cache.Cache
	cfg     cache.Config
	hmacKey string
}

func New(c cache.Cache) *SprintCache {
	return &SprintCache{
		c:       c,
		cfg:     c.GetConfig(),
		hmacKey: c.GetConfig().HMACKey,
	}
}
