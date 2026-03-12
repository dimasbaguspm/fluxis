package httpx

import (
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
)

// === PATH PARAMETERS ===

// PathString extracts a string from the URL path
func PathString(r *http.Request, key string) string {
	return r.PathValue(key)
}

// PathUUID extracts a UUID from the URL path and validates it
// Returns an error if the UUID is invalid
func PathUUID(r *http.Request, key string) (pgtype.UUID, error) {
	var id pgtype.UUID
	if err := id.Scan(r.PathValue(key)); err != nil {
		return pgtype.UUID{}, BadRequest("invalid " + key)
	}
	return id, nil
}

// === QUERY PARAMETERS - SINGLE VALUES ===
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

// QueryString retrieves a single string from query parameters
func QueryString(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

// QueryNumber retrieves a single integer from query parameters
// Returns 0 if not provided or invalid
func QueryNumber(r *http.Request, key string) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return 0
	}
	if n, err := strconv.Atoi(value); err == nil {
		return n
	}
	return 0
}

// QueryBoolean retrieves a single boolean from query parameters
// Accepts: true, 1, yes as true; everything else as false
// Returns false if not provided
func QueryBoolean(r *http.Request, key string) bool {
	value := r.URL.Query().Get(key)
	return value == "true" || value == "1" || value == "yes"
}

// === QUERY PARAMETERS - MULTIPLE VALUES (ARRAYS) ===

// QueryStrings retrieves multiple string values from query parameters
// Returns empty slice if not provided
func QueryStrings(r *http.Request, key string) []string {
	return r.URL.Query()[key]
}

// QueryNumbers retrieves multiple integer values from query parameters
// Skips invalid values, returns only valid integers
func QueryNumbers(r *http.Request, key string) []int {
	var numbers []int
	for _, val := range r.URL.Query()[key] {
		if n, err := strconv.Atoi(val); err == nil {
			numbers = append(numbers, n)
		}
	}
	return numbers
}

// QueryUUIDs retrieves multiple UUID values from query parameters
// Skips invalid UUIDs, returns only valid UUIDs
func QueryUUIDs(r *http.Request, key string) []pgtype.UUID {
	var uuids []pgtype.UUID
	for _, idStr := range r.URL.Query()[key] {
		var id pgtype.UUID
		if err := id.Scan(idStr); err == nil {
			uuids = append(uuids, id)
		}
	}
	return uuids
}

// QueryBooleans retrieves multiple boolean values from query parameters
// Accepts: true, 1, yes as true; everything else as false
func QueryBooleans(r *http.Request, key string) []bool {
	var bools []bool
	for _, val := range r.URL.Query()[key] {
		bools = append(bools, val == "true" || val == "1" || val == "yes")
	}
	return bools
}
