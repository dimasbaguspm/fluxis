package cache

import (
	"github.com/dimasbaguspm/fluxis/pkg/cache"
)

type BoardCache struct {
	c       cache.Cache
	cfg     cache.Config
	hmacKey string
}

func New(c cache.Cache) *BoardCache {
	return &BoardCache{
		c:       c,
		cfg:     c.GetConfig(),
		hmacKey: c.GetConfig().HMACKey,
	}
}
