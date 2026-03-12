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

	boardColumn := createBoardColumn(t, boardID, tokens.AccessToken, randomBoardColumnName(), 0)
	boardColumnID := uuidToString(boardColumn.ID)

	ticket := createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "story", "medium")
	ticketID := uuidToString(ticket.ID)

	// Move ticket to board column
	statusCode, resp := do[domain.TicketModel](t, "PATCH", "/tickets/"+ticketID+"/move-board-column", domain.TicketBoardMoveModel{
		BoardID:       stringToUUID(boardID),
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
}

func TestTicket_MoveToBoardColumn_MismatchedColumn(t *testing.T) {
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

	// Create sprint, boards, board columns, and ticket
	sprint := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprintID := uuidToString(sprint.ID)

	board1 := createBoard(t, sprintID, tokens.AccessToken, randomBoardName())
	boardID1 := uuidToString(board1.ID)

	board2 := createBoard(t, sprintID, tokens.AccessToken, randomBoardName())
	boardID2 := uuidToString(board2.ID)

	// Create column in board2
	boardColumn2 := createBoardColumn(t, boardID2, tokens.AccessToken, randomBoardColumnName(), 0)
	boardColumnID2 := uuidToString(boardColumn2.ID)

	ticket := createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "story", "medium")
	ticketID := uuidToString(ticket.ID)

	// Try to move ticket to board1 with column from board2
	statusCode, resp := do[domain.TicketModel](t, "PATCH", "/tickets/"+ticketID+"/move-board-column", domain.TicketBoardMoveModel{
		BoardID:       stringToUUID(boardID1),
		BoardColumnID: stringToUUID(boardColumnID2),
	}, tokens.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400 for mismatched column, got %d: %v", statusCode, resp.Error)
	}
}
