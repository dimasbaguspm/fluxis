package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestBoardColumn_Reorder_Success(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)
	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	project := createProject(t, uuidToString(orgResp.Data.ID), tokens.AccessToken, randomProjectKey(), "Test Project", "private")
	sprint := createSprint(t, uuidToString(project.ID), tokens.AccessToken, randomSprintName())
	board := createBoard(t, uuidToString(sprint.ID), tokens.AccessToken, randomBoardName())

	// Create 3 columns
	col1 := createBoardColumn(t, uuidToString(board.ID), tokens.AccessToken, "Column 1")
	col2 := createBoardColumn(t, uuidToString(board.ID), tokens.AccessToken, "Column 2")
	col3 := createBoardColumn(t, uuidToString(board.ID), tokens.AccessToken, "Column 3")

	// Verify initial order
	statusCode, listResp := do[[]domain.BoardColumnModel](t, "GET", "/boards/"+uuidToString(board.ID)+"/columns", nil, tokens.AccessToken)
	if statusCode != http.StatusOK || len(*listResp.Data) != 3 {
		t.Fatalf("expected 3 columns")
	}
	if (*listResp.Data)[0].ID != col1.ID || (*listResp.Data)[1].ID != col2.ID || (*listResp.Data)[2].ID != col3.ID {
		t.Fatalf("initial order incorrect")
	}

	// Reorder to [col3, col1, col2]
	reorderPayload := domain.BoardColumnReorderModel{col3.ID, col1.ID, col2.ID}
	statusCode, reorderResp := do[[]domain.BoardColumnModel](t, "PATCH", "/boards/"+uuidToString(board.ID)+"/columns/reorder", reorderPayload, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, reorderResp.Error)
	}

	if reorderResp.Data == nil {
		t.Fatal("expected columns data")
	}

	if len(*reorderResp.Data) != 3 {
		t.Fatalf("expected 3 columns in response, got %d", len(*reorderResp.Data))
	}

	// Verify positions were updated correctly
	if (*reorderResp.Data)[0].Position != 0 || (*reorderResp.Data)[0].ID != col3.ID {
		t.Fatalf("expected col3 at position 0")
	}
	if (*reorderResp.Data)[1].Position != 1 || (*reorderResp.Data)[1].ID != col1.ID {
		t.Fatalf("expected col1 at position 1")
	}
	if (*reorderResp.Data)[2].Position != 2 || (*reorderResp.Data)[2].ID != col2.ID {
		t.Fatalf("expected col2 at position 2")
	}

	// Verify list endpoint returns same order
	statusCode, listResp = do[[]domain.BoardColumnModel](t, "GET", "/boards/"+uuidToString(board.ID)+"/columns", nil, tokens.AccessToken)
	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200")
	}
	if (*listResp.Data)[0].ID != col3.ID || (*listResp.Data)[1].ID != col1.ID || (*listResp.Data)[2].ID != col2.ID {
		t.Fatalf("list order does not match reorder result")
	}
}

func TestBoardColumn_Reorder_NonExistentBoard(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	nonExistentBoardID := "550e8400-e29b-41d4-a716-446655440000"
	col1ID := "550e8400-e29b-41d4-a716-446655440001"
	col2ID := "550e8400-e29b-41d4-a716-446655440002"

	reorderPayload := domain.BoardColumnReorderModel{stringToUUID(col1ID), stringToUUID(col2ID)}
	statusCode, _ := do[[]domain.BoardColumnModel](t, "PATCH", "/boards/"+nonExistentBoardID+"/columns/reorder", reorderPayload, tokens.AccessToken)

	if statusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", statusCode)
	}
}

func TestBoardColumn_Reorder_InvalidBoardUUID(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	col1ID := "550e8400-e29b-41d4-a716-446655440001"
	col2ID := "550e8400-e29b-41d4-a716-446655440002"

	reorderPayload := domain.BoardColumnReorderModel{stringToUUID(col1ID), stringToUUID(col2ID)}
	statusCode, _ := do[[]domain.BoardColumnModel](t, "PATCH", "/boards/not-a-uuid/columns/reorder", reorderPayload, tokens.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}
}

func TestBoardColumn_Reorder_NonExistentColumn(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)
	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	project := createProject(t, uuidToString(orgResp.Data.ID), tokens.AccessToken, randomProjectKey(), "Test Project", "private")
	sprint := createSprint(t, uuidToString(project.ID), tokens.AccessToken, randomSprintName())
	board := createBoard(t, uuidToString(sprint.ID), tokens.AccessToken, randomBoardName())

	col1 := createBoardColumn(t, uuidToString(board.ID), tokens.AccessToken, "Column 1")
	nonExistentColumnID := "550e8400-e29b-41d4-a716-446655440000"

	reorderPayload := domain.BoardColumnReorderModel{col1.ID, stringToUUID(nonExistentColumnID)}
	statusCode, resp := do[[]domain.BoardColumnModel](t, "PATCH", "/boards/"+uuidToString(board.ID)+"/columns/reorder", reorderPayload, tokens.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %v", statusCode, resp.Error)
	}
}

func TestBoardColumn_Reorder_ColumnFromDifferentBoard(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)
	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	project := createProject(t, uuidToString(orgResp.Data.ID), tokens.AccessToken, randomProjectKey(), "Test Project", "private")
	sprint := createSprint(t, uuidToString(project.ID), tokens.AccessToken, randomSprintName())

	// Create 2 boards
	board1 := createBoard(t, uuidToString(sprint.ID), tokens.AccessToken, randomBoardName())
	board2 := createBoard(t, uuidToString(sprint.ID), tokens.AccessToken, randomBoardName())

	// Create columns in both boards
	col1 := createBoardColumn(t, uuidToString(board1.ID), tokens.AccessToken, "Column 1")
	col2FromBoard2 := createBoardColumn(t, uuidToString(board2.ID), tokens.AccessToken, "Column From Board 2")

	// Try to reorder board1's columns using column from board2
	reorderPayload := domain.BoardColumnReorderModel{col1.ID, col2FromBoard2.ID}
	statusCode, resp := do[[]domain.BoardColumnModel](t, "PATCH", "/boards/"+uuidToString(board1.ID)+"/columns/reorder", reorderPayload, tokens.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %v", statusCode, resp.Error)
	}
}

func TestBoardColumn_Reorder_EmptyArray(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)
	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	project := createProject(t, uuidToString(orgResp.Data.ID), tokens.AccessToken, randomProjectKey(), "Test Project", "private")
	sprint := createSprint(t, uuidToString(project.ID), tokens.AccessToken, randomSprintName())
	board := createBoard(t, uuidToString(sprint.ID), tokens.AccessToken, randomBoardName())

	createBoardColumn(t, uuidToString(board.ID), tokens.AccessToken, "Column 1")

	reorderPayload := domain.BoardColumnReorderModel{}
	statusCode, resp := do[[]domain.BoardColumnModel](t, "PATCH", "/boards/"+uuidToString(board.ID)+"/columns/reorder", reorderPayload, tokens.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400 for empty array, got %d: %v", statusCode, resp.Error)
	}
}

func TestBoardColumn_Reorder_PartialUpdate(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)
	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	project := createProject(t, uuidToString(orgResp.Data.ID), tokens.AccessToken, randomProjectKey(), "Test Project", "private")
	sprint := createSprint(t, uuidToString(project.ID), tokens.AccessToken, randomSprintName())
	board := createBoard(t, uuidToString(sprint.ID), tokens.AccessToken, randomBoardName())

	// Create 3 columns
	col1 := createBoardColumn(t, uuidToString(board.ID), tokens.AccessToken, "Column 1")
	_ = createBoardColumn(t, uuidToString(board.ID), tokens.AccessToken, "Column 2") // col2 not used in reorder
	col3 := createBoardColumn(t, uuidToString(board.ID), tokens.AccessToken, "Column 3")

	// Try to reorder only 2 of 3 columns
	reorderPayload := domain.BoardColumnReorderModel{col1.ID, col3.ID}
	statusCode, resp := do[[]domain.BoardColumnModel](t, "PATCH", "/boards/"+uuidToString(board.ID)+"/columns/reorder", reorderPayload, tokens.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400 for partial update, got %d: %v", statusCode, resp.Error)
	}
}

func TestBoardColumn_Reorder_Unauthenticated(t *testing.T) {
	boardID := "550e8400-e29b-41d4-a716-446655440000"
	col1ID := "550e8400-e29b-41d4-a716-446655440001"
	col2ID := "550e8400-e29b-41d4-a716-446655440002"

	reorderPayload := domain.BoardColumnReorderModel{stringToUUID(col1ID), stringToUUID(col2ID)}
	statusCode, _ := do[[]domain.BoardColumnModel](t, "PATCH", "/boards/"+boardID+"/columns/reorder", reorderPayload, "")

	if statusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", statusCode)
	}
}
