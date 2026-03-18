package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestBoard_Create_ConflictOnlyOnePerSprint(t *testing.T) {
	// Verify that only one board can exist per sprint
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

	// Create first board - should succeed
	board1 := createBoard(t, sprintID, tokens.AccessToken, randomBoardName())
	if uuidToString(board1.ID) == "" {
		t.Fatal("failed to create first board")
	}

	// Try to create second board in same sprint - should return 409 Conflict
	statusCode, resp := do[domain.BoardModel](t, "POST", "/boards", domain.BoardCreateModel{
		Name:     randomBoardName(),
		SprintID: stringToUUID(sprintID),
	}, tokens.AccessToken)

	if statusCode != http.StatusConflict {
		t.Fatalf("expected status 409 Conflict, got %d", statusCode)
	}

	if resp.Error == nil {
		t.Fatalf("expected error response")
	}

	if resp.Error.Message != "sprint already has a board" {
		t.Fatalf("expected error message 'sprint already has a board', got '%s'", resp.Error.Message)
	}
}

func TestBoard_Create_DifferentSprintAllowsBoth(t *testing.T) {
	// Verify that different sprints can each have their own board
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

	// Create board in sprint 1
	board1 := createBoard(t, sprint1ID, tokens.AccessToken, randomBoardName())
	if uuidToString(board1.ID) == "" {
		t.Fatal("failed to create board in sprint 1")
	}

	// Create board in sprint 2 - should succeed
	statusCode, resp := do[domain.BoardModel](t, "POST", "/boards", domain.BoardCreateModel{
		Name:     randomBoardName(),
		SprintID: stringToUUID(sprint2ID),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated {
		t.Fatalf("expected status 201 for different sprint, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil || uuidToString(resp.Data.ID) == "" {
		t.Fatalf("expected board to be created in sprint 2")
	}
}

func TestBoard_Create_OnePerSprintAfterDelete(t *testing.T) {
	// Verify that after deleting a board, another board can be created in the same sprint
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

	// Create first board
	board1 := createBoard(t, sprintID, tokens.AccessToken, randomBoardName())
	board1ID := uuidToString(board1.ID)

	// Delete the board
	deleteCode, _ := do[interface{}](t, "DELETE", "/boards/"+board1ID, nil, tokens.AccessToken)
	if deleteCode != http.StatusNoContent {
		t.Fatalf("failed to delete board: status %d", deleteCode)
	}

	// Create second board - should succeed
	statusCode, resp := do[domain.BoardModel](t, "POST", "/boards", domain.BoardCreateModel{
		Name:     randomBoardName(),
		SprintID: stringToUUID(sprintID),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated {
		t.Fatalf("expected status 201 after deleting first board, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil || uuidToString(resp.Data.ID) == "" {
		t.Fatalf("expected board to be created after deletion")
	}
}
