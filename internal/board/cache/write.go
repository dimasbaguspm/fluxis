package cache

import (
	"context"

	"github.com/dimasbaguspm/fluxis/pkg/cache"
	"github.com/jackc/pgx/v5/pgtype"
)

func (bc *BoardCache) InvalidateSingleBoard(ctx context.Context, boardID pgtype.UUID) {
	_ = bc.c.Delete(ctx, cache.KeySingleBoard(bc.hmacKey, boardID))
}

func (bc *BoardCache) InvalidatePagedBoardColumns(ctx context.Context) {
	_ = bc.c.Delete(ctx, cache.KeyPagedBoardColumns(bc.hmacKey, nil))
}

func (bc *BoardCache) InvalidatePagedBoards(ctx context.Context) {
	_ = bc.c.Delete(ctx, cache.KeyPagedBoards(bc.hmacKey, nil))
}
