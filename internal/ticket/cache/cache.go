package cache

import (
	"github.com/dimasbaguspm/fluxis/pkg/cache"
)

type TicketCache struct {
	c       cache.Cache
	cfg     cache.Config
	hmacKey string
}

func New(c cache.Cache) *TicketCache {
	return &TicketCache{
		c:       c,
		cfg:     c.GetConfig(),
		hmacKey: c.GetConfig().HMACKey,
	}
}
