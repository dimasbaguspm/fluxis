package apitest_test

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
)

// API Response envelope
type apiResponse[T any] struct {
	Data  *T        `json:"data"`
	Error *apiError `json:"error"`
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Auth response models
type authTokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// User response model
type userModel struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
}

// Org response model
type orgModel struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	TotalMembers int64  `json:"totalMembers"`
}

// Org member response model
type orgMemberModel struct {
	UserID string `json:"userId"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

// Generic request/response helper
func do[T any](tb testing.TB, method, path string, body interface{}, token string) (int, apiResponse[T]) {
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			tb.Fatalf("failed to marshal body: %v", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, testServer.URL+path, bodyReader)
	if err != nil {
		tb.Fatalf("failed to create request: %v", err)
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		tb.Fatalf("failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		tb.Fatalf("failed to read response body: %v", err)
	}

	var result apiResponse[T]
	if err := json.Unmarshal(respBody, &result); err != nil {
		tb.Logf("Response body: %s", string(respBody))
		tb.Fatalf("failed to unmarshal response: %v", err)
	}

	return resp.StatusCode, result
}

// Auth helpers
func register(tb testing.TB, email, displayName, password string) authTokens {
	statusCode, resp := do[authTokens](tb, "POST", "/auth/register", map[string]string{
		"email":       email,
		"displayName": displayName,
		"password":    password,
	}, "")

	if statusCode != http.StatusCreated {
		tb.Fatalf("register failed: got status %d, error: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		tb.Fatalf("register returned nil data")
	}

	return *resp.Data
}

func login(tb testing.TB, email, password string) authTokens {
	statusCode, resp := do[authTokens](tb, "POST", "/auth/login", map[string]string{
		"email":    email,
		"password": password,
	}, "")

	if statusCode != http.StatusOK {
		tb.Fatalf("login failed: got status %d, error: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		tb.Fatalf("login returned nil data")
	}

	return *resp.Data
}

// Random string generation
func randomString(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

func randomEmail() string {
	return fmt.Sprintf("user_%s@example.com", randomString(8))
}

// Helper to extract UUID from string (for comparisons)
func parseUUID(tb testing.TB, s string) pgtype.UUID {
	var uuid pgtype.UUID
	if err := uuid.Scan(s); err != nil {
		tb.Fatalf("failed to parse UUID %s: %v", s, err)
	}
	return uuid
}
