package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestAuth_Login_Success(t *testing.T) {
	email := randomEmail()
	displayName := "Test User"
	password := "SecurePassword123!"

	register(t, email, displayName, password)

	statusCode, resp := do[domain.AuthModel](t, "POST", "/auth/login", domain.AuthLoginModel{
		Email:    email,
		Password: password,
	}, "")

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil || resp.Data.AccessToken == "" {
		t.Fatal("expected valid tokens")
	}
}

func TestAuth_Login_WrongPassword(t *testing.T) {
	email := randomEmail()
	displayName := "Test User"
	password := "SecurePassword123!"

	register(t, email, displayName, password)

	statusCode, _ := do[domain.AuthModel](t, "POST", "/auth/login", domain.AuthLoginModel{
		Email:    email,
		Password: "WrongPassword123!",
	}, "")

	if statusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", statusCode)
	}
}

func TestAuth_Login_NonExistentUser(t *testing.T) {
	statusCode, _ := do[domain.AuthModel](t, "POST", "/auth/login", domain.AuthLoginModel{
		Email:    randomEmail(),
		Password: "SomePassword123!",
	}, "")

	if statusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", statusCode)
	}
}

func TestAuth_Login_MissingEmail(t *testing.T) {
	statusCode, resp := do[domain.AuthModel](t, "POST", "/auth/login", domain.AuthLoginModel{
		Email:    "",
		Password: "SomePassword123!",
	}, "")

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}

	if resp.Error == nil {
		t.Fatalf("expected error response")
	}
}

func TestAuth_Login_MissingPassword(t *testing.T) {
	statusCode, resp := do[domain.AuthModel](t, "POST", "/auth/login", domain.AuthLoginModel{
		Email:    randomEmail(),
		Password: "",
	}, "")

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}

	if resp.Error == nil {
		t.Fatalf("expected error response")
	}
}

func TestAuth_Login_InvalidEmailFormat(t *testing.T) {
	statusCode, resp := do[domain.AuthModel](t, "POST", "/auth/login", domain.AuthLoginModel{
		Email:    "not-an-email",
		Password: "SomePassword123!",
	}, "")

	if statusCode != http.StatusBadRequest && statusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 400 or 401, got %d", statusCode)
	}

	if resp.Error == nil {
		t.Fatalf("expected error response")
	}
}
