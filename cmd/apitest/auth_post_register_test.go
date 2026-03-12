package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestAuth_Register_Success(t *testing.T) {
	email := randomEmail()
	displayName := "Test User"
	password := "SecurePassword123!"

	statusCode, resp := do[domain.AuthModel](t, "POST", "/auth/register", domain.AuthRegisterModel{
		UserCreateModel: domain.UserCreateModel{
			Email:       email,
			DisplayName: displayName,
			Password:    password,
		},
	}, "")

	if statusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected data, got nil")
	}

	if resp.Data.AccessToken == "" {
		t.Fatal("expected non-empty access token")
	}

	if resp.Data.RefreshToken == "" {
		t.Fatal("expected non-empty refresh token")
	}
}

func TestAuth_Register_DuplicateEmail(t *testing.T) {
	email := randomEmail()
	displayName := "Test User"
	password := "SecurePassword123!"

	// First registration should succeed
	register(t, email, displayName, password)

	// Second registration with same email should fail
	statusCode, resp := do[domain.AuthModel](t, "POST", "/auth/register", domain.AuthRegisterModel{
		UserCreateModel: domain.UserCreateModel{
			Email:       email,
			DisplayName: "Another User",
			Password:    password,
		},
	}, "")

	if statusCode != http.StatusConflict {
		t.Fatalf("expected status 409, got %d", statusCode)
	}

	if resp.Error == nil || resp.Error.Code != "email_taken" {
		t.Fatalf("expected error code 'email_taken', got %v", resp.Error)
	}
}

func TestAuth_Register_InvalidEmail(t *testing.T) {
	statusCode, resp := do[domain.AuthModel](t, "POST", "/auth/register", domain.AuthRegisterModel{
		UserCreateModel: domain.UserCreateModel{
			Email:       "not-an-email",
			DisplayName: "Test User",
			Password:    "SecurePassword123!",
		},
	}, "")

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}

	if resp.Error == nil {
		t.Fatalf("expected error response")
	}
}

func TestAuth_Register_MissingPassword(t *testing.T) {
	statusCode, resp := do[domain.AuthModel](t, "POST", "/auth/register", domain.AuthRegisterModel{
		UserCreateModel: domain.UserCreateModel{
			Email:       randomEmail(),
			DisplayName: "Test User",
			Password:    "",
		},
	}, "")

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}

	if resp.Error == nil {
		t.Fatalf("expected error response")
	}
}

func TestAuth_Register_MissingEmail(t *testing.T) {
	statusCode, resp := do[domain.AuthModel](t, "POST", "/auth/register", domain.AuthRegisterModel{
		UserCreateModel: domain.UserCreateModel{
			Email:       "",
			DisplayName: "Test User",
			Password:    "SecurePassword123!",
		},
	}, "")

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}

	if resp.Error == nil {
		t.Fatalf("expected error response")
	}
}

func TestAuth_Register_MissingDisplayName(t *testing.T) {
	statusCode, resp := do[domain.AuthModel](t, "POST", "/auth/register", domain.AuthRegisterModel{
		UserCreateModel: domain.UserCreateModel{
			Email:       randomEmail(),
			DisplayName: "",
			Password:    "SecurePassword123!",
		},
	}, "")

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}

	if resp.Error == nil {
		t.Fatalf("expected error response")
	}
}

func TestAuth_Register_EmptyBody(t *testing.T) {
	statusCode, resp := do[domain.AuthModel](t, "POST", "/auth/register", map[string]interface{}{}, "")

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}

	if resp.Error == nil {
		t.Fatalf("expected error response")
	}
}
