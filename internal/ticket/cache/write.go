package cache

import (
	"context"

	"github.com/dimasbaguspm/fluxis/pkg/cache"
	"github.com/jackc/pgx/v5/pgtype"
)

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
