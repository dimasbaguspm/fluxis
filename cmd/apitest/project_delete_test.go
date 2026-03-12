package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestProject_Delete_Success(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	// Create org and project
	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := uuidToString(orgResp.Data.ID)
	projResp := createProject(t, orgID, tokens.AccessToken, randomProjectKey(), "Test Project", "private")
	projectID := uuidToString(projResp.ID)

	// Delete the project
	statusCode, _ = do[struct{}](t, "DELETE", "/projects/"+projectID, nil, tokens.AccessToken)

	if statusCode != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", statusCode)
	}

	// Verify it's deleted
	statusCode, _ = do[domain.ProjectModel](t, "GET", "/projects/"+projectID, nil, tokens.AccessToken)

	if statusCode != http.StatusNotFound {
		t.Fatalf("expected status 404 after delete, got %d", statusCode)
	}
}
