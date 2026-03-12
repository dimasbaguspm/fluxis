package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestBoard_Update_Success(t *testing.T) {
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
		Name: updatedName,
	}, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil || resp.Data.Name != updatedName {
		t.Fatalf("expected name '%s', got %v", updatedName, resp.Data)
	}
}

func TestBoard_Update_WithSprintChange_Success(t *testing.T) {
	// Create org, project, and multiple sprints
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

	sprint1 := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprintID1 := uuidToString(sprint1.ID)
	sprint2 := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprintID2 := uuidToString(sprint2.ID)

	board := createBoard(t, sprintID1, tokens.AccessToken, "Original Board")
	boardID := uuidToString(board.ID)

	// Update board to move to different sprint (concurrent validation)
	updatedName := "Updated Board Name"
	newSprintID := stringToUUID(sprintID2)
	statusCode, resp := do[domain.BoardModel](t, "PATCH", "/boards/"+boardID, domain.BoardUpdateModel{
		Name:     updatedName,
		SprintID: newSprintID,
	}, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil || resp.Data.Name != updatedName {
		t.Fatalf("expected name '%s', got '%s'", updatedName, resp.Data.Name)
	}

	if uuidToString(resp.Data.SprintID) != sprintID2 {
		t.Fatalf("expected sprint ID '%s', got '%s'", sprintID2, uuidToString(resp.Data.SprintID))
	}
}

func TestBoard_Update_InvalidSprintID(t *testing.T) {
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

	board := createBoard(t, sprintID, tokens.AccessToken, "Original Board")
	boardID := uuidToString(board.ID)

	// Try to update with non-existent sprint ID
	invalidSprintID := stringToUUID("550e8400-e29b-41d4-a716-446655440000")
	updatedName := "Updated Board"
	statusCode, resp := do[domain.BoardModel](t, "PATCH", "/boards/"+boardID, domain.BoardUpdateModel{
		Name:     updatedName,
		SprintID: invalidSprintID,
	}, tokens.AccessToken)

	if statusCode != http.StatusNotFound {
		t.Fatalf("expected status 404 for non-existent sprint, got %d", statusCode)
	}

	if resp.Error == nil {
		t.Fatalf("expected error response for invalid sprint")
	}
}

func TestBoard_Update_AllowsPartialUpdate(t *testing.T) {
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
	sprint1 := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprintID1 := uuidToString(sprint1.ID)
	sprint2 := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprintID2 := uuidToString(sprint2.ID)

	originalName := "Original Board"
	board := createBoard(t, sprintID1, tokens.AccessToken, originalName)
	boardID := uuidToString(board.ID)

	// Update only sprint, omitting name (partial update should work)
	status, resp := do[domain.BoardModel](t, "PATCH", "/boards/"+boardID, domain.BoardUpdateModel{
		SprintID: stringToUUID(sprintID2),
	}, tokens.AccessToken)

	if status != http.StatusOK {
		t.Fatalf("expected status 200, got %d", status)
	}

	// Sprint should be updated
	if uuidToString(resp.Data.SprintID) != sprintID2 {
		t.Fatalf("expected sprint ID '%s', got '%s'", sprintID2, uuidToString(resp.Data.SprintID))
	}

	// Name should remain unchanged
	if resp.Data.Name != originalName {
		t.Fatalf("expected name to remain '%s', got '%s'", originalName, resp.Data.Name)
	}
}

func TestBoard_Update_NotFound(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	nonExistentID := "550e8400-e29b-41d4-a716-446655440000"

	statusCode, _ := do[domain.BoardModel](t, "PATCH", "/boards/"+nonExistentID, domain.BoardUpdateModel{
		Name: "Updated Board",
	}, tokens.AccessToken)

	if statusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", statusCode)
	}
}

func TestBoard_Update_InvalidUUID(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, _ := do[domain.BoardModel](t, "PATCH", "/boards/not-a-uuid", domain.BoardUpdateModel{
		Name: "Updated Board",
	}, tokens.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}
}

func TestBoard_Update_Unauthenticated(t *testing.T) {
	boardID := "550e8400-e29b-41d4-a716-446655440000"

	statusCode, _ := do[domain.BoardModel](t, "PATCH", "/boards/"+boardID, domain.BoardUpdateModel{
		Name: "Updated Board",
	}, "")

	if statusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", statusCode)
	}
}
