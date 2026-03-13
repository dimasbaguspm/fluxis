package cache

import (
	"context"

	"github.com/dimasbaguspm/fluxis/pkg/cache"
	"github.com/dimasbaguspm/fluxis/pkg/domain"
	"github.com/jackc/pgx/v5/pgtype"
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

func (tc *TicketCache) GetSingleTicket(ctx context.Context, ticketID pgtype.UUID, fetch func(context.Context) (domain.TicketModel, error)) (domain.TicketModel, error) {
	key := cache.KeySingleTicket(tc.hmacKey, ticketID)
	return cache.ReadThrough(ctx, tc.c, key, tc.cfg.DefaultTTL, func() (domain.TicketModel, error) {
		return fetch(ctx)
	})
}

func (tc *TicketCache) GetPagedBoardTickets(ctx context.Context, params interface{}, fetch func(context.Context) (domain.TicketsPagedModel, error)) (domain.TicketsPagedModel, error) {
	key := cache.KeyPagedBoardTickets(tc.hmacKey, params)
	return cache.ReadThrough(ctx, tc.c, key, tc.cfg.DefaultTTL, func() (domain.TicketsPagedModel, error) {
		return fetch(ctx)
	})
}

func (tc *TicketCache) GetPagedSprintTickets(ctx context.Context, params interface{}, fetch func(context.Context) (domain.TicketsPagedModel, error)) (domain.TicketsPagedModel, error) {
	key := cache.KeyPagedSprintTickets(tc.hmacKey, params)
	return cache.ReadThrough(ctx, tc.c, key, tc.cfg.DefaultTTL, func() (domain.TicketsPagedModel, error) {
		return fetch(ctx)
	})
}

func (tc *TicketCache) GetPagedProjectBacklog(ctx context.Context, params interface{}, fetch func(context.Context) (domain.TicketsPagedModel, error)) (domain.TicketsPagedModel, error) {
	key := cache.KeyPagedProjectBacklog(tc.hmacKey, params)
	return cache.ReadThrough(ctx, tc.c, key, tc.cfg.DefaultTTL, func() (domain.TicketsPagedModel, error) {
		return fetch(ctx)
	})
}

func (tc *TicketCache) InvalidateSingleTicket(ctx context.Context, ticketID pgtype.UUID) {
	_ = tc.c.Delete(ctx, cache.KeySingleTicket(tc.hmacKey, ticketID))
}

func (tc *TicketCache) InvalidatePagedBoardTickets(ctx context.Context) {
	_ = tc.c.Delete(ctx, cache.KeyPagedBoardTickets(tc.hmacKey, nil))
}

func (tc *TicketCache) InvalidatePagedSprintTickets(ctx context.Context) {
	_ = tc.c.Delete(ctx, cache.KeyPagedSprintTickets(tc.hmacKey, nil))
}

func (tc *TicketCache) InvalidatePagedProjectBacklog(ctx context.Context) {
	_ = tc.c.Delete(ctx, cache.KeyPagedProjectBacklog(tc.hmacKey, nil))
}
