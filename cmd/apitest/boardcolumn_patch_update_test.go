package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestBoardColumn_Update_Success(t *testing.T) {
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

	newName := "Updated Column"
	statusCode, resp := do[domain.BoardColumnModel](t, "PATCH", "/boards/"+uuidToString(board.ID)+"/columns/"+uuidToString(column.ID), domain.BoardColumnUpdateModel{
		Name: newName,
	}, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected column data")
	}

	if resp.Data.Name != newName {
		t.Fatalf("expected name '%s', got '%s'", newName, resp.Data.Name)
	}
}

func TestBoardColumn_Update_NotFound(t *testing.T) {
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
	newName := "Updated Column"
	code, _ := do[domain.BoardColumnModel](t, "PATCH", "/boards/"+uuidToString(board.ID)+"/columns/"+invalidColumnID, domain.BoardColumnUpdateModel{
		Name: newName,
	}, tokens.AccessToken)

	if code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", code)
	}
}
