package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestTicket_List_Success(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := uuidToString(orgResp.Data.ID)
	project := createProject(t, orgID, tokens.AccessToken, randomProjectKey(), "Test Project "+randomString(8), "private")
	projectID := uuidToString(project.ID)

	// Create multiple tickets
	createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "story", "medium")
	createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "task", "low")
	createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "epic", "high")

	// List tickets
	statusCode, resp := do[[]domain.TicketModel](t, "GET", "/tickets?projectId="+projectID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected tickets data")
	}

	if len(*resp.Data) != 3 {
		t.Fatalf("expected 3 tickets, got %d", len(*resp.Data))
	}
}
