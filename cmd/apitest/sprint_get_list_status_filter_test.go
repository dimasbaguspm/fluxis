package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestSprint_List_FilterByStatus_Planned(t *testing.T) {
	// Create org, project, and sprints with different statuses
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
	sprint1 := createSprint(t, projectID, tokens.AccessToken, "Sprint Planned 1")
	sprint1ID := uuidToString(sprint1.ID)

	createSprint(t, projectID, tokens.AccessToken, "Sprint Planned 2")
	sprint3 := createSprint(t, projectID, tokens.AccessToken, "Sprint Active")
	sprint3ID := uuidToString(sprint3.ID)

	// Start sprint3 to make it active
	do[domain.SprintModel](t, "POST", "/sprints/"+sprint3ID+"/start", nil, tokens.AccessToken)

	// List only planned sprints
	statusCode, resp := do[domain.SprintsPagedModel](t, "GET", "/sprints?projectId="+projectID+"&status=planned", nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected sprint list data")
	}

	if len(resp.Data.Items) < 2 {
		t.Fatalf("expected at least 2 planned sprints, got %d", len(resp.Data.Items))
	}

	// Verify all returned sprints have "planned" status
	for _, sprint := range resp.Data.Items {
		if sprint.Status != "planned" {
			t.Fatalf("expected all sprints to have status 'planned', but got '%s' for sprint %s", sprint.Status, uuidToString(sprint.ID))
		}
	}

	// Verify sprint1 is in the results
	found := false
	for _, sprint := range resp.Data.Items {
		if uuidToString(sprint.ID) == sprint1ID {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected sprint1 to be in planned list")
	}
}

func TestSprint_List_FilterByStatus_Active(t *testing.T) {
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

	// Create sprints and set different statuses
	sprint1 := createSprint(t, projectID, tokens.AccessToken, "Sprint Active 1")
	sprint1ID := uuidToString(sprint1.ID)

	createSprint(t, projectID, tokens.AccessToken, "Sprint Planned")

	// Start sprint1 to make it active
	do[domain.SprintModel](t, "POST", "/sprints/"+sprint1ID+"/start", nil, tokens.AccessToken)

	// List only active sprints
	statusCode, resp := do[domain.SprintsPagedModel](t, "GET", "/sprints?projectId="+projectID+"&status=active", nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected sprint list data")
	}

	// Should have exactly 1 active sprint per project constraint
	if len(resp.Data.Items) != 1 {
		t.Fatalf("expected exactly 1 active sprint, got %d", len(resp.Data.Items))
	}

	if resp.Data.Items[0].Status != "active" {
		t.Fatalf("expected status 'active', got '%s'", resp.Data.Items[0].Status)
	}

	if uuidToString(resp.Data.Items[0].ID) != sprint1ID {
		t.Fatalf("expected active sprint to be sprint1")
	}
}

func TestSprint_List_FilterByStatus_Completed(t *testing.T) {
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

	// Create sprints with different statuses
	sprint1 := createSprint(t, projectID, tokens.AccessToken, "Sprint Completed 1")
	sprint1ID := uuidToString(sprint1.ID)

	sprint2 := createSprint(t, projectID, tokens.AccessToken, "Sprint Completed 2")
	sprint2ID := uuidToString(sprint2.ID)

	createSprint(t, projectID, tokens.AccessToken, "Sprint Planned")

	// Complete both sprints
	do[domain.SprintModel](t, "POST", "/sprints/"+sprint1ID+"/completed", nil, tokens.AccessToken)
	do[domain.SprintModel](t, "POST", "/sprints/"+sprint2ID+"/completed", nil, tokens.AccessToken)

	// List only completed sprints
	statusCode, resp := do[domain.SprintsPagedModel](t, "GET", "/sprints?projectId="+projectID+"&status=completed", nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected sprint list data")
	}

	if len(resp.Data.Items) < 2 {
		t.Fatalf("expected at least 2 completed sprints, got %d", len(resp.Data.Items))
	}

	// Verify all returned sprints have "completed" status
	for _, sprint := range resp.Data.Items {
		if sprint.Status != "completed" {
			t.Fatalf("expected all sprints to have status 'completed', but got '%s'", sprint.Status)
		}
	}
}

func TestSprint_List_FilterByStatus_InvalidStatus(t *testing.T) {
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

	// Try to list with invalid status
	statusCode, resp := do[domain.SprintsPagedModel](t, "GET", "/sprints?projectId="+projectID+"&status=invalid", nil, tokens.AccessToken)

	// Should return 400 Bad Request for invalid status enum value
	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400 for invalid status, got %d: %v", statusCode, resp.Error)
	}

	if resp.Error == nil {
		t.Fatalf("expected error response for invalid status")
	}
}
