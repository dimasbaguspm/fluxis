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

func TestTicket_Update_AllowsPartialUpdate(t *testing.T) {
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

	originalPriority := "medium"
	ticket := createTicket(t, projectID, tokens.AccessToken, "Test Ticket", "story", originalPriority)
	ticketID := uuidToString(ticket.ID)

	// Update only priority, omitting title
	newPriority := "high"
	status, resp := do[domain.TicketModel](t, "PATCH", "/tickets/"+ticketID, domain.TicketUpdateModel{
		Priority: newPriority,
	}, tokens.AccessToken)

	if status != http.StatusOK {
		t.Fatalf("expected status 200, got %d", status)
	}

	// Priority should be updated
	if resp.Data.Priority != newPriority {
		t.Fatalf("expected priority '%s', got '%s'", newPriority, resp.Data.Priority)
	}
}

func TestTicket_Update_InvalidType(t *testing.T) {
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

	ticket := createTicket(t, projectID, tokens.AccessToken, "Test Ticket", "story", "medium")
	ticketID := uuidToString(ticket.ID)

	status, _ := do[domain.TicketModel](t, "PATCH", "/tickets/"+ticketID, domain.TicketUpdateModel{
		Type: "invalid",
	}, tokens.AccessToken)

	if status != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", status)
	}
}

func TestTicket_Update_InvalidPriority(t *testing.T) {
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

	ticket := createTicket(t, projectID, tokens.AccessToken, "Test Ticket", "story", "medium")
	ticketID := uuidToString(ticket.ID)

	status, _ := do[domain.TicketModel](t, "PATCH", "/tickets/"+ticketID, domain.TicketUpdateModel{
		Priority: "ultra",
	}, tokens.AccessToken)

	if status != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", status)
	}
}

func TestTicket_Update_NotFound(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	nonExistentID := "550e8400-e29b-41d4-a716-446655440000"

	statusCode, _ := do[domain.TicketModel](t, "PATCH", "/tickets/"+nonExistentID, domain.TicketUpdateModel{
		Title: "Updated Title",
	}, tokens.AccessToken)

	if statusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", statusCode)
	}
}

func TestTicket_Update_InvalidUUID(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, _ := do[domain.TicketModel](t, "PATCH", "/tickets/not-a-uuid", domain.TicketUpdateModel{
		Title: "Updated Title",
	}, tokens.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}
}

func TestTicket_Update_Unauthenticated(t *testing.T) {
	ticketID := "550e8400-e29b-41d4-a716-446655440000"

	statusCode, _ := do[domain.TicketModel](t, "PATCH", "/tickets/"+ticketID, domain.TicketUpdateModel{
		Title: "Updated Title",
	}, "")

	if statusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", statusCode)
	}
}
