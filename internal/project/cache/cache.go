package cache

import (
	"github.com/dimasbaguspm/fluxis/pkg/cache"
)

type ProjectCache struct {
	c       cache.Cache
	cfg     cache.Config
	hmacKey string
}

func New(c cache.Cache) *ProjectCache {
	return &ProjectCache{
		c:       c,
		cfg:     c.GetConfig(),
		hmacKey: c.GetConfig().HMACKey,
	}
}
