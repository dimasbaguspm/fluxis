package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestBoardColumn_List_Success(t *testing.T) {
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

	// Create multiple columns
	col1 := createBoardColumn(t, uuidToString(board.ID), tokens.AccessToken, "Column 1")
	col2 := createBoardColumn(t, uuidToString(board.ID), tokens.AccessToken, "Column 2")

	statusCode, resp := do[domain.BoardColumnsPagedModel](t, "GET", "/boards/"+uuidToString(board.ID)+"/columns", nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected columns data")
	}

	// 6 default columns + 2 created = 8 total
	if len(resp.Data.Items) != 8 {
		t.Fatalf("expected 8 columns (6 default + 2 created), got %d", len(resp.Data.Items))
	}

	// Verify created columns are at the end (after default columns)
	if resp.Data.Items[6].ID != col1.ID {
		t.Fatalf("expected 7th column to be col1")
	}
	if resp.Data.Items[7].ID != col2.ID {
		t.Fatalf("expected 8th column to be col2")
	}

	if resp.Data.TotalCount != 8 {
		t.Fatalf("expected total count 8, got %d", resp.Data.TotalCount)
	}

	if resp.Data.PageNumber != 1 {
		t.Fatalf("expected page number 1, got %d", resp.Data.PageNumber)
	}
}

func TestBoardColumn_List_EmptyBoard(t *testing.T) {
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

	statusCode, resp := do[domain.BoardColumnsPagedModel](t, "GET", "/boards/"+uuidToString(board.ID)+"/columns", nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected columns data")
	}

	// Board is created with 6 default columns
	if len(resp.Data.Items) != 6 {
		t.Fatalf("expected 6 default columns, got %d", len(resp.Data.Items))
	}

	if resp.Data.TotalCount != 6 {
		t.Fatalf("expected total count 6, got %d", resp.Data.TotalCount)
	}
}

func TestBoardColumn_List_WithNameFilter(t *testing.T) {
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
	boardID := uuidToString(board.ID)

	createBoardColumn(t, boardID, tokens.AccessToken, "To Do Column")
	createBoardColumn(t, boardID, tokens.AccessToken, "In Progress Column")
	createBoardColumn(t, boardID, tokens.AccessToken, "Done Column")

	// Filter by "In Progress"
	statusCode, resp := do[domain.BoardColumnsPagedModel](t, "GET", "/boards/"+boardID+"/columns?name=In+Progress", nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected columns data")
	}

	// Should match both the default "In Progress" column and "In Progress Column"
	if len(resp.Data.Items) != 2 {
		t.Fatalf("expected 2 columns matching 'In Progress' (1 default + 1 created), got %d", len(resp.Data.Items))
	}

	// Check that "In Progress Column" is in the results
	found := false
	for _, item := range resp.Data.Items {
		if item.Name == "In Progress Column" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected 'In Progress Column' in results")
	}

	if resp.Data.TotalCount != 2 {
		t.Fatalf("expected total count 2, got %d", resp.Data.TotalCount)
	}
}

func TestBoardColumn_List_WithNameFilterNoMatch(t *testing.T) {
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
	createBoardColumn(t, uuidToString(board.ID), tokens.AccessToken, "Column 2")

	// Filter by non-existent name
	statusCode, resp := do[domain.BoardColumnsPagedModel](t, "GET", "/boards/"+uuidToString(board.ID)+"/columns?name=NonExistent", nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected columns data")
	}

	if len(resp.Data.Items) != 0 {
		t.Fatalf("expected 0 columns for non-existent name, got %d", len(resp.Data.Items))
	}

	if resp.Data.TotalCount != 0 {
		t.Fatalf("expected total count 0, got %d", resp.Data.TotalCount)
	}
}

func TestBoardColumn_List_WithPagination(t *testing.T) {
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
	boardID := uuidToString(board.ID)

	// Create 4 columns (board already has 6 default columns, so total will be 10)
	for i := 0; i < 4; i++ {
		createBoardColumn(t, boardID, tokens.AccessToken, "Column "+string(rune('1'+i)))
	}

	// Get first page with size 2
	statusCode, resp := do[domain.BoardColumnsPagedModel](t, "GET", "/boards/"+boardID+"/columns?pageNumber=1&pageSize=2", nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if len(resp.Data.Items) != 2 {
		t.Fatalf("expected 2 columns on page 1, got %d", len(resp.Data.Items))
	}

	// Total count should be 10 (6 default + 4 created)
	if resp.Data.TotalCount != 10 {
		t.Fatalf("expected total count 10, got %d", resp.Data.TotalCount)
	}

	if resp.Data.TotalPages != 5 {
		t.Fatalf("expected 5 total pages, got %d", resp.Data.TotalPages)
	}

	// Get second page
	statusCode, resp2 := do[domain.BoardColumnsPagedModel](t, "GET", "/boards/"+boardID+"/columns?pageNumber=2&pageSize=2", nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if len(resp2.Data.Items) != 2 {
		t.Fatalf("expected 2 columns on page 2, got %d", len(resp2.Data.Items))
	}

	if resp2.Data.PageNumber != 2 {
		t.Fatalf("expected page number 2, got %d", resp2.Data.PageNumber)
	}
}

func TestBoardColumn_List_InvalidBoardID(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, _ := do[domain.BoardColumnsPagedModel](t, "GET", "/boards/invalid-id/columns", nil, tokens.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}
}

func TestBoardColumn_List_NonExistentBoard(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	nonExistentBoardID := "550e8400-e29b-41d4-a716-446655440000"
	statusCode, _ := do[domain.BoardColumnsPagedModel](t, "GET", "/boards/"+nonExistentBoardID+"/columns", nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}
}

func TestBoardColumn_List_FilterByID(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	project := createProject(t, uuidToString(orgResp.Data.ID), tokens.AccessToken, randomProjectKey(), "Test Project", "private")
	sprint := createSprint(t, uuidToString(project.ID), tokens.AccessToken, randomSprintName())
	board := createBoard(t, uuidToString(sprint.ID), tokens.AccessToken, "Test Board")
	boardID := uuidToString(board.ID)

	// Create columns
	c1 := createBoardColumn(t, boardID, tokens.AccessToken, "Column 1")
	createBoardColumn(t, boardID, tokens.AccessToken, "Column 2")

	// Filter by specific column ID
	c1ID := uuidToString(c1.ID)
	statusCode, resp := do[domain.BoardColumnsPagedModel](t, "GET", "/boards/"+boardID+"/columns?id="+c1ID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if len(resp.Data.Items) != 1 {
		t.Fatalf("expected 1 column, got %d", len(resp.Data.Items))
	}

	if resp.Data.Items[0].ID != c1.ID {
		t.Fatalf("expected column %s, got %s", c1ID, uuidToString(resp.Data.Items[0].ID))
	}
}

func TestBoardColumn_List_FilterByMultipleIDs(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	project := createProject(t, uuidToString(orgResp.Data.ID), tokens.AccessToken, randomProjectKey(), "Test Project", "private")
	sprint := createSprint(t, uuidToString(project.ID), tokens.AccessToken, randomSprintName())
	board := createBoard(t, uuidToString(sprint.ID), tokens.AccessToken, "Test Board")
	boardID := uuidToString(board.ID)

	// Create columns
	c1 := createBoardColumn(t, boardID, tokens.AccessToken, "Column 1")
	c2 := createBoardColumn(t, boardID, tokens.AccessToken, "Column 2")
	createBoardColumn(t, boardID, tokens.AccessToken, "Column 3")

	// Filter by multiple column IDs
	c1ID := uuidToString(c1.ID)
	c2ID := uuidToString(c2.ID)
	statusCode, resp := do[domain.BoardColumnsPagedModel](t, "GET", "/boards/"+boardID+"/columns?id="+c1ID+"&id="+c2ID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if len(resp.Data.Items) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(resp.Data.Items))
	}
}

func TestBoardColumn_List_Unauthenticated(t *testing.T) {
	boardID := "550e8400-e29b-41d4-a716-446655440000"
	statusCode, _ := do[domain.BoardColumnsPagedModel](t, "GET", "/boards/"+boardID+"/columns", nil, "")

	if statusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", statusCode)
	}
}
