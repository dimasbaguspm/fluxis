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

	"github.com/dimasbaguspm/fluxis/pkg/domain"
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

func uuidToString(u pgtype.UUID) string {
	bytes, _ := u.MarshalJSON()
	return string(bytes[1 : len(bytes)-1])
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
	// Skip unmarshaling for empty bodies (e.g., 204 No Content)
	if len(respBody) > 0 {
		if err := json.Unmarshal(respBody, &result); err != nil {
			tb.Logf("Response body: %s", string(respBody))
			tb.Fatalf("failed to unmarshal response: %v", err)
		}
	}

	return resp.StatusCode, result
}

// Auth helpers
func register(tb testing.TB, email, displayName, password string) domain.AuthModel {
	statusCode, resp := do[domain.AuthModel](tb, "POST", "/auth/register", map[string]string{
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

// Project helpers
func createProject(tb testing.TB, orgID string, token string, key, name, visibility string) domain.ProjectModel {
	statusCode, resp := do[domain.ProjectModel](tb, "POST", "/projects?orgId="+orgID, map[string]string{
		"key":        key,
		"name":       name,
		"visibility": visibility,
	}, token)

	if statusCode != http.StatusCreated {
		tb.Fatalf("create project failed: got status %d, error: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		tb.Fatalf("create project returned nil data")
	}

	return *resp.Data
}

// Generate unique project key with max 10 chars
func randomProjectKey() string {
	return "p" + randomString(4)
}

// Sprint helpers
func createSprint(tb testing.TB, projectID string, token string, name string) domain.SprintModel {
	statusCode, resp := do[domain.SprintModel](tb, "POST", "/sprints?projectId="+projectID, map[string]string{
		"name": name,
	}, token)

	if statusCode != http.StatusCreated {
		tb.Fatalf("create sprint failed: got status %d, error: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		tb.Fatalf("create sprint returned nil data")
	}

	return *resp.Data
}

func randomSprintName() string {
	return "Sprint " + randomString(4)
}

// Board helpers
func createBoard(tb testing.TB, sprintID string, token string, name string) domain.BoardModel {
	statusCode, resp := do[domain.BoardModel](tb, "POST", "/boards?sprintId="+sprintID, map[string]string{
		"name": name,
	}, token)

	if statusCode != http.StatusCreated {
		tb.Fatalf("create board failed: got status %d, error: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		tb.Fatalf("create board returned nil data")
	}

	return *resp.Data
}

func randomBoardName() string {
	return "Board " + randomString(4)
}
