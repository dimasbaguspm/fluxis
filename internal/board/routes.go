package board

import (
	"net/http"

	"github.com/dimasbaguspm/fluxis/internal/board/handler"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
)

type Module struct {
	handler *handler.Handler
}

func NewModule(h *handler.Handler) *Module {
	return &Module{
		handler: h,
	}
}

func (m *Module) Routes(mux *http.ServeMux) {
	mux.HandleFunc("POST /boards", httpx.RequireAuth(m.handler.CreateBoard))
	mux.HandleFunc("GET /boards", httpx.RequireAuth(m.handler.ListBoards))
	mux.HandleFunc("GET /boards/{boardId}", httpx.RequireAuth(m.handler.GetBoard))
	mux.HandleFunc("PATCH /boards/{boardId}", httpx.RequireAuth(m.handler.UpdateBoard))
	mux.HandleFunc("PATCH /boards/reorder", httpx.RequireAuth(m.handler.ReorderBoards))
	mux.HandleFunc("DELETE /boards/{boardId}", httpx.RequireAuth(m.handler.DeleteBoard))
}
