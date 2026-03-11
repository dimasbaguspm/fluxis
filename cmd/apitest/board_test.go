package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestBoards_Create_Success(t *testing.T) {
	// Create org, project, and sprint first
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org: status=%d", statusCode)
	}

	orgID := uuidToString(orgResp.Data.ID)
	project := createProject(t, orgID, tokens.AccessToken, randomProjectKey(), "Test Project", "private")
	projectID := uuidToString(project.ID)

	sprint := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprintID := uuidToString(sprint.ID)

	// Create board
	boardName := randomBoardName()
	statusCode, resp := do[domain.BoardModel](t, "POST", "/boards?sprintId="+sprintID, domain.BoardCreateModel{
		Name: &boardName,
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected board data")
	}

	if resp.Data.Name != boardName {
		t.Fatalf("expected name '%s', got '%s'", boardName, resp.Data.Name)
	}

	if uuidToString(resp.Data.ID) == "" {
		t.Fatal("expected non-empty ID")
	}

	if resp.Data.Position != 0 {
		t.Fatalf("expected position 0 for first board, got %d", resp.Data.Position)
	}
}

func TestBoards_Create_Unauthenticated(t *testing.T) {
	sprintID := "550e8400-e29b-41d4-a716-446655440000"

	name := "Test Board"
	statusCode, _ := do[domain.BoardModel](t, "POST", "/boards?sprintId="+sprintID, domain.BoardCreateModel{
		Name: &name,
	}, "")

	if statusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", statusCode)
	}
}

func TestBoards_Create_MissingSprintId(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	name := "Test Board"
	statusCode, _ := do[domain.BoardModel](t, "POST", "/boards", domain.BoardCreateModel{
		Name: &name,
	}, tokens.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}
}

func TestBoards_Create_MissingName(t *testing.T) {
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

	code, _ := do[domain.BoardModel](t, "POST", "/boards?sprintId="+sprintID, map[string]interface{}{}, tokens.AccessToken)

	if code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", code)
	}
}

func TestBoards_List_BySprint(t *testing.T) {
	// Create org, project, sprint, and boards
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

	// Create multiple boards
	createBoard(t, sprintID, tokens.AccessToken, randomBoardName())
	createBoard(t, sprintID, tokens.AccessToken, randomBoardName())

	// List boards
	statusCode, resp := do[[]domain.BoardModel](t, "GET", "/boards?sprintId="+sprintID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected board list data")
	}

	if len(*resp.Data) < 2 {
		t.Fatalf("expected at least 2 boards, got %d", len(*resp.Data))
	}
}

func TestBoards_GetByID_Success(t *testing.T) {
	// Create org, project, sprint, and board
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

	boardName := randomBoardName()
	board := createBoard(t, sprintID, tokens.AccessToken, boardName)
	boardID := uuidToString(board.ID)

	// Get the board
	statusCode, resp := do[domain.BoardModel](t, "GET", "/boards/"+boardID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil || uuidToString(resp.Data.ID) != boardID {
		t.Fatalf("expected board ID %s, got %v", boardID, resp.Data)
	}

	if resp.Data.Name != boardName {
		t.Fatalf("expected name '%s', got '%s'", boardName, resp.Data.Name)
	}
}

func TestBoards_GetByID_NotFound(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	nonExistentID := "550e8400-e29b-41d4-a716-446655440000"

	statusCode, _ := do[domain.BoardModel](t, "GET", "/boards/"+nonExistentID, nil, tokens.AccessToken)

	if statusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", statusCode)
	}
}

func TestBoards_GetByID_InvalidUUID(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, _ := do[domain.BoardModel](t, "GET", "/boards/not-a-uuid", nil, tokens.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}
}

func TestBoards_Update_Success(t *testing.T) {
	// Create org, project, sprint, and board
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

	board := createBoard(t, sprintID, tokens.AccessToken, "Original Board Name")
	boardID := uuidToString(board.ID)

	// Update the board
	updatedName := "Updated Board Name " + randomString(4)
	statusCode, resp := do[domain.BoardModel](t, "PATCH", "/boards/"+boardID, domain.BoardUpdateModel{
		Name: &updatedName,
	}, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil || resp.Data.Name != updatedName {
		t.Fatalf("expected name '%s', got %v", updatedName, resp.Data)
	}
}

func TestBoards_Reorder_Success(t *testing.T) {
	// Create org, project, sprint, and boards
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

	// Reorder boards
	statusCode, resp := do[[]domain.BoardModel](t, "PATCH", "/boards/reorder", map[string]interface{}{
		"boards": []map[string]interface{}{
			{
				"id":       uuidToString(board1.ID),
				"position": 1,
			},
			{
				"id":       uuidToString(board2.ID),
				"position": 0,
			},
		},
	}, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil || len(*resp.Data) != 2 {
		t.Fatalf("expected 2 boards in response")
	}
}

func TestBoards_Delete_Success(t *testing.T) {
	// Create org, project, sprint, and board
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

	board := createBoard(t, sprintID, tokens.AccessToken, randomBoardName())
	boardID := uuidToString(board.ID)

	// Delete the board
	code, _ := do[domain.BoardModel](t, "DELETE", "/boards/"+boardID, nil, tokens.AccessToken)

	if code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", code)
	}

	// Verify board is deleted
	code, _ = do[domain.BoardModel](t, "GET", "/boards/"+boardID, nil, tokens.AccessToken)
	if code != http.StatusNotFound {
		t.Fatalf("expected deleted board to return 404, got %d", code)
	}
}

func TestBoards_Delete_NotFound(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	nonExistentID := "550e8400-e29b-41d4-a716-446655440000"

	statusCode, _ := do[domain.BoardModel](t, "DELETE", "/boards/"+nonExistentID, nil, tokens.AccessToken)

	if statusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", statusCode)
	}
}

func TestBoards_Create_MultipleInSprint_AutoPosition(t *testing.T) {
	// Create org, project, and sprint
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

	// Create multiple boards and verify positions
	board1 := createBoard(t, sprintID, tokens.AccessToken, randomBoardName())
	if board1.Position != 0 {
		t.Fatalf("expected first board position 0, got %d", board1.Position)
	}

	board2 := createBoard(t, sprintID, tokens.AccessToken, randomBoardName())
	if board2.Position != 1 {
		t.Fatalf("expected second board position 1, got %d", board2.Position)
	}

	board3 := createBoard(t, sprintID, tokens.AccessToken, randomBoardName())
	if board3.Position != 2 {
		t.Fatalf("expected third board position 2, got %d", board3.Position)
	}
}
