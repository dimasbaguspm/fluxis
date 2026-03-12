package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestSprints_Create_Success(t *testing.T) {
	// Create org and project first
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org: status=%d", statusCode)
	}

	orgID := uuidToString(orgResp.Data.ID)
	project := createProject(t, orgID, tokens.AccessToken, randomProjectKey(), "Test Project", "private")
	projectID := uuidToString(project.ID)

	// Create sprint
	sprintName := randomSprintName()
	statusCode, resp := do[domain.SprintModel](t, "POST", "/sprints", domain.SprintCreateModel{
		Name:      sprintName,
		ProjectID: stringToUUID(projectID),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected sprint data")
	}

	if resp.Data.Name != sprintName {
		t.Fatalf("expected name '%s', got '%s'", sprintName, resp.Data.Name)
	}

	if resp.Data.Status != "planned" {
		t.Fatalf("expected status 'planned', got '%s'", resp.Data.Status)
	}

	if uuidToString(resp.Data.ID) == "" {
		t.Fatal("expected non-empty ID")
	}
}

func TestSprints_Create_WithGoal(t *testing.T) {
	// Create org and project
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

	// Create sprint with goal
	sprintName := randomSprintName()
	goal := "Complete user authentication feature"
	statusCode, resp := do[domain.SprintModel](t, "POST", "/sprints", domain.SprintCreateModel{
		Name:      sprintName,
		ProjectID: stringToUUID(projectID),
		Goal:      goal,
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || resp.Data == nil {
		t.Fatalf("expected status 201, got %d", statusCode)
	}

	if resp.Data.Goal != goal {
		t.Fatalf("expected goal '%s', got '%s'", goal, resp.Data.Goal)
	}
}

func TestSprints_Create_Unauthenticated(t *testing.T) {
	projectID := "550e8400-e29b-41d4-a716-446655440000"

	name := "Test Sprint"
	statusCode, _ := do[domain.SprintModel](t, "POST", "/sprints", domain.SprintCreateModel{
		Name:      name,
		ProjectID: stringToUUID(projectID),
	}, "")

	if statusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", statusCode)
	}
}

func TestSprints_Create_MissingProjectId(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	name := "Test Sprint"
	statusCode, resp := do[domain.SprintModel](t, "POST", "/sprints", domain.SprintCreateModel{
		Name: name,
	}, tokens.AccessToken)

	if statusCode != http.StatusNotFound {
		t.Fatalf("expected status 400, got %d: %v", statusCode, resp.Error)
	}
}

func TestSprints_List_ByProject(t *testing.T) {
	// Create org, project, and sprints
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

	// Create multiple sprints
	createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	createSprint(t, projectID, tokens.AccessToken, randomSprintName())

	// List sprints
	statusCode, resp := do[[]domain.SprintModel](t, "GET", "/sprints?projectId="+projectID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected sprint list data")
	}

	if len(*resp.Data) < 2 {
		t.Fatalf("expected at least 2 sprints, got %d", len(*resp.Data))
	}
}

func TestSprints_GetByID_Success(t *testing.T) {
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

	sprintName := randomSprintName()
	sprint := createSprint(t, projectID, tokens.AccessToken, sprintName)
	sprintID := uuidToString(sprint.ID)

	// Get the sprint
	statusCode, resp := do[domain.SprintModel](t, "GET", "/sprints/"+sprintID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil || uuidToString(resp.Data.ID) != sprintID {
		t.Fatalf("expected sprint ID %s, got %v", sprintID, resp.Data)
	}

	if resp.Data.Name != sprintName {
		t.Fatalf("expected name '%s', got '%s'", sprintName, resp.Data.Name)
	}
}

func TestSprints_GetByID_NotFound(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	nonExistentID := "550e8400-e29b-41d4-a716-446655440000"

	statusCode, _ := do[domain.SprintModel](t, "GET", "/sprints/"+nonExistentID, nil, tokens.AccessToken)

	if statusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", statusCode)
	}
}

func TestSprints_GetByID_InvalidUUID(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, _ := do[domain.SprintModel](t, "GET", "/sprints/not-a-uuid", nil, tokens.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}
}

func TestSprints_Update_Success(t *testing.T) {
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

func TestSprints_Start_Success(t *testing.T) {
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

	sprint := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprintID := uuidToString(sprint.ID)

	// Start the sprint
	statusCode, resp := do[domain.SprintModel](t, "POST", "/sprints/"+sprintID+"/start", nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil || resp.Data.Status != "active" {
		t.Fatalf("expected status 'active', got '%s'", resp.Data.Status)
	}

	if resp.Data.StartedAt == nil {
		t.Fatal("expected StartedAt to be set")
	}
}

func TestSprints_Complete_Success(t *testing.T) {
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

	sprint := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprintID := uuidToString(sprint.ID)

	// Start the sprint first
	do[domain.SprintModel](t, "POST", "/sprints/"+sprintID+"/start", nil, tokens.AccessToken)

	// Complete the sprint
	statusCode, resp := do[domain.SprintModel](t, "POST", "/sprints/"+sprintID+"/completed", nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil || resp.Data.Status != "completed" {
		t.Fatalf("expected status 'completed', got '%s'", resp.Data.Status)
	}

	if resp.Data.CompletedAt == nil {
		t.Fatal("expected CompletedAt to be set")
	}
}

func TestSprints_Update_PartialFields(t *testing.T) {
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
