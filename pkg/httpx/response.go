package httpx

import (
	"encoding/json"
	"net/http"
)

// envelope is the consistent response shape for every API response.
//
// Success:  { "data": <payload>,  "meta": <optional> }
// Error:    { "error": { "message": "...", "code": "..." } }
//
// Frontend always knows where to look — no guessing between
// { message } vs { error } vs { errors } vs bare payloads.

type envelope struct {
	Data  any       `json:"data,omitempty"`
	Error *ErrBlock `json:"error,omitempty"`
}

type ErrBlock struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"` // machine-readable e.g. "email_taken"
}

func OK(w http.ResponseWriter, data any) {
	write(w, http.StatusOK, envelope{Data: data})
}

func Created(w http.ResponseWriter, data any) {
	write(w, http.StatusCreated, envelope{Data: data})
}

func Error(w http.ResponseWriter, status int, message string) {
	write(w, status, envelope{Error: &ErrBlock{Message: message}})
}

func ErrorCode(w http.ResponseWriter, status int, message, code string) {
	write(w, status, envelope{Error: &ErrBlock{Message: message, Code: code}})
}

func InternalError(w http.ResponseWriter, err error) {
	write(w, http.StatusInternalServerError, envelope{
		Error: &ErrBlock{Message: "something went wrong"},
	})
}

func write(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
