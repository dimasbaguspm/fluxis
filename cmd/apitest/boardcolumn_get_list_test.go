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
	col1 := createBoardColumn(t, uuidToString(board.ID), tokens.AccessToken, "Column 1", 0)
	col2 := createBoardColumn(t, uuidToString(board.ID), tokens.AccessToken, "Column 2", 1)

	statusCode, resp := do[[]domain.BoardColumnModel](t, "GET", "/boards/"+uuidToString(board.ID)+"/columns", nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected columns data")
	}

	if len(*resp.Data) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(*resp.Data))
	}

	// Verify columns are in order
	if (*resp.Data)[0].ID != col1.ID {
		t.Fatalf("expected first column to be col1")
	}
	if (*resp.Data)[1].ID != col2.ID {
		t.Fatalf("expected second column to be col2")
	}
}
