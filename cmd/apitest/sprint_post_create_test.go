package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestSprint_Create_Success(t *testing.T) {
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

func TestSprint_Create_WithGoal(t *testing.T) {
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

func TestSprint_Create_Unauthenticated(t *testing.T) {
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

func TestSprint_Create_MissingProjectId(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	name := "Test Sprint"
	statusCode, resp := do[domain.SprintModel](t, "POST", "/sprints", domain.SprintCreateModel{
		Name: name,
	}, tokens.AccessToken)

	if statusCode != http.StatusNotFound {
		t.Fatalf("expected status 400, got %d: %v", statusCode, resp.Error)
	}
}

func TestSprint_Create_MissingName(t *testing.T) {
	// Create org and project first
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

	status, _ := do[domain.SprintModel](t, "POST", "/sprints", domain.SprintCreateModel{
		Name:      "",
		ProjectID: stringToUUID(projectID),
	}, tokens.AccessToken)

	if status != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", status)
	}
}

func TestSprint_Create_InvalidProjectIdFormat(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, _ := do[domain.SprintModel](t, "POST", "/sprints", domain.SprintCreateModel{
		Name:      "Test Sprint",
		ProjectID: stringToUUID("not-a-uuid"),
	}, tokens.AccessToken)

	// Invalid UUID format is treated as not found, not bad request
	if statusCode != http.StatusNotFound && statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 404 or 400, got %d", statusCode)
	}
}

func TestSprint_Create_NonExistentProject(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	nonExistentProjectID := "550e8400-e29b-41d4-a716-446655440000"

	statusCode, _ := do[domain.SprintModel](t, "POST", "/sprints", domain.SprintCreateModel{
		Name:      "Test Sprint",
		ProjectID: stringToUUID(nonExistentProjectID),
	}, tokens.AccessToken)

	if statusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", statusCode)
	}
}
