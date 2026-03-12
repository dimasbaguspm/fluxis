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
