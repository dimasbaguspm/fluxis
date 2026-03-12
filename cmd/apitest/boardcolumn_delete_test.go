package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestBoardColumn_Delete_Success(t *testing.T) {
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
	column := createBoardColumn(t, uuidToString(board.ID), tokens.AccessToken, randomBoardColumnName())

	code, _ := do[interface{}](t, "DELETE", "/boards/"+uuidToString(board.ID)+"/columns/"+uuidToString(column.ID), nil, tokens.AccessToken)

	if code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", code)
	}
}

func TestBoardColumn_Delete_NotFound(t *testing.T) {
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

	invalidColumnID := "550e8400-e29b-41d4-a716-446655440000"
	code, _ := do[interface{}](t, "DELETE", "/boards/"+uuidToString(board.ID)+"/columns/"+invalidColumnID, nil, tokens.AccessToken)

	if code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", code)
	}
}

func TestBoardColumn_Delete_InvalidBoard(t *testing.T) {
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
	column := createBoardColumn(t, uuidToString(board.ID), tokens.AccessToken, randomBoardColumnName())

	// Try to delete column with wrong board ID
	wrongBoardID := "550e8400-e29b-41d4-a716-446655440000"
	code, _ := do[interface{}](t, "DELETE", "/boards/"+wrongBoardID+"/columns/"+uuidToString(column.ID), nil, tokens.AccessToken)

	if code != http.StatusNotFound {
		t.Fatalf("expected status 404 for column not in board, got %d", code)
	}
}
