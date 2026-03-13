package httpx

import (
	"encoding/json"
	"log/slog"
)

// EncodePayload marshals a value to JSON and returns it as a payload map.
// The payload contains a "data" key with the JSON-encoded value.
// Errors are logged internally and an empty map is returned on failure.
// This is best-effort: encoding failures should not prevent operation.
func EncodePayload(v any) map[string]string {
	data, err := json.Marshal(v)
	if err != nil {
		slog.Warn("[httpx]: failed to encode payload", "error", err)
		return map[string]string{}
	}
	return map[string]string{"data": string(data)}
}
