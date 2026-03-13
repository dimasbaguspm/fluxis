package user

import (
	"context"
	"log/slog"
	"net/http"

	usercache "github.com/dimasbaguspm/fluxis/internal/user/cache"
	"github.com/dimasbaguspm/fluxis/internal/user/handler"
	"github.com/dimasbaguspm/fluxis/pkg/httpx"
	"github.com/dimasbaguspm/fluxis/pkg/pubsub"
)

type Module struct {
	h          *handler.Handler
	userCache  *usercache.UserCache
	bus        pubsub.Bus
}

func NewModule(h *handler.Handler, c *usercache.UserCache, bus pubsub.Bus) *Module {
	return &Module{
		h:         h,
		userCache: c,
		bus:       bus,
	}
}

func (m *Module) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /users/me", httpx.RequireAuth(m.h.GetCurrentUser))
}

func (m *Module) StartSubscriber(ctx context.Context) {
	slog.Info("[UserModule]: starting bus subscriber")
	handler := func(ctx context.Context, e pubsub.Event) error {
		switch e.Type {
		case pubsub.UserUpdated, pubsub.UserDeleted:
			if userID, ok := pubsub.UUIDFromPayload(e, "id"); ok {
				m.userCache.InvalidateSingleUser(ctx, userID)
			}
		}
		return nil
	}

	m.bus.Subscribe(ctx, pubsub.Channel(pubsub.User), handler)
}
