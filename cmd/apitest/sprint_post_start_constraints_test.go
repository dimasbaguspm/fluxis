package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestSprint_Start_ConflictMultipleActive(t *testing.T) {
	// Verify that only one sprint can be active per project
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

	// Create two sprints
	sprint1 := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprint1ID := uuidToString(sprint1.ID)

	sprint2 := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprint2ID := uuidToString(sprint2.ID)

	// Start first sprint - should succeed
	statusCode, resp1 := do[domain.SprintModel](t, "POST", "/sprints/"+sprint1ID+"/start", nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp1.Error)
	}

	if resp1.Data == nil || resp1.Data.Status != "active" {
		t.Fatalf("expected sprint1 to be active")
	}

	// Try to start second sprint - should return 409 Conflict
	statusCode, resp2 := do[domain.SprintModel](t, "POST", "/sprints/"+sprint2ID+"/start", nil, tokens.AccessToken)

	if statusCode != http.StatusConflict {
		t.Fatalf("expected status 409 Conflict, got %d", statusCode)
	}

	if resp2.Error == nil {
		t.Fatalf("expected error response")
	}

	if resp2.Error.Message != "project already has an active sprint" {
		t.Fatalf("expected error message 'project already has an active sprint', got '%s'", resp2.Error.Message)
	}
}

func TestSprint_Start_AllowAfterComplete(t *testing.T) {
	// Verify that after completing a sprint, another sprint can be activated
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

	// Create two sprints
	sprint1 := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprint1ID := uuidToString(sprint1.ID)

	sprint2 := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprint2ID := uuidToString(sprint2.ID)

	// Start first sprint
	code1, _ := do[domain.SprintModel](t, "POST", "/sprints/"+sprint1ID+"/start", nil, tokens.AccessToken)
	if code1 != http.StatusOK {
		t.Fatalf("failed to start sprint1")
	}

	// Complete first sprint
	code2, _ := do[domain.SprintModel](t, "POST", "/sprints/"+sprint1ID+"/completed", nil, tokens.AccessToken)
	if code2 != http.StatusOK {
		t.Fatalf("failed to complete sprint1")
	}

	// Now start second sprint - should succeed
	statusCode, resp := do[domain.SprintModel](t, "POST", "/sprints/"+sprint2ID+"/start", nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200 after completing first sprint, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil || resp.Data.Status != "active" {
		t.Fatalf("expected sprint2 to be active after sprint1 completed")
	}
}

func TestSprint_Start_DifferentProjectAllowsBothActive(t *testing.T) {
	// Verify that different projects can each have their own active sprint
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := uuidToString(orgResp.Data.ID)

	// Create two projects
	project1 := createProject(t, orgID, tokens.AccessToken, randomProjectKey(), "Project 1", "private")
	project1ID := uuidToString(project1.ID)

	project2 := createProject(t, orgID, tokens.AccessToken, randomProjectKey(), "Project 2", "private")
	project2ID := uuidToString(project2.ID)

	// Create sprints in each project
	sprint1 := createSprint(t, project1ID, tokens.AccessToken, randomSprintName())
	sprint1ID := uuidToString(sprint1.ID)

	sprint2 := createSprint(t, project2ID, tokens.AccessToken, randomSprintName())
	sprint2ID := uuidToString(sprint2.ID)

	// Start sprint in project 1
	statusCode3, _ := do[domain.SprintModel](t, "POST", "/sprints/"+sprint1ID+"/start", nil, tokens.AccessToken)
	if statusCode3 != http.StatusOK {
		t.Fatalf("failed to start sprint in project 1")
	}

	// Start sprint in project 2 - should succeed
	statusCode, resp := do[domain.SprintModel](t, "POST", "/sprints/"+sprint2ID+"/start", nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200 for different project, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil || resp.Data.Status != "active" {
		t.Fatalf("expected sprint in project 2 to be active")
	}
}

// TestSprint_Start_StatusTransition verifies status transitions through lifecycle
func TestSprint_Start_StatusTransition(t *testing.T) {
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

	// Create sprint and verify initial status
	sprint := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprintID := uuidToString(sprint.ID)

	if sprint.Status != "planned" {
		t.Fatalf("expected initial status 'planned', got '%s'", sprint.Status)
	}

	// Get sprint and verify status is still planned
	_, getResp := do[domain.SprintModel](t, "GET", "/sprints/"+sprintID, nil, tokens.AccessToken)
	if getResp.Data.Status != "planned" {
		t.Fatalf("expected GET response status 'planned', got '%s'", getResp.Data.Status)
	}

	// Start sprint
	_, startResp := do[domain.SprintModel](t, "POST", "/sprints/"+sprintID+"/start", nil, tokens.AccessToken)
	if startResp.Data.Status != "active" {
		t.Fatalf("expected status 'active' after start, got '%s'", startResp.Data.Status)
	}

	if startResp.Data.StartedAt == nil {
		t.Fatal("expected StartedAt to be set after start")
	}

	// Complete sprint
	_, completeResp := do[domain.SprintModel](t, "POST", "/sprints/"+sprintID+"/completed", nil, tokens.AccessToken)
	if completeResp.Data.Status != "completed" {
		t.Fatalf("expected status 'completed' after complete, got '%s'", completeResp.Data.Status)
	}

	if completeResp.Data.CompletedAt == nil {
		t.Fatal("expected CompletedAt to be set after complete")
	}
}
