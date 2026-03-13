package user

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/dimasbaguspm/fluxis/internal/user/handler"
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
	mux.HandleFunc("GET /users/me", httpx.RequireAuth(m.h.GetCurrentUser))
}

func (m *Module) StartSubscriber(ctx context.Context) {
	slog.Info("[UserModule]: starting bus subscriber")
	handler := func(ctx context.Context, e pubsub.Event) error {
		slog.Info("[UserModule]: received event", "type", string(e.Type), "payload", e.Payload)
		return nil
	}

	m.bus.Subscribe(ctx, pubsub.Channel(pubsub.User), handler)
}
