package cache

import (
	"github.com/dimasbaguspm/fluxis/pkg/cache"
)

type UserCache struct {
	c       cache.Cache
	cfg     cache.Config
	hmacKey string
}

func New(c cache.Cache) *UserCache {
	return &UserCache{
		c:       c,
		cfg:     c.GetConfig(),
		hmacKey: c.GetConfig().HMACKey,
	}
}
