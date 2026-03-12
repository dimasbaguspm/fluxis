package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestProject_Create_Success(t *testing.T) {
	// Create org first
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org: status=%d", statusCode)
	}

	orgID := uuidToString(orgResp.Data.ID)

	// Create project
	projectKey := randomProjectKey()
	projectName := "Test Project " + randomString(8)
	statusCode, resp := do[domain.ProjectModel](t, "POST", "/projects?orgId="+orgID, domain.ProjectCreateModel{
		Key:        projectKey,
		Name:       projectName,
		Visibility: "private",
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected project data")
	}

	if resp.Data.Key != projectKey {
		t.Fatalf("expected key '%s', got '%s'", projectKey, resp.Data.Key)
	}

	if resp.Data.Name != projectName {
		t.Fatalf("expected name '%s', got '%s'", projectName, resp.Data.Name)
	}

	if resp.Data.Visibility != "private" {
		t.Fatalf("expected visibility 'private', got '%s'", resp.Data.Visibility)
	}

	if uuidToString(resp.Data.ID) == "" {
		t.Fatal("expected non-empty ID")
	}
}

func TestProject_Create_Unauthenticated(t *testing.T) {
	orgID := "550e8400-e29b-41d4-a716-446655440000"

	statusCode, _ := do[domain.ProjectModel](t, "POST", "/projects?orgId="+orgID, domain.ProjectCreateModel{
		Key:        randomProjectKey(),
		Name:       "Test Project",
		Visibility: "private",
	}, "")

	if statusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", statusCode)
	}
}

func TestProject_Create_MissingOrgId(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, _ := do[domain.ProjectModel](t, "POST", "/projects", domain.ProjectCreateModel{
		Key:        randomProjectKey(),
		Name:       "Test Project",
		Visibility: "private",
	}, tokens.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}
}

func TestProject_Create_DuplicateKey(t *testing.T) {
	// Create org and first project
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := uuidToString(orgResp.Data.ID)
	projectKey := randomProjectKey()

	// Create first project
	createProject(t, orgID, tokens.AccessToken, projectKey, "Project 1", "private")

	// Try to create another project with the same key
	statusCode, _ = do[domain.ProjectModel](t, "POST", "/projects?orgId="+orgID, domain.ProjectCreateModel{
		Key:        projectKey,
		Name:       "Project 2",
		Visibility: "private",
	}, tokens.AccessToken)

	if statusCode != http.StatusConflict {
		t.Fatalf("expected status 409, got %d", statusCode)
	}
}

func TestProject_Create_WithDescription(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	// Create org
	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := uuidToString(orgResp.Data.ID)

	// Create project with description
	description := "This is a test project description"
	statusCode, resp := do[domain.ProjectModel](t, "POST", "/projects?orgId="+orgID, domain.ProjectCreateModel{
		Key:         randomProjectKey(),
		Name:        "Project with Desc",
		Description: description,
		Visibility:  "public",
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || resp.Data == nil {
		t.Fatalf("expected status 201, got %d", statusCode)
	}

	if resp.Data.Description != description {
		t.Fatalf("expected description '%s', got '%s'", description, resp.Data.Description)
	}
}
