package auth

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/dimasbaguspm/fluxis/internal/auth/handler"
	"github.com/dimasbaguspm/fluxis/internal/auth/service"
	"github.com/dimasbaguspm/fluxis/pkg/pubsub"
)

type Module struct {
	h   *handler.Handler
	svc *service.Service
	bus pubsub.Bus
}

func NewModule(svc *service.Service, h *handler.Handler, bus pubsub.Bus) *Module {
	return &Module{svc: svc, h: h, bus: bus}
}

func (m *Module) Routes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/register", m.h.Register)
	mux.HandleFunc("POST /auth/login", m.h.Login)
	mux.HandleFunc("POST /auth/refresh", m.h.Refresh)
}

func (m *Module) Service() *service.Service {
	return m.svc
}

func (m *Module) StartSubscriber(ctx context.Context) {
	slog.Info("[AuthModule]: starting bus subscriber")
	handler := func(ctx context.Context, e pubsub.Event) error {
		slog.Info("[AuthModule]: received event", "type", string(e.Type), "payload", e.Payload)
		return nil
	}

	m.bus.Subscribe(ctx, pubsub.Channel(pubsub.Auth), handler)
}
