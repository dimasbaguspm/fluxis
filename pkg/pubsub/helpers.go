package pubsub

import (
	"github.com/jackc/pgx/v5/pgtype"
)

func UUIDFromPayload(e Event, key string) (pgtype.UUID, bool) {
	val, ok := e.Payload[key]
	if !ok {
		return pgtype.UUID{}, false
	}
	var id pgtype.UUID
	if err := id.Scan(val); err != nil {
		return pgtype.UUID{}, false
	}
	return id, true
}

func StringFromPayload(e Event, key string) (string, bool) {
	val, ok := e.Payload[key]
	return val, ok
}
