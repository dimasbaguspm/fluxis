package cache

import (
	"context"

	"github.com/dimasbaguspm/fluxis/pkg/cache"
	"github.com/dimasbaguspm/fluxis/pkg/domain"
	"github.com/jackc/pgx/v5/pgtype"
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

func (bc *BoardCache) GetSingleBoard(ctx context.Context, boardID pgtype.UUID, fetch func(context.Context) (domain.BoardModel, error)) (domain.BoardModel, error) {
	key := cache.KeySingleBoard(bc.hmacKey, boardID)
	return cache.ReadThrough(ctx, bc.c, key, bc.cfg.DefaultTTL, func() (domain.BoardModel, error) {
		return fetch(ctx)
	})
}

func (bc *BoardCache) GetPagedBoardColumns(ctx context.Context, params interface{}, fetch func(context.Context) ([]domain.BoardColumnModel, error)) ([]domain.BoardColumnModel, error) {
	key := cache.KeyPagedBoardColumns(bc.hmacKey, params)
	return cache.ReadThrough(ctx, bc.c, key, bc.cfg.DefaultTTL, func() ([]domain.BoardColumnModel, error) {
		return fetch(ctx)
	})
}

func (bc *BoardCache) GetPagedBoards(ctx context.Context, params interface{}, fetch func(context.Context) (domain.BoardsPagedModel, error)) (domain.BoardsPagedModel, error) {
	key := cache.KeyPagedBoards(bc.hmacKey, params)
	return cache.ReadOrWrite(ctx, bc.c, key, bc.cfg.DefaultTTL, func(ctx context.Context) (domain.BoardsPagedModel, error) {
		return fetch(ctx)
	})
}

func (bc *BoardCache) InvalidateSingleBoard(ctx context.Context, boardID pgtype.UUID) {
	_ = bc.c.Delete(ctx, cache.KeySingleBoard(bc.hmacKey, boardID))
}

func (bc *BoardCache) InvalidatePagedBoardColumns(ctx context.Context) {
	_ = bc.c.Delete(ctx, cache.KeyPagedBoardColumns(bc.hmacKey, nil))
}
