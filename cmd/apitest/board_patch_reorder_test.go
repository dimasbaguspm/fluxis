package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestBoard_Reorder_Success(t *testing.T) {
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

	statusCode, resp := do[[]domain.BoardModel](t, "PATCH", "/boards/reorder", domain.BoardReorderModel{
		Boards: []domain.BoardPositionUpdate{
			{
				ID:       board1.ID,
				Position: 1,
			},
			{
				ID:       board2.ID,
				Position: 0,
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
