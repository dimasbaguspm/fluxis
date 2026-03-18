package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestTicket_MoveToBoardColumn_Success(t *testing.T) {
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

	// Create sprint, board, board column, and ticket
	sprint := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprintID := uuidToString(sprint.ID)

	board := createBoard(t, sprintID, tokens.AccessToken, randomBoardName())
	boardID := uuidToString(board.ID)

	boardColumn := createBoardColumn(t, boardID, tokens.AccessToken, randomBoardColumnName())
	boardColumnID := uuidToString(boardColumn.ID)

	ticket := createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "story", "medium")
	ticketID := uuidToString(ticket.ID)

	// Move ticket to board column (only boardColumnId required)
	statusCode, resp := do[domain.TicketModel](t, "PATCH", "/tickets/"+ticketID+"/move-board-column", domain.TicketMoveBoardColumnModel{
		BoardColumnID: stringToUUID(boardColumnID),
	}, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected ticket data")
	}

	if uuidToString(resp.Data.BoardColumnID) != boardColumnID {
		t.Fatalf("expected board column ID '%s', got '%s'", boardColumnID, uuidToString(resp.Data.BoardColumnID))
	}

	if uuidToString(resp.Data.BoardID) != boardID {
		t.Fatalf("expected board ID '%s', got '%s'", boardID, uuidToString(resp.Data.BoardID))
	}
}

func TestTicket_MoveToBoardColumn_ColumnNotFound(t *testing.T) {
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

	// Try to move ticket to non-existent column
	invalidColumnID := "550e8400-e29b-41d4-a716-446655440000"
	statusCode, resp := do[domain.TicketModel](t, "PATCH", "/tickets/"+ticketID+"/move-board-column", domain.TicketMoveBoardColumnModel{
		BoardColumnID: stringToUUID(invalidColumnID),
	}, tokens.AccessToken)

	if statusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d: %v", statusCode, resp.Error)
	}
}
