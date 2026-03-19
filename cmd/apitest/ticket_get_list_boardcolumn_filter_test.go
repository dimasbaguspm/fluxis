package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestTicket_List_FilterByBoardColumn(t *testing.T) {
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

	// Create sprint and board
	sprint := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprintID := uuidToString(sprint.ID)
	board := createBoard(t, sprintID, tokens.AccessToken, randomBoardName())
	boardID := uuidToString(board.ID)

	// Create board columns
	col1 := createBoardColumn(t, boardID, tokens.AccessToken, randomBoardColumnName())
	col1ID := uuidToString(col1.ID)

	// Create tickets and move to board column
	t1 := createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "story", "medium")
	t2 := createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "task", "low")

	// Move first ticket to sprint then to column
	do[domain.TicketModel](t, "PATCH", "/tickets/"+uuidToString(t1.ID)+"/move-to-sprint", domain.TicketMoveToSprintModel{
		SprintID: stringToUUID(sprintID),
	}, tokens.AccessToken)
	do[domain.TicketModel](t, "PATCH", "/tickets/"+uuidToString(t1.ID)+"/move-board-column", domain.TicketMoveBoardColumnModel{
		BoardColumnID: stringToUUID(col1ID),
	}, tokens.AccessToken)

	// Move second ticket to sprint but keep in different column (will get default column)
	do[domain.TicketModel](t, "PATCH", "/tickets/"+uuidToString(t2.ID)+"/move-to-sprint", domain.TicketMoveToSprintModel{
		SprintID: stringToUUID(sprintID),
	}, tokens.AccessToken)

	// List tickets by board column
	statusCode, resp := do[domain.TicketsPagedModel](t, "GET", "/tickets?projectId="+projectID+"&boardColumnId="+col1ID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if resp.Data == nil {
		t.Fatal("expected tickets data")
	}

	if len(resp.Data.Items) != 1 {
		t.Fatalf("expected 1 ticket in column, got %d", len(resp.Data.Items))
	}

	if resp.Data.TotalCount != 1 {
		t.Fatalf("expected total count 1, got %d", resp.Data.TotalCount)
	}

	// Verify the correct ticket is returned
	if uuidToString(resp.Data.Items[0].ID) != uuidToString(t1.ID) {
		t.Fatal("expected ticket 1 in column")
	}
}

func TestTicket_List_FilterByMultipleBoardColumns(t *testing.T) {
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

	// Create sprint and board
	sprint := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprintID := uuidToString(sprint.ID)
	board := createBoard(t, sprintID, tokens.AccessToken, randomBoardName())
	boardID := uuidToString(board.ID)

	// Create multiple board columns
	col1 := createBoardColumn(t, boardID, tokens.AccessToken, randomBoardColumnName())
	col2 := createBoardColumn(t, boardID, tokens.AccessToken, randomBoardColumnName())
	col1ID := uuidToString(col1.ID)
	col2ID := uuidToString(col2.ID)

	// Create tickets
	t1 := createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "story", "medium")
	t2 := createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "task", "low")
	t3 := createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "bug", "high")

	// Move tickets to different columns
	for _, ticket := range []domain.TicketModel{t1, t2, t3} {
		do[domain.TicketModel](t, "PATCH", "/tickets/"+uuidToString(ticket.ID)+"/move-to-sprint", domain.TicketMoveToSprintModel{
			SprintID: stringToUUID(sprintID),
		}, tokens.AccessToken)
	}

	do[domain.TicketModel](t, "PATCH", "/tickets/"+uuidToString(t1.ID)+"/move-board-column", domain.TicketMoveBoardColumnModel{
		BoardColumnID: stringToUUID(col1ID),
	}, tokens.AccessToken)

	do[domain.TicketModel](t, "PATCH", "/tickets/"+uuidToString(t2.ID)+"/move-board-column", domain.TicketMoveBoardColumnModel{
		BoardColumnID: stringToUUID(col2ID),
	}, tokens.AccessToken)

	// Query for multiple columns
	statusCode, resp := do[domain.TicketsPagedModel](t, "GET", "/tickets?projectId="+projectID+"&boardColumnId="+col1ID+"&boardColumnId="+col2ID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if len(resp.Data.Items) != 2 {
		t.Fatalf("expected 2 tickets in both columns, got %d", len(resp.Data.Items))
	}

	if resp.Data.TotalCount != 2 {
		t.Fatalf("expected total count 2, got %d", resp.Data.TotalCount)
	}
}

func TestTicket_List_FilterByBoardColumnWithSprintFilter(t *testing.T) {
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

	// Create 2 sprints and boards
	sprint1 := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprint2 := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprint1ID := uuidToString(sprint1.ID)
	sprint2ID := uuidToString(sprint2.ID)

	board1 := createBoard(t, sprint1ID, tokens.AccessToken, randomBoardName())
	_ = createBoard(t, sprint2ID, tokens.AccessToken, randomBoardName())
	board1ID := uuidToString(board1.ID)

	// Create columns only in board1
	col1 := createBoardColumn(t, board1ID, tokens.AccessToken, randomBoardColumnName())
	col1ID := uuidToString(col1.ID)

	// Create tickets
	t1 := createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "story", "medium")
	t2 := createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "task", "low")

	// Move t1 to sprint1+col1
	do[domain.TicketModel](t, "PATCH", "/tickets/"+uuidToString(t1.ID)+"/move-to-sprint", domain.TicketMoveToSprintModel{
		SprintID: stringToUUID(sprint1ID),
	}, tokens.AccessToken)
	do[domain.TicketModel](t, "PATCH", "/tickets/"+uuidToString(t1.ID)+"/move-board-column", domain.TicketMoveBoardColumnModel{
		BoardColumnID: stringToUUID(col1ID),
	}, tokens.AccessToken)

	// Move t2 to sprint2
	do[domain.TicketModel](t, "PATCH", "/tickets/"+uuidToString(t2.ID)+"/move-to-sprint", domain.TicketMoveToSprintModel{
		SprintID: stringToUUID(sprint2ID),
	}, tokens.AccessToken)

	// Query for sprint1 AND boardColumnId
	statusCode, resp := do[domain.TicketsPagedModel](t, "GET", "/tickets?projectId="+projectID+"&sprintId="+sprint1ID+"&boardColumnId="+col1ID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if len(resp.Data.Items) != 1 {
		t.Fatalf("expected 1 ticket, got %d", len(resp.Data.Items))
	}

	if uuidToString(resp.Data.Items[0].ID) != uuidToString(t1.ID) {
		t.Fatal("expected ticket 1")
	}
}

func TestTicket_List_FilterByBoardColumnEmptyResult(t *testing.T) {
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

	// Create sprint and board with column
	sprint := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprintID := uuidToString(sprint.ID)
	board := createBoard(t, sprintID, tokens.AccessToken, randomBoardName())
	boardID := uuidToString(board.ID)
	col1 := createBoardColumn(t, boardID, tokens.AccessToken, randomBoardColumnName())
	col1ID := uuidToString(col1.ID)

	// Create ticket but don't move to the column
	_ = createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "story", "medium")

	// Query for non-existent tickets in column
	statusCode, resp := do[domain.TicketsPagedModel](t, "GET", "/tickets?projectId="+projectID+"&boardColumnId="+col1ID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if len(resp.Data.Items) != 0 {
		t.Fatalf("expected 0 tickets, got %d", len(resp.Data.Items))
	}

	if resp.Data.TotalCount != 0 {
		t.Fatalf("expected total count 0, got %d", resp.Data.TotalCount)
	}
}

func TestTicket_List_FilterByBoardColumnWithPagination(t *testing.T) {
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

	// Create sprint and board
	sprint := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprintID := uuidToString(sprint.ID)
	board := createBoard(t, sprintID, tokens.AccessToken, randomBoardName())
	boardID := uuidToString(board.ID)

	// Create board column
	col1 := createBoardColumn(t, boardID, tokens.AccessToken, randomBoardColumnName())
	col1ID := uuidToString(col1.ID)

	// Create 5 tickets in the column
	for i := 0; i < 5; i++ {
		ticket := createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "story", "medium")
		do[domain.TicketModel](t, "PATCH", "/tickets/"+uuidToString(ticket.ID)+"/move-to-sprint", domain.TicketMoveToSprintModel{
			SprintID: stringToUUID(sprintID),
		}, tokens.AccessToken)
		do[domain.TicketModel](t, "PATCH", "/tickets/"+uuidToString(ticket.ID)+"/move-board-column", domain.TicketMoveBoardColumnModel{
			BoardColumnID: stringToUUID(col1ID),
		}, tokens.AccessToken)
	}

	// Get first page with size 2
	statusCode, resp := do[domain.TicketsPagedModel](t, "GET", "/tickets?projectId="+projectID+"&boardColumnId="+col1ID+"&pageSize=2&pageNumber=1", nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if len(resp.Data.Items) != 2 {
		t.Fatalf("expected 2 tickets on page 1, got %d", len(resp.Data.Items))
	}

	if resp.Data.TotalCount != 5 {
		t.Fatalf("expected total count 5, got %d", resp.Data.TotalCount)
	}

	if resp.Data.TotalPages != 3 {
		t.Fatalf("expected 3 total pages, got %d", resp.Data.TotalPages)
	}

	// Get second page
	statusCode, resp2 := do[domain.TicketsPagedModel](t, "GET", "/tickets?projectId="+projectID+"&boardColumnId="+col1ID+"&pageSize=2&pageNumber=2", nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if len(resp2.Data.Items) != 2 {
		t.Fatalf("expected 2 tickets on page 2, got %d", len(resp2.Data.Items))
	}
}

func TestTicket_List_FilterByBoardColumnInvalidId(t *testing.T) {
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

	// Create a ticket
	createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "story", "medium")

	// Query with invalid board column ID
	// Invalid UUIDs are treated as empty filters (consistent with other array parameters)
	invalidStatusCode, resp := do[domain.TicketsPagedModel](t, "GET", "/tickets?projectId="+projectID+"&boardColumnId=invalid-id", nil, tokens.AccessToken)

	if invalidStatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", invalidStatusCode)
	}

	// With invalid ID treated as empty filter, should return all tickets for the project
	if len(resp.Data.Items) != 1 {
		t.Fatalf("expected 1 ticket (invalid filter treated as empty), got %d", len(resp.Data.Items))
	}
}

func TestTicket_List_FilterByBoardColumnNonExistentId(t *testing.T) {
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

	// Non-existent board column ID
	nonExistentColumnID := "550e8400-e29b-41d4-a716-446655440000"

	// Query with non-existent board column ID
	nonExistentStatusCode, resp := do[domain.TicketsPagedModel](t, "GET", "/tickets?projectId="+projectID+"&boardColumnId="+nonExistentColumnID, nil, tokens.AccessToken)

	if nonExistentStatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", nonExistentStatusCode)
	}

	// Should return no tickets
	if len(resp.Data.Items) != 0 {
		t.Fatalf("expected 0 tickets for non-existent column, got %d", len(resp.Data.Items))
	}
}
