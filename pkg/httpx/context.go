package httpx

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type contextKey string

const (
	keyUserID    contextKey = "user_id"
	keyRequestID contextKey = "request_id"
	keyUserAgent contextKey = "user_agent"
	keyRemoteIP  contextKey = "remote_ip"
)

func RequestIDFrom(ctx context.Context) string {
	v, _ := ctx.Value(keyRequestID).(string)
	return v
}

func UserAgentFrom(ctx context.Context) string {
	v, _ := ctx.Value(keyUserAgent).(string)
	return v
}

func UserIDFrom(ctx context.Context) (pgtype.UUID, bool) {
	id, ok := ctx.Value(keyUserID).(pgtype.UUID)
	return id, ok
}

func RemoteIPFrom(ctx context.Context) string {
	v, _ := ctx.Value(keyRemoteIP).(string)
	return v
}

func MustUserID(ctx context.Context) pgtype.UUID {
	id, ok := ctx.Value(keyUserID).(pgtype.UUID)
	if !ok {
		panic("httpx.MustUserID: called outside of RequireAuth middleware")
	}
	return id
}
