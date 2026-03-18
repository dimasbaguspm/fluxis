package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestTicket_MoveToSprint_WithBoard_Success(t *testing.T) {
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

	// Create sprint and board
	sprint := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprintID := uuidToString(sprint.ID)

	board := createBoard(t, sprintID, tokens.AccessToken, randomBoardName())
	boardID := uuidToString(board.ID)

	// Create ticket
	ticket := createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "story", "medium")
	ticketID := uuidToString(ticket.ID)

	// Move ticket to sprint with body
	statusCode, resp := do[domain.TicketModel](t, "PATCH", "/tickets/"+ticketID+"/move-to-sprint", domain.TicketMoveToSprintModel{
		SprintID: stringToUUID(sprintID),
	}, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected ticket data")
	}

	// Verify all three fields are populated
	if uuidToString(resp.Data.SprintID) != sprintID {
		t.Fatalf("expected sprint ID '%s', got '%s'", sprintID, uuidToString(resp.Data.SprintID))
	}

	if uuidToString(resp.Data.BoardID) != boardID {
		t.Fatalf("expected board ID '%s', got '%s'", boardID, uuidToString(resp.Data.BoardID))
	}

	if !resp.Data.BoardColumnID.Valid {
		t.Fatal("expected board column ID to be valid")
	}
}

func TestTicket_MoveToSprint_ToBacklog_Success(t *testing.T) {
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

	// Create sprint, board, and ticket
	sprint := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprintID := uuidToString(sprint.ID)

	_ = createBoard(t, sprintID, tokens.AccessToken, randomBoardName())

	ticket := createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "story", "medium")
	ticketID := uuidToString(ticket.ID)

	// Move to sprint first
	statusCode, _ = do[domain.TicketModel](t, "PATCH", "/tickets/"+ticketID+"/move-to-sprint", domain.TicketMoveToSprintModel{
		SprintID: stringToUUID(sprintID),
	}, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("failed to move ticket to sprint")
	}

	// Move to backlog (empty body)
	statusCode, resp := do[domain.TicketModel](t, "PATCH", "/tickets/"+ticketID+"/move-to-sprint", domain.TicketMoveToSprintModel{}, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected ticket data")
	}

	// Verify all three fields are cleared
	if resp.Data.SprintID.Valid {
		t.Fatal("expected sprint ID to be nil")
	}

	if resp.Data.BoardID.Valid {
		t.Fatal("expected board ID to be nil")
	}

	if resp.Data.BoardColumnID.Valid {
		t.Fatal("expected board column ID to be nil")
	}
}

func TestTicket_MoveToSprint_SprintNotFound(t *testing.T) {
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

	ticket := createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "story", "medium")
	ticketID := uuidToString(ticket.ID)

	// Try to move to non-existent sprint
	invalidSprintID := "550e8400-e29b-41d4-a716-446655440000"
	statusCode, resp := do[domain.TicketModel](t, "PATCH", "/tickets/"+ticketID+"/move-to-sprint", domain.TicketMoveToSprintModel{
		SprintID: stringToUUID(invalidSprintID),
	}, tokens.AccessToken)

	if statusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d: %v", statusCode, resp.Error)
	}
}
