package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestTickets_Create_Success(t *testing.T) {
	// Setup: Create user, org, project
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", map[string]string{
		"name": "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := uuidToString(orgResp.Data.ID)
	project := createProject(t, orgID, tokens.AccessToken, randomProjectKey(), "Test Project "+randomString(8), "private")
	projectID := uuidToString(project.ID)

	// Create ticket
	title := randomTicketTitle()
	statusCode, resp := do[domain.TicketModel](t, "POST", "/tickets?projectId="+projectID, map[string]string{
		"title":    title,
		"type":     "story",
		"priority": "high",
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected ticket data")
	}

	if resp.Data.Title != title {
		t.Fatalf("expected title '%s', got '%s'", title, resp.Data.Title)
	}

	if resp.Data.Type != "story" {
		t.Fatalf("expected type 'story', got '%s'", resp.Data.Type)
	}

	if resp.Data.Priority != "high" {
		t.Fatalf("expected priority 'high', got '%s'", resp.Data.Priority)
	}

	if resp.Data.Key == "" {
		t.Fatal("expected non-empty key")
	}

	// Key should be in format PROJECT_KEY-NUMBER
	projectKey := project.Key
	if resp.Data.Key[:len(projectKey)] != projectKey {
		t.Fatalf("expected key to start with '%s', got '%s'", projectKey, resp.Data.Key)
	}
}

func TestTickets_Get_Success(t *testing.T) {
	// Setup
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", map[string]string{
		"name": "Test Org " + randomString(8),
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

func TestTickets_List_Success(t *testing.T) {
	// Setup
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", map[string]string{
		"name": "Test Org " + randomString(8),
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

func TestTickets_Update_Success(t *testing.T) {
	// Setup
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", map[string]string{
		"name": "Test Org " + randomString(8),
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
	statusCode, resp := do[domain.TicketModel](t, "PATCH", "/tickets/"+ticketID, map[string]string{
		"title":    newTitle,
		"priority": "critical",
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

func TestTickets_MoveToSprint_Success(t *testing.T) {
	// Setup
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", map[string]string{
		"name": "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := uuidToString(orgResp.Data.ID)
	project := createProject(t, orgID, tokens.AccessToken, randomProjectKey(), "Test Project "+randomString(8), "private")
	projectID := uuidToString(project.ID)

	// Create sprint and ticket
	sprint := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprintID := uuidToString(sprint.ID)

	ticket := createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "story", "medium")
	ticketID := uuidToString(ticket.ID)

	// Move ticket to sprint
	statusCode, resp := do[domain.TicketModel](t, "PATCH", "/tickets/"+ticketID+"/move-to-sprint?sprintId="+sprintID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected ticket data")
	}

	if uuidToString(resp.Data.SprintID) != sprintID {
		t.Fatalf("expected sprint ID '%s', got '%s'", sprintID, uuidToString(resp.Data.SprintID))
	}
}

func TestTickets_Delete_Success(t *testing.T) {
	// Setup
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", map[string]string{
		"name": "Test Org " + randomString(8),
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

func TestTickets_Create_Unauthenticated(t *testing.T) {
	projectID := "550e8400-e29b-41d4-a716-446655440000"

	statusCode, _ := do[domain.TicketModel](t, "POST", "/tickets?projectId="+projectID, map[string]string{
		"title":    "Test Ticket",
		"type":     "story",
		"priority": "medium",
	}, "")

	if statusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", statusCode)
	}
}

func TestTickets_Create_InvalidProjectId(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, _ := do[domain.TicketModel](t, "POST", "/tickets?projectId=invalid-uuid", map[string]string{
		"title":    "Test Ticket",
		"type":     "story",
		"priority": "medium",
	}, tokens.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}
}

func TestTickets_Create_InvalidTicketType(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	status, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", map[string]string{
		"name": "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if status != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := uuidToString(orgResp.Data.ID)
	project := createProject(t, orgID, tokens.AccessToken, randomProjectKey(), "Test Project "+randomString(8), "private")
	projectID := uuidToString(project.ID)

	statusCode, _ := do[domain.TicketModel](t, "POST", "/tickets?projectId="+projectID, map[string]string{
		"title":    "Test Ticket",
		"type":     "invalid",
		"priority": "medium",
	}, tokens.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}
}
