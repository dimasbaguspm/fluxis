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
	statusCode, resp := do[[]domain.ProjectModel](t, "GET", "/projects?orgId="+orgID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected project list data")
	}

	if len(*resp.Data) < 2 {
		t.Fatalf("expected at least 2 projects, got %d", len(*resp.Data))
	}
}
