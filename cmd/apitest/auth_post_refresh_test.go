package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestAuth_Refresh_Success(t *testing.T) {
	email := randomEmail()
	displayName := "Test User"
	password := "SecurePassword123!"

	tokens := register(t, email, displayName, password)

	statusCode, resp := do[domain.AuthModel](t, "POST", "/auth/refresh", domain.AuthRefreshModel{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, "")

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil || resp.Data.AccessToken == "" {
		t.Fatal("expected valid tokens")
	}
}

func TestAuth_Refresh_InvalidToken(t *testing.T) {
	statusCode, _ := do[domain.AuthModel](t, "POST", "/auth/refresh", domain.AuthRefreshModel{
		AccessToken:  "invalid.token.here",
		RefreshToken: "invalid.refresh.here",
	}, "")

	if statusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", statusCode)
	}
}

func TestAuth_Refresh_MissingAccessToken(t *testing.T) {
	statusCode, resp := do[domain.AuthModel](t, "POST", "/auth/refresh", domain.AuthRefreshModel{
		AccessToken:  "",
		RefreshToken: "some.token.here",
	}, "")

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}

	if resp.Error == nil {
		t.Fatalf("expected error response")
	}
}

func TestAuth_Refresh_MissingRefreshToken(t *testing.T) {
	statusCode, resp := do[domain.AuthModel](t, "POST", "/auth/refresh", domain.AuthRefreshModel{
		AccessToken:  "some.token.here",
		RefreshToken: "",
	}, "")

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}

	if resp.Error == nil {
		t.Fatalf("expected error response")
	}
}

func TestAuth_Refresh_EmptyBody(t *testing.T) {
	statusCode, resp := do[domain.AuthModel](t, "POST", "/auth/refresh", map[string]interface{}{}, "")

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}

	if resp.Error == nil {
		t.Fatalf("expected error response")
	}
}
