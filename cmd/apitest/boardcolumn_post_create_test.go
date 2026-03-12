package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestBoardColumn_Create_Success(t *testing.T) {
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

	columnName := randomBoardColumnName()
	position := int32(0)
	statusCode, resp := do[domain.BoardColumnModel](t, "POST", "/boards/"+uuidToString(board.ID)+"/columns", domain.BoardColumnCreateModel{
		Name:     columnName,
		Position: position,
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected column data")
	}

	if resp.Data.Name != columnName {
		t.Fatalf("expected name '%s', got '%s'", columnName, resp.Data.Name)
	}

	if resp.Data.Position != position {
		t.Fatalf("expected position %d, got %d", position, resp.Data.Position)
	}

	if uuidToString(resp.Data.ID) == "" {
		t.Fatal("expected non-empty ID")
	}
}

func TestBoardColumn_Create_MissingName(t *testing.T) {
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

	status, _ := do[domain.BoardColumnModel](t, "POST", "/boards/"+uuidToString(board.ID)+"/columns", domain.BoardColumnCreateModel{
		Name:     "",
		Position: 0,
	}, tokens.AccessToken)

	if status != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", status)
	}
}

func TestBoardColumn_Create_InvalidBoardUUID(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, _ := do[domain.BoardColumnModel](t, "POST", "/boards/not-a-uuid/columns", domain.BoardColumnCreateModel{
		Name:     "Test Column",
		Position: 0,
	}, tokens.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}
}

func TestBoardColumn_Create_NonExistentBoard(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	nonExistentBoardID := "550e8400-e29b-41d4-a716-446655440000"

	statusCode, _ := do[domain.BoardColumnModel](t, "POST", "/boards/"+nonExistentBoardID+"/columns", domain.BoardColumnCreateModel{
		Name:     "Test Column",
		Position: 0,
	}, tokens.AccessToken)

	if statusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", statusCode)
	}
}

func TestBoardColumn_Create_Unauthenticated(t *testing.T) {
	boardID := "550e8400-e29b-41d4-a716-446655440000"

	statusCode, _ := do[domain.BoardColumnModel](t, "POST", "/boards/"+boardID+"/columns", domain.BoardColumnCreateModel{
		Name:     "Test Column",
		Position: 0,
	}, "")

	if statusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", statusCode)
	}
}
