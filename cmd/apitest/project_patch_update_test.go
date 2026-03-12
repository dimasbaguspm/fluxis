package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestProject_Update_Success(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	// Create org and project
	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := uuidToString(orgResp.Data.ID)
	projResp := createProject(t, orgID, tokens.AccessToken, randomProjectKey(), "Original Name", "private")
	projectID := uuidToString(projResp.ID)

	// Update the project
	updatedName := "Updated Name " + randomString(8)
	statusCode, resp := do[domain.ProjectModel](t, "PATCH", "/projects/"+projectID, domain.ProjectUpdateModel{
		Name:        updatedName,
		Description: "New description",
	}, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil || resp.Data.Name != updatedName {
		t.Fatalf("expected name '%s', got %v", updatedName, resp.Data)
	}

	if resp.Data.Description != "New description" {
		t.Fatalf("expected description 'New description', got '%s'", resp.Data.Description)
	}
}

func TestProject_Update_MissingName(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := uuidToString(orgResp.Data.ID)
	projResp := createProject(t, orgID, tokens.AccessToken, randomProjectKey(), "Test Project", "private")
	projectID := uuidToString(projResp.ID)

	status, _ := do[domain.ProjectModel](t, "PATCH", "/projects/"+projectID, domain.ProjectUpdateModel{
		Name: "",
	}, tokens.AccessToken)

	if status != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", status)
	}
}

func TestProject_Update_NotFound(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	nonExistentID := "550e8400-e29b-41d4-a716-446655440000"

	statusCode, _ := do[domain.ProjectModel](t, "PATCH", "/projects/"+nonExistentID, domain.ProjectUpdateModel{
		Name: "Updated Name",
	}, tokens.AccessToken)

	if statusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", statusCode)
	}
}

func TestProject_Update_InvalidUUID(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, _ := do[domain.ProjectModel](t, "PATCH", "/projects/not-a-uuid", domain.ProjectUpdateModel{
		Name: "Updated Name",
	}, tokens.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}
}

func TestProject_Update_Unauthenticated(t *testing.T) {
	projectID := "550e8400-e29b-41d4-a716-446655440000"

	statusCode, _ := do[domain.ProjectModel](t, "PATCH", "/projects/"+projectID, domain.ProjectUpdateModel{
		Name: "Updated Name",
	}, "")

	if statusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", statusCode)
	}
}
