package cache

import (
	"context"

	"github.com/dimasbaguspm/fluxis/pkg/cache"
	"github.com/dimasbaguspm/fluxis/pkg/domain"
	"github.com/jackc/pgx/v5/pgtype"
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

func (uc *UserCache) GetSingleUser(ctx context.Context, userID pgtype.UUID, fetch func(context.Context) (domain.UserModel, error)) (domain.UserModel, error) {
	key := cache.KeySingleUser(uc.hmacKey, userID)
	return cache.ReadThrough(ctx, uc.c, key, uc.cfg.DefaultTTL, func() (domain.UserModel, error) {
		return fetch(ctx)
	})
}

func (uc *UserCache) InvalidateSingleUser(ctx context.Context, userID pgtype.UUID) {
	_ = uc.c.Delete(ctx, cache.KeySingleUser(uc.hmacKey, userID))
}
