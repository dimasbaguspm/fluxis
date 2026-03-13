package ticket

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/dimasbaguspm/fluxis/internal/ticket/handler"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
	"github.com/dimasbaguspm/fluxis/pkg/pubsub"
)

type Module struct {
	h   *handler.Handler
	bus pubsub.Bus
}

func NewModule(h *handler.Handler, bus pubsub.Bus) *Module {
	return &Module{h, bus}
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
	slog.Info("[TicketModule]: starting subscriber")
	handler := func(ctx context.Context, e pubsub.Event) error {
		slog.Info("[TicketModule]: received event", "type", string(e.Type), "payload", e.Payload)
		return nil
	}

	m.bus.Subscribe(ctx, pubsub.Channel(pubsub.Ticket), handler)
}
