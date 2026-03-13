package cache

import (
	"github.com/dimasbaguspm/fluxis/pkg/cache"
)

type OrgCache struct {
	c       cache.Cache
	cfg     cache.Config
	hmacKey string
}

func New(c cache.Cache) *OrgCache {
	return &OrgCache{
		c:       c,
		cfg:     c.GetConfig(),
		hmacKey: c.GetConfig().HMACKey,
	}
}
