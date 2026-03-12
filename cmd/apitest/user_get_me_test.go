package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestUser_GetMe_Success(t *testing.T) {
	email := randomEmail()
	displayName := "Test User"
	password := "SecurePassword123!"

	tokens := register(t, email, displayName, password)

	statusCode, resp := do[domain.UserModel](t, "GET", "/users/me", nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected user data")
	}

	if resp.Data.Email != email {
		t.Fatalf("expected email %s, got %s", email, resp.Data.Email)
	}

	if resp.Data.DisplayName != displayName {
		t.Fatalf("expected display name %s, got %s", displayName, resp.Data.DisplayName)
	}
}

func TestUser_GetMe_Unauthenticated(t *testing.T) {
	statusCode, _ := do[domain.UserModel](t, "GET", "/users/me", nil, "")

	if statusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", statusCode)
	}
}

func TestUser_GetMe_InvalidToken(t *testing.T) {
	statusCode, _ := do[domain.UserModel](t, "GET", "/users/me", nil, "invalid.token.here")

	if statusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", statusCode)
	}
}
