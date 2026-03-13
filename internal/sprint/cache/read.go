package cache

import (
	"context"

	"github.com/dimasbaguspm/fluxis/pkg/cache"
	"github.com/dimasbaguspm/fluxis/pkg/domain"
	"github.com/jackc/pgx/v5/pgtype"
)

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
