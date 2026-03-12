package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestProject_List_ByOrg(t *testing.T) {
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
	statusCode, resp := do[domain.ProjectsPagedModel](t, "GET", "/projects?orgId="+orgID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected project list data")
	}

	if len(resp.Data.Items) < 2 {
		t.Fatalf("expected at least 2 projects, got %d", len(resp.Data.Items))
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

func TestProject_List_WithPagination(t *testing.T) {
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
	createProject(t, orgID, tokens.AccessToken, randomProjectKey(), "Project 3", "private")

	// List projects with custom pageSize
	statusCode, resp := do[domain.ProjectsPagedModel](t, "GET", "/projects?orgId="+orgID+"&pageSize=2&pageNumber=1", nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if resp.Data == nil {
		t.Fatal("expected project list data")
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

func TestProject_List_WithNameFilter(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	// Create org
	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := uuidToString(orgResp.Data.ID)

	// Create projects with different names
	createProject(t, orgID, tokens.AccessToken, randomProjectKey(), "Alpha Project", "private")
	createProject(t, orgID, tokens.AccessToken, randomProjectKey(), "Beta Project", "public")
	createProject(t, orgID, tokens.AccessToken, randomProjectKey(), "Alpha Integration", "private")

	// List projects filtering by name
	statusCode, resp := do[domain.ProjectsPagedModel](t, "GET", "/projects?orgId="+orgID+"&name=Alpha", nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if resp.Data == nil {
		t.Fatal("expected project list data")
	}

	if len(resp.Data.Items) != 2 {
		t.Fatalf("expected 2 projects matching 'Alpha', got %d", len(resp.Data.Items))
	}
}

func TestProject_List_EmptyOrg(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	// Create org
	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := uuidToString(orgResp.Data.ID)

	// List projects - should return empty list
	statusCode, resp := do[domain.ProjectsPagedModel](t, "GET", "/projects?orgId="+orgID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if resp.Data == nil {
		t.Fatal("expected project list data")
	}

	if len(resp.Data.Items) != 0 {
		t.Fatalf("expected 0 projects, got %d", len(resp.Data.Items))
	}

	if resp.Data.TotalCount != 0 {
		t.Fatalf("expected totalCount 0, got %d", resp.Data.TotalCount)
	}
}

func TestProject_List_FilterByID(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	// Create org
	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := uuidToString(orgResp.Data.ID)

	// Create projects
	_, p1 := do[domain.ProjectModel](t, "POST", "/projects?orgId="+orgID, domain.ProjectCreateModel{
		Key:        randomProjectKey(),
		Name:       "Project 1",
		Visibility: "private",
	}, tokens.AccessToken)

	_, p2 := do[domain.ProjectModel](t, "POST", "/projects?orgId="+orgID, domain.ProjectCreateModel{
		Key:        randomProjectKey(),
		Name:       "Project 2",
		Visibility: "public",
	}, tokens.AccessToken)

	if p1.Data == nil || p2.Data == nil {
		t.Fatal("failed to create projects")
	}

	// List projects filtering by specific ID
	p1ID := uuidToString(p1.Data.ID)
	statusCode, resp := do[domain.ProjectsPagedModel](t, "GET", "/projects?orgId="+orgID+"&id="+p1ID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if resp.Data == nil {
		t.Fatal("expected project list data")
	}

	if len(resp.Data.Items) != 1 {
		t.Fatalf("expected 1 project, got %d", len(resp.Data.Items))
	}

	if resp.Data.Items[0].ID != p1.Data.ID {
		t.Fatalf("expected project %s, got %s", p1ID, uuidToString(resp.Data.Items[0].ID))
	}
}

func TestProject_List_FilterByMultipleIDs(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	// Create org
	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := uuidToString(orgResp.Data.ID)

	// Create projects
	_, p1 := do[domain.ProjectModel](t, "POST", "/projects?orgId="+orgID, domain.ProjectCreateModel{
		Key:        randomProjectKey(),
		Name:       "Project 1",
		Visibility: "private",
	}, tokens.AccessToken)

	_, p2 := do[domain.ProjectModel](t, "POST", "/projects?orgId="+orgID, domain.ProjectCreateModel{
		Key:        randomProjectKey(),
		Name:       "Project 2",
		Visibility: "public",
	}, tokens.AccessToken)

	_, p3 := do[domain.ProjectModel](t, "POST", "/projects?orgId="+orgID, domain.ProjectCreateModel{
		Key:        randomProjectKey(),
		Name:       "Project 3",
		Visibility: "private",
	}, tokens.AccessToken)

	if p1.Data == nil || p2.Data == nil || p3.Data == nil {
		t.Fatal("failed to create projects")
	}

	// List projects filtering by multiple IDs
	p1ID := uuidToString(p1.Data.ID)
	p2ID := uuidToString(p2.Data.ID)
	statusCode, resp := do[domain.ProjectsPagedModel](t, "GET", "/projects?orgId="+orgID+"&id="+p1ID+"&id="+p2ID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if resp.Data == nil {
		t.Fatal("expected project list data")
	}

	if len(resp.Data.Items) != 2 {
		t.Fatalf("expected 2 projects, got %d", len(resp.Data.Items))
	}
}
