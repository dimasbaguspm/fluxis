package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func createBoard(tb testing.TB, sprintID string, token string, name string) domain.BoardModel {
	statusCode, resp := do[domain.BoardModel](tb, "POST", "/boards?sprintId="+sprintID, domain.BoardCreateModel{
		Name: &name,
	}, token)

	if statusCode != http.StatusCreated {
		tb.Fatalf("create board failed: got status %d, error: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		tb.Fatalf("create board returned nil data")
	}

	return *resp.Data
}

func randomBoardName() string {
	return "Board " + randomString(4)
}

func createBoardColumn(tb testing.TB, boardID string, token string, name string, position int32) domain.BoardColumnModel {
	statusCode, resp := do[domain.BoardColumnModel](tb, "POST", "/boards/"+boardID+"/columns", domain.BoardColumnCreateModel{
		Name:     &name,
		Position: &position,
	}, token)

	if statusCode != http.StatusCreated {
		tb.Fatalf("create board column failed: got status %d, error: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		tb.Fatalf("create board column returned nil data")
	}

	return *resp.Data
}

func randomBoardColumnName() string {
	return "Column " + randomString(4)
}
