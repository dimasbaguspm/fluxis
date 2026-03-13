package ticket

import (
	"context"
	"log/slog"
	"net/http"

	ticketcache "github.com/dimasbaguspm/fluxis/internal/ticket/cache"
	"github.com/dimasbaguspm/fluxis/internal/ticket/handler"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
	"github.com/dimasbaguspm/fluxis/pkg/pubsub"
)

type Module struct {
	h            *handler.Handler
	ticketCache  *ticketcache.TicketCache
	bus          pubsub.Bus
}

func NewModule(h *handler.Handler, c *ticketcache.TicketCache, bus pubsub.Bus) *Module {
	return &Module{
		h:           h,
		ticketCache: c,
		bus:         bus,
	}
}

func (m *Module) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /tickets", httpx.RequireAuth(m.h.ListTickets))
	mux.HandleFunc("GET /tickets/{ticketId}", httpx.RequireAuth(m.h.GetTicket))
	mux.HandleFunc("POST /tickets", httpx.RequireAuth(m.h.CreateTicket))
	mux.HandleFunc("PATCH /tickets/{ticketId}", httpx.RequireAuth(m.h.UpdateTicket))
	mux.HandleFunc("PATCH /tickets/{ticketId}/move-to-board", httpx.RequireAuth(m.h.MoveTicketToBoard))
	mux.HandleFunc("PATCH /tickets/{ticketId}/move-to-sprint", httpx.RequireAuth(m.h.MoveTicketToSprint))
	mux.HandleFunc("PATCH /tickets/{ticketId}/move-board-column", httpx.RequireAuth(m.h.MoveTicketToBoardColumn))
	mux.HandleFunc("DELETE /tickets/{ticketId}", httpx.RequireAuth(m.h.DeleteTicket))
}

func (m *Module) StartSubscriber(ctx context.Context) {
	slog.Info("[TicketModule]: starting bus subscriber")
	ticketHandler := func(ctx context.Context, e pubsub.Event) error {
		switch e.Type {
		case pubsub.TicketMovedToBoard:
			m.ticketCache.InvalidatePagedBoardTickets(ctx)
			m.ticketCache.InvalidatePagedProjectBacklog(ctx)
		case pubsub.TicketMovedToBoardColumn:
			m.ticketCache.InvalidatePagedBoardTickets(ctx)
		case pubsub.TicketMovedToSprint:
			m.ticketCache.InvalidatePagedSprintTickets(ctx)
			m.ticketCache.InvalidatePagedProjectBacklog(ctx)
		}
		return nil
	}

	sprintHandler := func(ctx context.Context, e pubsub.Event) error {
		switch e.Type {
		case pubsub.SprintCompleted:
			m.ticketCache.InvalidatePagedSprintTickets(ctx)
		}
		return nil
	}

	m.bus.Subscribe(ctx, pubsub.Channel(pubsub.Ticket), ticketHandler)
	m.bus.Subscribe(ctx, pubsub.Channel(pubsub.Sprint), sprintHandler)
}
