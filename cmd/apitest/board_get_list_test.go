package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestBoard_List_BySprint(t *testing.T) {
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

	// Create board
	board := createBoard(t, sprintID, tokens.AccessToken, randomBoardName())

	// List boards for this sprint
	statusCode, resp := do[domain.BoardsPagedModel](t, "GET", "/boards?sprintId="+sprintID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected board list data")
	}

	if len(resp.Data.Items) != 1 {
		t.Fatalf("expected exactly 1 board for sprint, got %d", len(resp.Data.Items))
	}

	if resp.Data.Items[0].ID != board.ID {
		t.Fatalf("expected board ID to match")
	}

	if resp.Data.PageNumber != 1 {
		t.Fatalf("expected page number 1, got %d", resp.Data.PageNumber)
	}
}

func TestBoard_List_EmptySprint(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)
	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	project := createProject(t, uuidToString(orgResp.Data.ID), tokens.AccessToken, randomProjectKey(), "Test Project", "private")
	sprint := createSprint(t, uuidToString(project.ID), tokens.AccessToken, randomSprintName())

	statusCode, resp := do[domain.BoardsPagedModel](t, "GET", "/boards?sprintId="+uuidToString(sprint.ID), nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected board list data")
	}

	if len(resp.Data.Items) != 0 {
		t.Fatalf("expected 0 boards, got %d", len(resp.Data.Items))
	}

	if resp.Data.TotalCount != 0 {
		t.Fatalf("expected total count 0, got %d", resp.Data.TotalCount)
	}
}

func TestBoard_List_WithNameFilter(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)
	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	project := createProject(t, uuidToString(orgResp.Data.ID), tokens.AccessToken, randomProjectKey(), "Test Project", "private")
	projectID := uuidToString(project.ID)

	// Create sprints and boards
	sprint1 := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprint1ID := uuidToString(sprint1.ID)
	createBoard(t, sprint1ID, tokens.AccessToken, "Frontend Board")

	sprint2 := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprint2ID := uuidToString(sprint2.ID)
	createBoard(t, sprint2ID, tokens.AccessToken, "Backend Board")

	sprint3 := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprint3ID := uuidToString(sprint3.ID)
	createBoard(t, sprint3ID, tokens.AccessToken, "Testing Board")

	// Filter by "Frontend"
	statusCode, resp := do[domain.BoardsPagedModel](t, "GET", "/boards?name=Frontend", nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected board list data")
	}

	if len(resp.Data.Items) < 1 {
		t.Fatalf("expected at least 1 board matching 'Frontend', got %d", len(resp.Data.Items))
	}

	found := false
	for _, board := range resp.Data.Items {
		if board.Name == "Frontend Board" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected to find 'Frontend Board' in results")
	}
}

func TestBoard_List_WithNameFilterNoMatch(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)
	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	project := createProject(t, uuidToString(orgResp.Data.ID), tokens.AccessToken, randomProjectKey(), "Test Project", "private")
	projectID := uuidToString(project.ID)

	sprint1 := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprint1ID := uuidToString(sprint1.ID)
	createBoard(t, sprint1ID, tokens.AccessToken, "Frontend Board")

	sprint2 := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprint2ID := uuidToString(sprint2.ID)
	createBoard(t, sprint2ID, tokens.AccessToken, "Backend Board")

	// Filter by non-existent name
	statusCode, resp := do[domain.BoardsPagedModel](t, "GET", "/boards?name=NonExistent", nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected board list data")
	}

	if len(resp.Data.Items) != 0 {
		t.Fatalf("expected 0 boards for non-existent name, got %d", len(resp.Data.Items))
	}
}

func TestBoard_List_WithPagination(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)
	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	project := createProject(t, uuidToString(orgResp.Data.ID), tokens.AccessToken, randomProjectKey(), "Test Project", "private")
	projectID := uuidToString(project.ID)

	// Create 5 boards (one per sprint)
	for i := 0; i < 5; i++ {
		sprint := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
		sprintID := uuidToString(sprint.ID)
		createBoard(t, sprintID, tokens.AccessToken, "Board "+string(rune('A'+i)))
	}

	// Get first page with size 2
	statusCode, resp := do[domain.BoardsPagedModel](t, "GET", "/boards?pageNumber=1&pageSize=2", nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if len(resp.Data.Items) != 2 {
		t.Fatalf("expected 2 boards on page 1, got %d", len(resp.Data.Items))
	}

	if resp.Data.TotalCount < 5 {
		t.Fatalf("expected total count at least 5, got %d", resp.Data.TotalCount)
	}

	if resp.Data.TotalPages < 3 {
		t.Fatalf("expected at least 3 total pages, got %d", resp.Data.TotalPages)
	}

	// Get second page
	statusCode, resp2 := do[domain.BoardsPagedModel](t, "GET", "/boards?pageNumber=2&pageSize=2", nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if len(resp2.Data.Items) != 2 {
		t.Fatalf("expected 2 boards on page 2, got %d", len(resp2.Data.Items))
	}

	if resp2.Data.PageNumber != 2 {
		t.Fatalf("expected page number 2, got %d", resp2.Data.PageNumber)
	}
}

func TestBoard_List_FilterByID(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	project := createProject(t, uuidToString(orgResp.Data.ID), tokens.AccessToken, randomProjectKey(), "Test Project", "private")
	projectID := uuidToString(project.ID)

	// Create boards in different sprints
	sprint1 := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprint1ID := uuidToString(sprint1.ID)
	b1 := createBoard(t, sprint1ID, tokens.AccessToken, "Board 1")

	sprint2 := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprint2ID := uuidToString(sprint2.ID)
	createBoard(t, sprint2ID, tokens.AccessToken, "Board 2")

	// Filter by specific board ID
	b1ID := uuidToString(b1.ID)
	statusCode, resp := do[domain.BoardsPagedModel](t, "GET", "/boards?id="+b1ID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if len(resp.Data.Items) != 1 {
		t.Fatalf("expected 1 board, got %d", len(resp.Data.Items))
	}

	if resp.Data.Items[0].ID != b1.ID {
		t.Fatalf("expected board %s, got %s", b1ID, uuidToString(resp.Data.Items[0].ID))
	}
}

func TestBoard_List_FilterByMultipleIDs(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	project := createProject(t, uuidToString(orgResp.Data.ID), tokens.AccessToken, randomProjectKey(), "Test Project", "private")
	projectID := uuidToString(project.ID)

	// Create boards in different sprints
	sprint1 := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprint1ID := uuidToString(sprint1.ID)
	b1 := createBoard(t, sprint1ID, tokens.AccessToken, "Board 1")

	sprint2 := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprint2ID := uuidToString(sprint2.ID)
	b2 := createBoard(t, sprint2ID, tokens.AccessToken, "Board 2")

	sprint3 := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprint3ID := uuidToString(sprint3.ID)
	createBoard(t, sprint3ID, tokens.AccessToken, "Board 3")

	// Filter by multiple board IDs
	b1ID := uuidToString(b1.ID)
	b2ID := uuidToString(b2.ID)
	statusCode, resp := do[domain.BoardsPagedModel](t, "GET", "/boards?id="+b1ID+"&id="+b2ID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if len(resp.Data.Items) != 2 {
		t.Fatalf("expected 2 boards, got %d", len(resp.Data.Items))
	}
}

func TestBoard_List_Unauthenticated(t *testing.T) {
	statusCode, _ := do[domain.BoardsPagedModel](t, "GET", "/boards?sprintId=550e8400-e29b-41d4-a716-446655440000", nil, "")

	if statusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", statusCode)
	}
}

func TestBoard_List_ByProjectId(t *testing.T) {
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

	// Create multiple sprints with boards
	sprint1 := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprint1ID := uuidToString(sprint1.ID)
	sprint2 := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprint2ID := uuidToString(sprint2.ID)
	sprint3 := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprint3ID := uuidToString(sprint3.ID)

	// Create boards in different sprints
	createBoard(t, sprint1ID, tokens.AccessToken, "Sprint 1 Board A")
	createBoard(t, sprint2ID, tokens.AccessToken, "Sprint 2 Board B")
	createBoard(t, sprint3ID, tokens.AccessToken, "Sprint 3 Board C")

	// List boards by projectId (should get all boards from all sprints)
	statusCode, resp := do[domain.BoardsPagedModel](t, "GET", "/boards?projectId="+projectID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected board list data")
	}

	if len(resp.Data.Items) != 3 {
		t.Fatalf("expected 3 boards from project, got %d", len(resp.Data.Items))
	}

	if resp.Data.TotalCount != 3 {
		t.Fatalf("expected total count 3, got %d", resp.Data.TotalCount)
	}
}

func TestBoard_List_ByProjectAndSprintId(t *testing.T) {
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

	// Create multiple sprints
	sprint1 := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprint1ID := uuidToString(sprint1.ID)
	sprint2 := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprint2ID := uuidToString(sprint2.ID)
	sprint3 := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprint3ID := uuidToString(sprint3.ID)

	// Create boards in different sprints
	createBoard(t, sprint1ID, tokens.AccessToken, "Sprint 1 Board A")
	createBoard(t, sprint2ID, tokens.AccessToken, "Sprint 2 Board B")
	createBoard(t, sprint3ID, tokens.AccessToken, "Sprint 3 Board C")

	// List boards by both projectId and sprintId (should get only boards from sprint1)
	statusCode, resp := do[domain.BoardsPagedModel](t, "GET", "/boards?projectId="+projectID+"&sprintId="+sprint1ID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected board list data")
	}

	if len(resp.Data.Items) != 1 {
		t.Fatalf("expected 1 board from sprint 1, got %d", len(resp.Data.Items))
	}

	if resp.Data.TotalCount != 1 {
		t.Fatalf("expected total count 1, got %d", resp.Data.TotalCount)
	}
}
