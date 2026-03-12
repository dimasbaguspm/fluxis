package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestSprint_List_ByProject(t *testing.T) {
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
	statusCode, resp := do[domain.SprintsPagedModel](t, "GET", "/sprints?projectId="+projectID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected sprint list data")
	}

	if len(resp.Data.Items) < 2 {
		t.Fatalf("expected at least 2 sprints, got %d", len(resp.Data.Items))
	}

	if resp.Data.TotalCount < 2 {
		t.Fatalf("expected totalCount >= 2, got %d", resp.Data.TotalCount)
	}

	if resp.Data.PageNumber != 1 {
		t.Fatalf("expected page 1, got %d", resp.Data.PageNumber)
	}

	if resp.Data.PageSize != 25 {
		t.Fatalf("expected pageSize 25, got %d", resp.Data.PageSize)
	}
}

func TestSprint_List_WithPagination(t *testing.T) {
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

	// Create some sprints
	createSprint(t, projectID, tokens.AccessToken, "Sprint 1")
	createSprint(t, projectID, tokens.AccessToken, "Sprint 2")
	createSprint(t, projectID, tokens.AccessToken, "Sprint 3")

	// List sprints with custom pageSize
	statusCode, resp := do[domain.SprintsPagedModel](t, "GET", "/sprints?projectId="+projectID+"&pageSize=2&pageNumber=1", nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if resp.Data == nil {
		t.Fatal("expected sprint list data")
	}

	if len(resp.Data.Items) != 2 {
		t.Fatalf("expected 2 items per page, got %d", len(resp.Data.Items))
	}

	if resp.Data.PageSize != 2 {
		t.Fatalf("expected pageSize 2, got %d", resp.Data.PageSize)
	}

	if resp.Data.TotalPages < 2 {
		t.Fatalf("expected at least 2 pages, got %d", resp.Data.TotalPages)
	}
}

func TestSprint_List_WithNameFilter(t *testing.T) {
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

	// Create sprints with different names
	createSprint(t, projectID, tokens.AccessToken, "Alpha Sprint")
	createSprint(t, projectID, tokens.AccessToken, "Beta Sprint")
	createSprint(t, projectID, tokens.AccessToken, "Alpha Planning")

	// List sprints filtering by name
	statusCode, resp := do[domain.SprintsPagedModel](t, "GET", "/sprints?projectId="+projectID+"&name=Alpha", nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if resp.Data == nil {
		t.Fatal("expected sprint list data")
	}

	if len(resp.Data.Items) != 2 {
		t.Fatalf("expected 2 sprints matching 'Alpha', got %d", len(resp.Data.Items))
	}
}

func TestSprint_List_EmptyProject(t *testing.T) {
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

	// List sprints - should return empty list
	statusCode, resp := do[domain.SprintsPagedModel](t, "GET", "/sprints?projectId="+projectID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if resp.Data == nil {
		t.Fatal("expected sprint list data")
	}

	if len(resp.Data.Items) != 0 {
		t.Fatalf("expected 0 sprints, got %d", len(resp.Data.Items))
	}

	if resp.Data.TotalCount != 0 {
		t.Fatalf("expected totalCount 0, got %d", resp.Data.TotalCount)
	}
}

func TestSprint_List_FilterByID(t *testing.T) {
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

	// Create sprints
	s1 := createSprint(t, projectID, tokens.AccessToken, "Sprint 1")
	createSprint(t, projectID, tokens.AccessToken, "Sprint 2")

	// List sprints filtering by specific ID
	s1ID := uuidToString(s1.ID)
	statusCode, resp := do[domain.SprintsPagedModel](t, "GET", "/sprints?projectId="+projectID+"&id="+s1ID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if resp.Data == nil {
		t.Fatal("expected sprint list data")
	}

	if len(resp.Data.Items) != 1 {
		t.Fatalf("expected 1 sprint, got %d", len(resp.Data.Items))
	}

	if resp.Data.Items[0].ID != s1.ID {
		t.Fatalf("expected sprint %s, got %s", s1ID, uuidToString(resp.Data.Items[0].ID))
	}
}

func TestSprint_List_FilterByMultipleIDs(t *testing.T) {
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

	// Create sprints
	s1 := createSprint(t, projectID, tokens.AccessToken, "Sprint 1")
	s2 := createSprint(t, projectID, tokens.AccessToken, "Sprint 2")
	createSprint(t, projectID, tokens.AccessToken, "Sprint 3")

	// List sprints filtering by multiple IDs
	s1ID := uuidToString(s1.ID)
	s2ID := uuidToString(s2.ID)
	statusCode, resp := do[domain.SprintsPagedModel](t, "GET", "/sprints?projectId="+projectID+"&id="+s1ID+"&id="+s2ID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if resp.Data == nil {
		t.Fatal("expected sprint list data")
	}

	if len(resp.Data.Items) != 2 {
		t.Fatalf("expected 2 sprints, got %d", len(resp.Data.Items))
	}
}
