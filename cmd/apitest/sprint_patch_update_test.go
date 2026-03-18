package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestSprint_Update_Success(t *testing.T) {
	// Create org, project, and sprint
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := uuidToString(orgResp.Data.ID)
	project := createProject(t, orgID, tokens.AccessToken, randomProjectKey(), "Test Project", "private")
	projectID := uuidToString(project.ID)

	sprint := createSprint(t, projectID, tokens.AccessToken, "Original Sprint Name")
	sprintID := uuidToString(sprint.ID)

	// Update the sprint
	updatedName := "Updated Sprint Name " + randomString(4)
	updatedGoal := "New sprint goal"
	statusCode, resp := do[domain.SprintModel](t, "PATCH", "/sprints/"+sprintID, domain.SprintUpdateModel{
		Name: updatedName,
		Goal: updatedGoal,
	}, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil || resp.Data.Name != updatedName {
		t.Fatalf("expected name '%s', got %v", updatedName, resp.Data)
	}

	if resp.Data.Goal != updatedGoal {
		t.Fatalf("expected goal '%s', got '%s'", updatedGoal, resp.Data.Goal)
	}
}

func TestSprint_Update_PartialFields(t *testing.T) {
	// Create org, project, and sprint
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := uuidToString(orgResp.Data.ID)
	project := createProject(t, orgID, tokens.AccessToken, randomProjectKey(), "Test Project", "private")
	projectID := uuidToString(project.ID)

	originalName := randomSprintName()
	sprint := createSprint(t, projectID, tokens.AccessToken, originalName)
	sprintID := uuidToString(sprint.ID)

	// Update only the goal (name should remain unchanged)
	updatedGoal := "Only goal changed"
	statusCode, resp := do[domain.SprintModel](t, "PATCH", "/sprints/"+sprintID, domain.SprintUpdateModel{
		Goal: updatedGoal,
	}, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if resp.Data.Name != originalName {
		t.Fatalf("expected name to remain '%s', got '%s'", originalName, resp.Data.Name)
	}

	if resp.Data.Goal != updatedGoal {
		t.Fatalf("expected goal '%s', got '%s'", updatedGoal, resp.Data.Goal)
	}
}

func TestSprint_Update_EmptyName(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := uuidToString(orgResp.Data.ID)
	project := createProject(t, orgID, tokens.AccessToken, randomProjectKey(), "Test Project", "private")
	projectID := uuidToString(project.ID)

	originalName := randomSprintName()
	sprint := createSprint(t, projectID, tokens.AccessToken, originalName)
	sprintID := uuidToString(sprint.ID)

	// Updating with empty name should not change it (partial updates)
	status, resp := do[domain.SprintModel](t, "PATCH", "/sprints/"+sprintID, domain.SprintUpdateModel{
		Name: "",
	}, tokens.AccessToken)

	if status != http.StatusOK {
		t.Fatalf("expected status 200, got %d", status)
	}

	// Name should remain unchanged
	if resp.Data.Name != originalName {
		t.Fatalf("expected name to remain '%s', got '%s'", originalName, resp.Data.Name)
	}
}

func TestSprint_Update_NotFound(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	nonExistentID := "550e8400-e29b-41d4-a716-446655440000"

	statusCode, _ := do[domain.SprintModel](t, "PATCH", "/sprints/"+nonExistentID, domain.SprintUpdateModel{
		Name: "Updated Name",
	}, tokens.AccessToken)

	if statusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", statusCode)
	}
}

func TestSprint_Update_InvalidUUID(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, _ := do[domain.SprintModel](t, "PATCH", "/sprints/not-a-uuid", domain.SprintUpdateModel{
		Name: "Updated Name",
	}, tokens.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}
}

func TestSprint_Update_Unauthenticated(t *testing.T) {
	sprintID := "550e8400-e29b-41d4-a716-446655440000"

	statusCode, _ := do[domain.SprintModel](t, "PATCH", "/sprints/"+sprintID, domain.SprintUpdateModel{
		Name: "Updated Name",
	}, "")

	if statusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", statusCode)
	}
}

// TestSprint_Update_StatusNotAllowed verifies that status field in PATCH request is not accepted
// (status can only be changed via /start and /complete endpoints)
func TestSprint_Update_StatusNotAllowed(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := uuidToString(orgResp.Data.ID)
	project := createProject(t, orgID, tokens.AccessToken, randomProjectKey(), "Test Project", "private")
	projectID := uuidToString(project.ID)

	// Create sprint (status should be "planned")
	sprint := createSprint(t, projectID, tokens.AccessToken, "Test Sprint")
	sprintID := uuidToString(sprint.ID)

	if sprint.Status != "planned" {
		t.Fatalf("expected initial status 'planned', got '%s'", sprint.Status)
	}

	// Try to update status via PATCH (should be rejected as unknown field)
	updatedName := "Updated Name"
	statusCode, resp := do[domain.SprintModel](t, "PATCH", "/sprints/"+sprintID, map[string]interface{}{
		"name":   updatedName,
		"status": "active", // This should be rejected
	}, tokens.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400 for unknown field 'status', got %d: %v", statusCode, resp.Error)
	}

	// Verify error message mentions unknown field
	if resp.Error == nil || resp.Error.Message == "" {
		t.Fatalf("expected error response")
	}
}

// TestSprint_Update_StatusCannotBeChanged verifies that even valid fields update correctly
// while status field remains absent from the model
func TestSprint_Update_StatusCannotBeChanged(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := uuidToString(orgResp.Data.ID)
	project := createProject(t, orgID, tokens.AccessToken, randomProjectKey(), "Test Project", "private")
	projectID := uuidToString(project.ID)

	// Create sprint and start it
	sprint := createSprint(t, projectID, tokens.AccessToken, "Test Sprint")
	sprintID := uuidToString(sprint.ID)

	_, _ = do[domain.SprintModel](t, "POST", "/sprints/"+sprintID+"/start", nil, tokens.AccessToken)

	// Update only name (not status)
	updatedName := "Updated Name After Start"
	statusCode, resp := do[domain.SprintModel](t, "PATCH", "/sprints/"+sprintID, domain.SprintUpdateModel{
		Name: updatedName,
	}, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	// Verify name was updated but status is still "active" (not changed by PATCH)
	if resp.Data.Name != updatedName {
		t.Fatalf("expected name '%s', got '%s'", updatedName, resp.Data.Name)
	}

	if resp.Data.Status != "active" {
		t.Fatalf("expected status to remain 'active', got '%s'", resp.Data.Status)
	}
}
