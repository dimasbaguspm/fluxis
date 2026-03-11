package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func createProject(tb testing.TB, orgID string, token string, key, name, visibility string) domain.ProjectModel {
	statusCode, resp := do[domain.ProjectModel](tb, "POST", "/projects?orgId="+orgID, domain.ProjectCreateModel{
		Key:        key,
		Name:       name,
		Visibility: visibility,
	}, token)

	if statusCode != http.StatusCreated {
		tb.Fatalf("create project failed: got status %d, error: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		tb.Fatalf("create project returned nil data")
	}

	return *resp.Data
}

func randomProjectKey() string {
	return "p" + randomString(4)
}
