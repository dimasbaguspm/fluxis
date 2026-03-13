package cache

import (
	"context"

	"github.com/dimasbaguspm/fluxis/pkg/cache"
	"github.com/jackc/pgx/v5/pgtype"
)

func (uc *UserCache) InvalidateSingleUser(ctx context.Context, userID pgtype.UUID) {
	_ = uc.c.Delete(ctx, cache.KeySingleUser(uc.hmacKey, userID))
}
