package cache

import (
	"context"

	"github.com/dimasbaguspm/fluxis/pkg/cache"
	"github.com/jackc/pgx/v5/pgtype"
)

func (sc *SprintCache) InvalidateSingleActiveSprint(ctx context.Context, projectID pgtype.UUID) {
	_ = sc.c.Delete(ctx, cache.KeySingleActiveSprint(sc.hmacKey, projectID))
}

func (sc *SprintCache) InvalidateSingleSprint(ctx context.Context, sprintID pgtype.UUID) {
	_ = sc.c.Delete(ctx, cache.KeySingleSprint(sc.hmacKey, sprintID))
}

func (sc *SprintCache) InvalidatePagedSprints(ctx context.Context) {
	_ = sc.c.Delete(ctx, cache.KeyPagedSprints(sc.hmacKey, nil))
}
