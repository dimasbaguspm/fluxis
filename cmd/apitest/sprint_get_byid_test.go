package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestSprint_GetByID_Success(t *testing.T) {
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

func TestSprint_GetByID_NotFound(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	nonExistentID := "550e8400-e29b-41d4-a716-446655440000"

	statusCode, _ := do[domain.SprintModel](t, "GET", "/sprints/"+nonExistentID, nil, tokens.AccessToken)

	if statusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", statusCode)
	}
}

func TestSprint_GetByID_InvalidUUID(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, _ := do[domain.SprintModel](t, "GET", "/sprints/not-a-uuid", nil, tokens.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}
}

func TestSprint_GetByID_ResponseStructure(t *testing.T) {
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

	// Create sprint with all fields
	sprintName := randomSprintName()
	goal := "Complete core features"

	statusCode, createResp := do[domain.SprintModel](t, "POST", "/sprints", domain.SprintCreateModel{
		Name:      sprintName,
		ProjectID: stringToUUID(projectID),
		Goal:      goal,
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || createResp.Data == nil {
		t.Fatalf("failed to create sprint")
	}

	sprintID := uuidToString(createResp.Data.ID)

	// Get the sprint and verify response structure
	statusCode, resp := do[domain.SprintModel](t, "GET", "/sprints/"+sprintID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if resp.Data == nil {
		t.Fatal("expected sprint data in response")
	}

	// Verify all expected fields are present
	if uuidToString(resp.Data.ID) == "" {
		t.Fatal("expected non-empty ID in response")
	}

	if resp.Data.Name == "" {
		t.Fatal("expected Name field in response")
	}

	if resp.Data.Name != sprintName {
		t.Fatalf("expected name '%s', got '%s'", sprintName, resp.Data.Name)
	}

	if resp.Data.Status == "" {
		t.Fatal("expected Status field in response")
	}

	if resp.Data.Status != "planned" {
		t.Fatalf("expected status 'planned', got '%s'", resp.Data.Status)
	}

	if resp.Data.Goal == "" {
		t.Logf("warning: Goal field empty in response (may be optional)")
	} else if resp.Data.Goal != goal {
		t.Fatalf("expected goal '%s', got '%s'", goal, resp.Data.Goal)
	}
}
