package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestProject_GetByID_Success(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	// Create org and project
	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := uuidToString(orgResp.Data.ID)
	projectName := "Test Project " + randomString(8)

	projResp := createProject(t, orgID, tokens.AccessToken, randomProjectKey(), projectName, "public")
	projectID := uuidToString(projResp.ID)

	// Get the project
	statusCode, resp := do[domain.ProjectModel](t, "GET", "/projects/"+projectID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil || uuidToString(resp.Data.ID) != projectID {
		t.Fatalf("expected project ID %s, got %v", projectID, resp.Data)
	}

	if resp.Data.Name != projectName {
		t.Fatalf("expected name '%s', got '%s'", projectName, resp.Data.Name)
	}
}

func TestProject_GetByID_NotFound(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	nonExistentID := "550e8400-e29b-41d4-a716-446655440000"

	statusCode, _ := do[domain.ProjectModel](t, "GET", "/projects/"+nonExistentID, nil, tokens.AccessToken)

	if statusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", statusCode)
	}
}

func TestProject_GetByID_InvalidUUID(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, _ := do[domain.ProjectModel](t, "GET", "/projects/not-a-uuid", nil, tokens.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}
}

func TestProject_GetByID_ResponseStructure(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	// Create org and project
	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := uuidToString(orgResp.Data.ID)
	projectKey := randomProjectKey()
	projectName := "Test Project " + randomString(8)

	projResp := createProject(t, orgID, tokens.AccessToken, projectKey, projectName, "private")
	projectID := uuidToString(projResp.ID)

	// Get the project and verify response structure
	statusCode, resp := do[domain.ProjectModel](t, "GET", "/projects/"+projectID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if resp.Data == nil {
		t.Fatal("expected project data in response")
	}

	// Verify all expected fields are present
	if uuidToString(resp.Data.ID) == "" {
		t.Fatal("expected non-empty ID in response")
	}

	if resp.Data.Key == "" {
		t.Fatal("expected Key field in response")
	}

	if resp.Data.Key != projectKey {
		t.Fatalf("expected key '%s', got '%s'", projectKey, resp.Data.Key)
	}

	if resp.Data.Name == "" {
		t.Fatal("expected Name field in response")
	}

	if resp.Data.Name != projectName {
		t.Fatalf("expected name '%s', got '%s'", projectName, resp.Data.Name)
	}

	if resp.Data.Visibility == "" {
		t.Fatal("expected Visibility field in response")
	}

	if resp.Data.Visibility != "private" {
		t.Fatalf("expected visibility 'private', got '%s'", resp.Data.Visibility)
	}
}
