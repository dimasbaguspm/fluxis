package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestBoard_Reorder_Success(t *testing.T) {
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
	sprint := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprintID := uuidToString(sprint.ID)

	board1 := createBoard(t, sprintID, tokens.AccessToken, randomBoardName())
	board2 := createBoard(t, sprintID, tokens.AccessToken, randomBoardName())

	// Reorder boards: board2 at position 0, board1 at position 1
	statusCode, resp := do[[]domain.BoardModel](t, "PATCH", "/boards/reorder?sprintId="+sprintID, domain.BoardReorderModel{
		board2.ID,
		board1.ID,
	}, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil || len(*resp.Data) != 2 {
		t.Fatalf("expected 2 boards in response")
	}

	// Verify board2 is at position 0, board1 is at position 1
	data := *resp.Data
	if data[0].ID != board2.ID || data[0].Position != 0 {
		t.Fatalf("expected board2 at position 0, got %v at position %d", data[0].ID, data[0].Position)
	}
	if data[1].ID != board1.ID || data[1].Position != 1 {
		t.Fatalf("expected board1 at position 1, got %v at position %d", data[1].ID, data[1].Position)
	}
}

func TestBoard_Reorder_SingleBoard(t *testing.T) {
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
	sprint := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprintID := uuidToString(sprint.ID)

	board1 := createBoard(t, sprintID, tokens.AccessToken, randomBoardName())

	// Reorder single board (should succeed)
	statusCode, resp := do[[]domain.BoardModel](t, "PATCH", "/boards/reorder?sprintId="+sprintID, domain.BoardReorderModel{
		board1.ID,
	}, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil || len(*resp.Data) != 1 {
		t.Fatalf("expected 1 board in response")
	}

	if (*resp.Data)[0].Position != 0 {
		t.Fatalf("expected board at position 0, got %d", (*resp.Data)[0].Position)
	}
}

func TestBoard_Reorder_ThreeBoards(t *testing.T) {
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
	sprint := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprintID := uuidToString(sprint.ID)

	board1 := createBoard(t, sprintID, tokens.AccessToken, randomBoardName())
	board2 := createBoard(t, sprintID, tokens.AccessToken, randomBoardName())
	board3 := createBoard(t, sprintID, tokens.AccessToken, randomBoardName())

	// Reorder: board3, board1, board2
	statusCode, resp := do[[]domain.BoardModel](t, "PATCH", "/boards/reorder?sprintId="+sprintID, domain.BoardReorderModel{
		board3.ID,
		board1.ID,
		board2.ID,
	}, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	data := *resp.Data
	if len(data) != 3 {
		t.Fatalf("expected 3 boards in response, got %d", len(data))
	}

	if data[0].ID != board3.ID || data[0].Position != 0 {
		t.Fatalf("expected board3 at position 0, got %v at position %d", data[0].ID, data[0].Position)
	}
	if data[1].ID != board1.ID || data[1].Position != 1 {
		t.Fatalf("expected board1 at position 1, got %v at position %d", data[1].ID, data[1].Position)
	}
	if data[2].ID != board2.ID || data[2].Position != 2 {
		t.Fatalf("expected board2 at position 2, got %v at position %d", data[2].ID, data[2].Position)
	}
}

func TestBoard_Reorder_EmptyArray(t *testing.T) {
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
	sprint := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprintID := uuidToString(sprint.ID)

	createBoard(t, sprintID, tokens.AccessToken, randomBoardName())

	// Try to reorder with empty array
	statusCode, resp := do[[]domain.BoardModel](t, "PATCH", "/boards/reorder?sprintId="+sprintID, domain.BoardReorderModel{}, tokens.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}

	if resp.Error == nil {
		t.Fatalf("expected error message for empty array")
	}
}

func TestBoard_Reorder_IncompleteList(t *testing.T) {
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
	sprint := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprintID := uuidToString(sprint.ID)

	board1 := createBoard(t, sprintID, tokens.AccessToken, randomBoardName())
	_ = createBoard(t, sprintID, tokens.AccessToken, randomBoardName()) // Create second board to ensure validation checks count

	// Try to reorder with only one board (should fail - must include all boards)
	statusCode, resp := do[[]domain.BoardModel](t, "PATCH", "/boards/reorder?sprintId="+sprintID, domain.BoardReorderModel{
		board1.ID,
	}, tokens.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}

	if resp.Error == nil {
		t.Fatalf("expected error message for incomplete list")
	}
}

func TestBoard_Reorder_InvalidBoardId(t *testing.T) {
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
	sprint := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprintID := uuidToString(sprint.ID)

	board1 := createBoard(t, sprintID, tokens.AccessToken, randomBoardName())

	// Create a fake UUID for non-existent board
	fakeUUID := pgtype.UUID{}
	fakeUUID.Bytes = [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	fakeUUID.Valid = true

	// Try to reorder with non-existent board ID
	code, _ := do[[]domain.BoardModel](t, "PATCH", "/boards/reorder?sprintId="+sprintID, domain.BoardReorderModel{
		board1.ID,
		fakeUUID,
	}, tokens.AccessToken)

	if code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", code)
	}
}

func TestBoard_Reorder_BoardFromDifferentSprint(t *testing.T) {
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

	// Create two sprints
	sprint1 := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprint1ID := uuidToString(sprint1.ID)
	sprint2 := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprint2ID := uuidToString(sprint2.ID)

	// Create boards in each sprint
	board1 := createBoard(t, sprint1ID, tokens.AccessToken, randomBoardName())
	board2 := createBoard(t, sprint2ID, tokens.AccessToken, randomBoardName())

	// Try to reorder boards from sprint1 but include a board from sprint2
	statusCode, resp := do[[]domain.BoardModel](t, "PATCH", "/boards/reorder?sprintId="+sprint1ID, domain.BoardReorderModel{
		board1.ID,
		board2.ID,
	}, tokens.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}

	if resp.Error == nil {
		t.Fatalf("expected error message for board from different sprint")
	}
}

func TestBoard_Reorder_InvalidSprintId(t *testing.T) {
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
	sprint := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprintID := uuidToString(sprint.ID)

	board1 := createBoard(t, sprintID, tokens.AccessToken, randomBoardName())

	// Use a fake UUID string for non-existent sprint
	fakeSprintID := "00000000-0000-0000-0000-000000000000"

	// Try to reorder with non-existent sprint ID
	code, _ := do[[]domain.BoardModel](t, "PATCH", "/boards/reorder?sprintId="+fakeSprintID, domain.BoardReorderModel{
		board1.ID,
	}, tokens.AccessToken)

	if code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", code)
	}
}

func TestBoard_Reorder_MissingSprintId(t *testing.T) {
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
	sprint := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprintID := uuidToString(sprint.ID)

	board1 := createBoard(t, sprintID, tokens.AccessToken, randomBoardName())

	// Try to reorder without sprintId query parameter
	code, _ := do[[]domain.BoardModel](t, "PATCH", "/boards/reorder", domain.BoardReorderModel{
		board1.ID,
	}, tokens.AccessToken)

	if code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", code)
	}
}

func TestBoard_Reorder_Unauthenticated(t *testing.T) {
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
	sprint := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprintID := uuidToString(sprint.ID)

	board1 := createBoard(t, sprintID, tokens.AccessToken, randomBoardName())

	// Try to reorder without authentication
	code, _ := do[[]domain.BoardModel](t, "PATCH", "/boards/reorder?sprintId="+sprintID, domain.BoardReorderModel{
		board1.ID,
	}, "")

	if code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", code)
	}
}
