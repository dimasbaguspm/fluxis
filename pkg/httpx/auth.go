package httpx

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

var authWrite domain.AuthWrite

func InitAuth(v domain.AuthWrite) {
	if authWrite != nil {
		panic("httpx.InitAuth called more than once")
	}
	authWrite = v
}

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, ok := bearerToken(r)
		if !ok {
			Error(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		claim, err := authWrite.ValidateAccessToken(r.Context(), token)
		if err != nil {
			var appErr *AppError
			if errors.As(err, &appErr) {
				ErrorCode(w, appErr.Status, appErr.Message, appErr.Code)
				return
			}
			Error(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), keyUserID, claim.ID)
		next(w, r.WithContext(ctx))
	}
}

func bearerToken(r *http.Request) (string, bool) {
	h := r.Header.Get("Authorization")
	if h == "" {
		return "", false
	}
	// SplitN to 2 — handles tokens that contain spaces (shouldn't happen but safe)
	parts := strings.SplitN(h, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return "", false
	}
	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", false
	}
	return token, true
}
