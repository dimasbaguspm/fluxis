package httpx

import (
	"errors"
	"log/slog"
	"net/http"
)

type AppError struct {
	Status  int    // HTTP status code
	Message string // safe to show to the client
	Code    string // optional machine-readable code e.g. "email_taken"
	Err     error  // original error for logging — never sent to client
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NotFound(msg string) *AppError {
	return &AppError{Status: http.StatusNotFound, Message: msg}
}

func Conflict(msg string) *AppError {
	return &AppError{Status: http.StatusConflict, Message: msg}
}

func BadRequest(msg string) *AppError {
	return &AppError{Status: http.StatusBadRequest, Message: msg}
}

func Unauthorized(msg string) *AppError {
	return &AppError{Status: http.StatusUnauthorized, Message: msg}
}

func Forbidden(msg string) *AppError {
	return &AppError{Status: http.StatusForbidden, Message: msg}
}

func TooManyRequests(msg string) *AppError {
	return &AppError{Status: http.StatusTooManyRequests, Message: msg}
}

func Unprocessable(msg string) *AppError {
	return &AppError{Status: http.StatusUnprocessableEntity, Message: msg}
}

func NotImplemented(msg string) *AppError {
	return &AppError{Status: http.StatusNotImplemented, Message: msg}
}

func (e *AppError) WithCode(code string) *AppError {
	e.Code = code
	return e
}

func (e *AppError) Wrap(err error) *AppError {
	e.Err = err
	return e
}

func Handle(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		ErrorCode(w, appErr.Status, appErr.Message, appErr.Code)
		return
	}

	slog.Error("unhandled error", "error", err)
	InternalError(w, err)
}
