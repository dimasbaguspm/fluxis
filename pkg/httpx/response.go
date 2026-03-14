package httpx

import (
	"encoding/json"
	"net/http"
)

// Error responses use an envelope with ErrBlock
// Success responses write data directly (no envelope)
//
// Success:  <payload> (written directly)
// Error:    { "error": { "message": "...", "code": "..." } }

type errorEnvelope struct {
	Error *ErrBlock `json:"error"`
}

type ErrBlock struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"` // machine-readable e.g. "email_taken"
}

func OK(w http.ResponseWriter, data any) {
	write(w, http.StatusOK, data)
}

func Created(w http.ResponseWriter, data any) {
	write(w, http.StatusCreated, data)
}

func Error(w http.ResponseWriter, status int, message string) {
	write(w, status, errorEnvelope{Error: &ErrBlock{Message: message}})
}

func ErrorCode(w http.ResponseWriter, status int, message, code string) {
	write(w, status, errorEnvelope{Error: &ErrBlock{Message: message, Code: code}})
}

func InternalError(w http.ResponseWriter, err error) {
	write(w, http.StatusInternalServerError, errorEnvelope{
		Error: &ErrBlock{Message: "something went wrong"},
	})
}

func write(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
