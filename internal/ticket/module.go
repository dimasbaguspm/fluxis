package ticket

import (
	"net/http"

	"github.com/dimasbaguspm/fluxis/internal/ticket/handler"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
)

type Module struct {
	h *handler.Handler
}

func NewModule(h *handler.Handler) *Module {
	return &Module{h}
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
