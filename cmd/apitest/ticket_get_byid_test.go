package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestTicket_Get_Success(t *testing.T) {
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
	ticket := createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "bug", "critical")
	ticketID := uuidToString(ticket.ID)

	// Get ticket
	statusCode, resp := do[domain.TicketModel](t, "GET", "/tickets/"+ticketID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected ticket data")
	}

	if uuidToString(resp.Data.ID) != ticketID {
		t.Fatalf("expected ID '%s', got '%s'", ticketID, uuidToString(resp.Data.ID))
	}

	if resp.Data.Type != "bug" {
		t.Fatalf("expected type 'bug', got '%s'", resp.Data.Type)
	}
}

func TestTicket_GetByID_ResponseStructure(t *testing.T) {
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

	// Create ticket with all fields
	title := randomTicketTitle()
	description := "Test ticket description"
	ticketType := "story"
	priority := "high"

	statusCode, createResp := do[domain.TicketModel](t, "POST", "/tickets?projectId="+projectID, domain.TicketCreateModel{
		Title:       title,
		Description: description,
		Type:        ticketType,
		Priority:    priority,
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || createResp.Data == nil {
		t.Fatalf("failed to create ticket")
	}

	ticketID := uuidToString(createResp.Data.ID)

	// Get ticket and verify response structure
	statusCode, resp := do[domain.TicketModel](t, "GET", "/tickets/"+ticketID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if resp.Data == nil {
		t.Fatal("expected ticket data in response")
	}

	// Verify all expected fields are present
	if uuidToString(resp.Data.ID) == "" {
		t.Fatal("expected non-empty ID in response")
	}

	if resp.Data.Title == "" {
		t.Fatal("expected Title field in response")
	}

	if resp.Data.Title != title {
		t.Fatalf("expected title '%s', got '%s'", title, resp.Data.Title)
	}

	if resp.Data.Type == "" {
		t.Fatal("expected Type field in response")
	}

	if resp.Data.Type != ticketType {
		t.Fatalf("expected type '%s', got '%s'", ticketType, resp.Data.Type)
	}

	if resp.Data.Priority == "" {
		t.Fatal("expected Priority field in response")
	}

	if resp.Data.Priority != priority {
		t.Fatalf("expected priority '%s', got '%s'", priority, resp.Data.Priority)
	}

	// Key should be auto-generated
	if resp.Data.Key == "" {
		t.Fatal("expected Key field in response")
	}
}
