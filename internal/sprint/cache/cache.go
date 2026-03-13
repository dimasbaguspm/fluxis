package cache

import (
	"context"

	"github.com/dimasbaguspm/fluxis/pkg/cache"
	"github.com/dimasbaguspm/fluxis/pkg/domain"
	"github.com/jackc/pgx/v5/pgtype"
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

func (sc *SprintCache) GetSingleSprint(ctx context.Context, sprintID pgtype.UUID, fetch func(context.Context) (domain.SprintModel, error)) (domain.SprintModel, error) {
	key := cache.KeySingleSprint(sc.hmacKey, sprintID)
	return cache.ReadThrough(ctx, sc.c, key, sc.cfg.DefaultTTL, func() (domain.SprintModel, error) {
		return fetch(ctx)
	})
}

func (sc *SprintCache) GetPagedSprints(ctx context.Context, params interface{}, fetch func(context.Context) (domain.SprintsPagedModel, error)) (domain.SprintsPagedModel, error) {
	key := cache.KeyPagedSprints(sc.hmacKey, params)
	return cache.ReadOrWrite(ctx, sc.c, key, sc.cfg.DefaultTTL, func(ctx context.Context) (domain.SprintsPagedModel, error) {
		return fetch(ctx)
	})
}

func (sc *SprintCache) InvalidateSingleActiveSprint(ctx context.Context, projectID pgtype.UUID) {
	_ = sc.c.Delete(ctx, cache.KeySingleActiveSprint(sc.hmacKey, projectID))
}

func (sc *SprintCache) InvalidateSingleSprint(ctx context.Context, sprintID pgtype.UUID) {
	_ = sc.c.Delete(ctx, cache.KeySingleSprint(sc.hmacKey, sprintID))
}
