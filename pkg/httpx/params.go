package httpx

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
)

// PathUUID extracts a UUID from the URL path and validates it
// Returns an error if the UUID is invalid
func PathUUID(r *http.Request, key string) (pgtype.UUID, error) {
	var id pgtype.UUID
	if err := id.Scan(r.PathValue(key)); err != nil {
		return pgtype.UUID{}, BadRequest("invalid " + key)
	}
	return id, nil
}

// PathString extracts a string from the URL path
func PathString(r *http.Request, key string) string {
	return r.PathValue(key)
}

// QueryUUID extracts and validates a UUID from query parameters
// Returns an error if the parameter is missing or invalid
func QueryUUID(r *http.Request, key string) (pgtype.UUID, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return pgtype.UUID{}, BadRequest(key + " query parameter is required")
	}

	var id pgtype.UUID
	if err := id.Scan(value); err != nil {
		return pgtype.UUID{}, BadRequest("invalid " + key)
	}
	return id, nil
}

// QueryString extracts a string from query parameters
// Returns an error if the parameter is missing
func QueryString(r *http.Request, key string) (string, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return "", BadRequest(key + " query parameter is required")
	}
	return value, nil
}

// QueryStringOptional extracts an optional string from query parameters
// Returns empty string if not provided
func QueryStringOptional(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}
