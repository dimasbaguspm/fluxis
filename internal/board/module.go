package board

import (
	"context"
	"log/slog"
	"net/http"

	boardcache "github.com/dimasbaguspm/fluxis/internal/board/cache"
	"github.com/dimasbaguspm/fluxis/internal/board/handler"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
	"github.com/dimasbaguspm/fluxis/pkg/pubsub"
)

type Module struct {
	handler    *handler.Handler
	boardCache *boardcache.BoardCache
	bus        pubsub.Bus
}

func NewModule(h *handler.Handler, c *boardcache.BoardCache, bus pubsub.Bus) *Module {
	return &Module{
		handler:    h,
		boardCache: c,
		bus:        bus,
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
	slog.Info("[BoardModule]: starting bus subscriber")
	handler := func(ctx context.Context, e pubsub.Event) error {
		switch e.Type {
		case pubsub.BoardUpdated, pubsub.BoardDeleted:
			if boardID, ok := pubsub.UUIDFromPayload(e, "id"); ok {
				m.boardCache.InvalidateSingleBoard(ctx, boardID)
			}
			m.boardCache.InvalidatePagedBoardColumns(ctx)
		case pubsub.BoardColumnCreated, pubsub.BoardColumnUpdated, pubsub.BoardColumnDeleted, pubsub.BoardColumnReordered:
			m.boardCache.InvalidatePagedBoardColumns(ctx)
		}
		return nil
	}

	m.bus.Subscribe(ctx, pubsub.Channel(pubsub.Board), handler)
}
