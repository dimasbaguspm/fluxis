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
