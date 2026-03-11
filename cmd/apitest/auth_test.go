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

func TestUsers_GetMe_Authenticated(t *testing.T) {
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

func TestUsers_GetMe_Unauthenticated(t *testing.T) {
	statusCode, _ := do[domain.UserModel](t, "GET", "/users/me", nil, "")

	if statusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", statusCode)
	}
}

func TestUsers_GetMe_InvalidToken(t *testing.T) {
	statusCode, _ := do[domain.UserModel](t, "GET", "/users/me", nil, "invalid.token.here")

	if statusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", statusCode)
	}
}
