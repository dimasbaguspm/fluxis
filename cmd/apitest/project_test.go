package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestProjects_Create_Success(t *testing.T) {
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

func TestProjects_Create_Unauthenticated(t *testing.T) {
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

func TestProjects_Create_MissingOrgId(t *testing.T) {
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

func TestProjects_Create_DuplicateKey(t *testing.T) {
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

func TestProjects_List_ByOrg(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	// Create org
	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := uuidToString(orgResp.Data.ID)

	// Create some projects
	createProject(t, orgID, tokens.AccessToken, randomProjectKey(), "Project 1", "private")
	createProject(t, orgID, tokens.AccessToken, randomProjectKey(), "Project 2", "public")

	// List projects
	statusCode, resp := do[[]domain.ProjectModel](t, "GET", "/projects?orgId="+orgID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected project list data")
	}

	if len(*resp.Data) < 2 {
		t.Fatalf("expected at least 2 projects, got %d", len(*resp.Data))
	}
}

func TestProjects_GetByID_Success(t *testing.T) {
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

func TestProjects_GetByID_NotFound(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	nonExistentID := "550e8400-e29b-41d4-a716-446655440000"

	statusCode, _ := do[domain.ProjectModel](t, "GET", "/projects/"+nonExistentID, nil, tokens.AccessToken)

	if statusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", statusCode)
	}
}

func TestProjects_GetByID_InvalidUUID(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, _ := do[domain.ProjectModel](t, "GET", "/projects/not-a-uuid", nil, tokens.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}
}

func TestProjects_Update_Success(t *testing.T) {
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

func TestProjects_UpdateVisibility_Success(t *testing.T) {
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

func TestProjects_UpdateVisibility_InvalidValue(t *testing.T) {
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

func TestProjects_Delete_Success(t *testing.T) {
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

	// Delete the project
	statusCode, _ = do[struct{}](t, "DELETE", "/projects/"+projectID, nil, tokens.AccessToken)

	if statusCode != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", statusCode)
	}

	// Verify it's deleted
	statusCode, _ = do[domain.ProjectModel](t, "GET", "/projects/"+projectID, nil, tokens.AccessToken)

	if statusCode != http.StatusNotFound {
		t.Fatalf("expected status 404 after delete, got %d", statusCode)
	}
}

func TestProjects_Create_WithDescription(t *testing.T) {
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
