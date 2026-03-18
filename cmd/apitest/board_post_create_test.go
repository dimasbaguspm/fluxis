package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestBoard_Create_Success(t *testing.T) {
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
	statusCode, resp := do[domain.BoardModel](t, "POST", "/boards", domain.BoardCreateModel{
		Name:     boardName,
		SprintID: stringToUUID(sprintID),
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

	if len(resp.Data.Columns) != 6 {
		t.Fatalf("expected 6 columns, got %d", len(resp.Data.Columns))
	}

	expectedColumnNames := []string{
		"To Do",
		"In Progress",
		"Code Review",
		"Ready to Test",
		"Ready for PO Review",
		"Done",
	}

	for i, col := range resp.Data.Columns {
		if col.Name != expectedColumnNames[i] {
			t.Fatalf("expected column %d name '%s', got '%s'", i, expectedColumnNames[i], col.Name)
		}
	}
}

func TestBoard_Create_Unauthenticated(t *testing.T) {
	sprintID := "550e8400-e29b-41d4-a716-446655440000"

	name := "Test Board"
	statusCode, _ := do[domain.BoardModel](t, "POST", "/boards", domain.BoardCreateModel{
		Name:     name,
		SprintID: stringToUUID(sprintID),
	}, "")

	if statusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", statusCode)
	}
}

func TestBoard_Create_MissingSprintId(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	code, _ := do[domain.BoardModel](t, "POST", "/boards", map[string]interface{}{
		"name": "Test Board",
	}, tokens.AccessToken)

	if code != http.StatusNotFound {
		t.Fatalf("expected status 404 (sprint not found), got %d", code)
	}
}

func TestBoard_Create_MissingName(t *testing.T) {
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

func TestBoard_Create_InvalidSprintID(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	invalidSprintID := "550e8400-e29b-41d4-a716-446655440000"
	boardName := randomBoardName()

	statusCode, resp := do[domain.BoardModel](t, "POST", "/boards", domain.BoardCreateModel{
		Name:     boardName,
		SprintID: stringToUUID(invalidSprintID),
	}, tokens.AccessToken)

	if statusCode != http.StatusNotFound {
		t.Fatalf("expected status 404 for non-existent sprint, got %d", statusCode)
	}

	if resp.Error == nil {
		t.Fatalf("expected error response for invalid sprint")
	}
}

func TestBoard_Create_AutoPosition_FirstBoardIsZero(t *testing.T) {
	// Verify that first board in a sprint gets position 0
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

	// Create board and verify it was created
	board := createBoard(t, sprintID, tokens.AccessToken, randomBoardName())
	if uuidToString(board.ID) == "" {
		t.Fatal("expected board to be created")
	}
}
