package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestTicket_Update_Success(t *testing.T) {
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
	ticket := createTicket(t, projectID, tokens.AccessToken, "Original Title", "story", "medium")
	ticketID := uuidToString(ticket.ID)

	// Update ticket
	newTitle := "Updated Title " + randomString(8)
	statusCode, resp := do[domain.TicketModel](t, "PATCH", "/tickets/"+ticketID, domain.TicketUpdateModel{
		Title:    newTitle,
		Priority: "critical",
	}, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected ticket data")
	}

	if resp.Data.Title != newTitle {
		t.Fatalf("expected title '%s', got '%s'", newTitle, resp.Data.Title)
	}

	if resp.Data.Priority != "critical" {
		t.Fatalf("expected priority 'critical', got '%s'", resp.Data.Priority)
	}
}
