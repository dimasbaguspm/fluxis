package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestSprint_Complete_Success(t *testing.T) {
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

// TestSprint_Complete_AllowsPlanned verifies that planned sprints can also be completed
// (no need to go active first)
func TestSprint_Complete_AllowsPlanned(t *testing.T) {
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

	// Verify sprint is in planned status
	if sprint.Status != "planned" {
		t.Fatalf("expected initial status 'planned', got '%s'", sprint.Status)
	}

	// Complete the sprint directly from planned (without starting)
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

// TestSprint_Complete_ReleasesActiveSlot verifies that completing a sprint frees up the active slot
// allowing another sprint in the project to be started
func TestSprint_Complete_ReleasesActiveSlot(t *testing.T) {
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

	sprint1 := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprint1ID := uuidToString(sprint1.ID)

	sprint2 := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprint2ID := uuidToString(sprint2.ID)

	// Start and complete first sprint
	_, _ = do[domain.SprintModel](t, "POST", "/sprints/"+sprint1ID+"/start", nil, tokens.AccessToken)
	_, _ = do[domain.SprintModel](t, "POST", "/sprints/"+sprint1ID+"/completed", nil, tokens.AccessToken)

	// Now should be able to start second sprint
	statusCode, resp := do[domain.SprintModel](t, "POST", "/sprints/"+sprint2ID+"/start", nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200 after completing active sprint, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil || resp.Data.Status != "active" {
		t.Fatalf("expected sprint2 to be active after sprint1 completed")
	}
}
