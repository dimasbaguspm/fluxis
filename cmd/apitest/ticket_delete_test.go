package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestTicket_Delete_Success(t *testing.T) {
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

	// Create ticket
	ticket := createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "story", "medium")
	ticketID := uuidToString(ticket.ID)

	// Delete ticket
	status, _ := do[domain.TicketModel](t, "DELETE", "/tickets/"+ticketID, nil, tokens.AccessToken)

	if status != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", status)
	}

	// Verify ticket is deleted (should return 404)
	status, _ = do[domain.TicketModel](t, "GET", "/tickets/"+ticketID, nil, tokens.AccessToken)

	if status != http.StatusNotFound {
		t.Fatalf("expected status 404 after delete, got %d", status)
	}
}
