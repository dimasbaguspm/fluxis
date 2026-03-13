package board

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/dimasbaguspm/fluxis/internal/board/handler"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
	"github.com/dimasbaguspm/fluxis/pkg/pubsub"
)

type Module struct {
	handler *handler.Handler
	bus     pubsub.Bus
}

func NewModule(h *handler.Handler, bus pubsub.Bus) *Module {
	return &Module{
		handler: h,
		bus:     bus,
	}
}

func (m *Module) Routes(mux *http.ServeMux) {
	mux.HandleFunc("POST /boards", httpx.RequireAuth(m.handler.CreateBoard))
	mux.HandleFunc("GET /boards", httpx.RequireAuth(m.handler.ListBoards))
	mux.HandleFunc("GET /boards/{boardId}", httpx.RequireAuth(m.handler.GetBoard))
	mux.HandleFunc("PATCH /boards/{boardId}", httpx.RequireAuth(m.handler.UpdateBoard))
	mux.HandleFunc("PATCH /boards/reorder", httpx.RequireAuth(m.handler.ReorderBoards))
	mux.HandleFunc("DELETE /boards/{boardId}", httpx.RequireAuth(m.handler.DeleteBoard))
	mux.HandleFunc("GET /boards/{boardId}/columns", httpx.RequireAuth(m.handler.ListBoardColumns))
	mux.HandleFunc("POST /boards/{boardId}/columns", httpx.RequireAuth(m.handler.CreateBoardColumn))
	mux.HandleFunc("PATCH /boards/{boardId}/columns/reorder", httpx.RequireAuth(m.handler.ReorderBoardColumns))
	mux.HandleFunc("PATCH /boards/{boardId}/columns/{boardColumnId}", httpx.RequireAuth(m.handler.UpdateBoardColumn))
	mux.HandleFunc("DELETE /boards/{boardId}/columns/{boardColumnId}", httpx.RequireAuth(m.handler.DeleteBoardColumn))
}

func (m *Module) StartSubscriber(ctx context.Context) {
	slog.Info("[BoardModule]: starting subscriber")
	handler := func(ctx context.Context, e pubsub.Event) error {
		slog.Info("[BoardModule]: received event", "type", string(e.Type), "payload", e.Payload)
		return nil
	}

	m.bus.Subscribe(ctx, pubsub.Channel(pubsub.Board), handler)
}
