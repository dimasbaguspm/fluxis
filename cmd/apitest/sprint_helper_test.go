package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func createSprint(tb testing.TB, projectID string, token string, name string) domain.SprintModel {
	statusCode, resp := do[domain.SprintModel](tb, "POST", "/sprints?projectId="+projectID, domain.SprintCreateModel{
		Name: &name,
	}, token)

	if statusCode != http.StatusCreated {
		tb.Fatalf("create sprint failed: got status %d, error: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		tb.Fatalf("create sprint returned nil data")
	}

	return *resp.Data
}

func randomSprintName() string {
	return "Sprint " + randomString(4)
}
