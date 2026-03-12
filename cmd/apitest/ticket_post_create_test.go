package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestTicket_Create_Success(t *testing.T) {
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
	title := randomTicketTitle()
	statusCode, resp := do[domain.TicketModel](t, "POST", "/tickets?projectId="+projectID, domain.TicketCreateModel{
		Title:    title,
		Type:     "story",
		Priority: "high",
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

func TestTicket_Create_Unauthenticated(t *testing.T) {
	projectID := "550e8400-e29b-41d4-a716-446655440000"

	statusCode, _ := do[domain.TicketModel](t, "POST", "/tickets?projectId="+projectID, domain.TicketCreateModel{
		Title:    "Test Ticket",
		Type:     "story",
		Priority: "medium",
	}, "")

	if statusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", statusCode)
	}
}

func TestTicket_Create_InvalidProjectId(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, _ := do[domain.TicketModel](t, "POST", "/tickets?projectId=invalid-uuid", domain.TicketCreateModel{
		Title:    "Test Ticket",
		Type:     "story",
		Priority: "medium",
	}, tokens.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}
}

func TestTicket_Create_InvalidTicketType(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	status, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if status != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := uuidToString(orgResp.Data.ID)
	project := createProject(t, orgID, tokens.AccessToken, randomProjectKey(), "Test Project "+randomString(8), "private")
	projectID := uuidToString(project.ID)

	statusCode, _ := do[domain.TicketModel](t, "POST", "/tickets?projectId="+projectID, domain.TicketCreateModel{
		Title:    "Test Ticket",
		Type:     "invalid",
		Priority: "medium",
	}, tokens.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}
}
