package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestProject_UpdateVisibility_Success(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	// Create org and project
	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := uuidToString(orgResp.Data.ID)
	projResp := createProject(t, orgID, tokens.AccessToken, randomProjectKey(), "Test Project", "private")
	projectID := uuidToString(projResp.ID)

	// Update visibility to public
	statusCode, resp := do[domain.ProjectModel](t, "PATCH", "/projects/"+projectID+"/visibility", domain.ProjectVisibilityModel{
		Visibility: "public",
	}, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil || resp.Data.Visibility != "public" {
		t.Fatalf("expected visibility 'public', got '%s'", resp.Data.Visibility)
	}
}

func TestProject_UpdateVisibility_InvalidValue(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	// Create org and project
	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := uuidToString(orgResp.Data.ID)
	projResp := createProject(t, orgID, tokens.AccessToken, randomProjectKey(), "Test Project", "private")
	projectID := uuidToString(projResp.ID)

	// Try to update with invalid visibility
	statusCode, _ = do[domain.ProjectModel](t, "PATCH", "/projects/"+projectID+"/visibility", domain.ProjectVisibilityModel{
		Visibility: "protected",
	}, tokens.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}
}

func TestProject_UpdateVisibility_MissingVisibility(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	// Create org and project
	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := uuidToString(orgResp.Data.ID)
	projResp := createProject(t, orgID, tokens.AccessToken, randomProjectKey(), "Test Project", "private")
	projectID := uuidToString(projResp.ID)

	// Try to update with empty visibility
	status, _ := do[domain.ProjectModel](t, "PATCH", "/projects/"+projectID+"/visibility", domain.ProjectVisibilityModel{
		Visibility: "",
	}, tokens.AccessToken)

	if status != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", status)
	}
}

func TestProject_UpdateVisibility_NotFound(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	nonExistentID := "550e8400-e29b-41d4-a716-446655440000"

	statusCode, _ := do[domain.ProjectModel](t, "PATCH", "/projects/"+nonExistentID+"/visibility", domain.ProjectVisibilityModel{
		Visibility: "public",
	}, tokens.AccessToken)

	if statusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", statusCode)
	}
}

func TestProject_UpdateVisibility_InvalidUUID(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, _ := do[domain.ProjectModel](t, "PATCH", "/projects/not-a-uuid/visibility", domain.ProjectVisibilityModel{
		Visibility: "public",
	}, tokens.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}
}

func TestProject_UpdateVisibility_Unauthenticated(t *testing.T) {
	projectID := "550e8400-e29b-41d4-a716-446655440000"

	statusCode, _ := do[domain.ProjectModel](t, "PATCH", "/projects/"+projectID+"/visibility", domain.ProjectVisibilityModel{
		Visibility: "public",
	}, "")

	if statusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", statusCode)
	}
}
