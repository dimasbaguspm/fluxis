package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// Decode decodes JSON body into dst without validation.
// Body is limited to 1MB — prevents memory exhaustion attacks.
func Decode(r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(nil, r.Body, 1<<20) // 1MB

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() // reject unexpected fields

	if err := dec.Decode(dst); err != nil {
		return handleDecodeError(err)
	}

	return nil
}

// DecodeAndValidate decodes JSON body into dst and runs struct validation.
// Returns a clean user-facing error string on failure.
// Body is limited to 1MB — prevents memory exhaustion attacks.
func DecodeAndValidate(r *http.Request, dst any) error {
	if err := Decode(r, dst); err != nil {
		return err
	}

	if err := validate.Struct(dst); err != nil {
		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			return formatValidationErrors(validationErrs)
		}
		return err
	}

	return nil
}

// handleDecodeError converts JSON decode errors into user-friendly messages.
func handleDecodeError(err error) error {
	var syntaxErr *json.SyntaxError
	var unmarshalErr *json.UnmarshalTypeError
	var maxBytesErr *http.MaxBytesError

	switch {
	case errors.As(err, &syntaxErr):
		return fmt.Errorf("malformed json at position %d", syntaxErr.Offset)
	case errors.As(err, &unmarshalErr):
		return fmt.Errorf("invalid type for field %q", unmarshalErr.Field)
	case errors.As(err, &maxBytesErr):
		return fmt.Errorf("request body too large")
	case errors.Is(err, io.EOF):
		return fmt.Errorf("request body is empty")
	default:
		return fmt.Errorf("decode error: %w", err)
	}
}

// formatValidationErrors turns validator's verbose errors
// into a single human-readable string.
// e.g. "email: must be a valid email address; password: minimum length is 8"

func formatValidationErrors(errs validator.ValidationErrors) error {
	msgs := make([]string, 0, len(errs))
	for _, e := range errs {
		field := strings.ToLower(e.Field())
		switch e.Tag() {
		case "required":
			msgs = append(msgs, fmt.Sprintf("%s is required", field))
		case "email":
			msgs = append(msgs, fmt.Sprintf("%s must be a valid email address", field))
		case "min":
			msgs = append(msgs, fmt.Sprintf("%s minimum length is %s", field, e.Param()))
		case "max":
			msgs = append(msgs, fmt.Sprintf("%s maximum length is %s", field, e.Param()))
		case "oneof":
			msgs = append(msgs, fmt.Sprintf("%s must be one of: %s", field, e.Param()))
		case "len":
			msgs = append(msgs, fmt.Sprintf("%s must be exactly %s characters", field, e.Param()))
		case "numeric":
			msgs = append(msgs, fmt.Sprintf("%s must contain only digits", field))
		case "url":
			msgs = append(msgs, fmt.Sprintf("%s must be a valid URL", field))
		default:
			msgs = append(msgs, fmt.Sprintf("%s is invalid (%s)", field, e.Tag()))
		}
	}
	return fmt.Errorf("%s", strings.Join(msgs, "; "))
}
