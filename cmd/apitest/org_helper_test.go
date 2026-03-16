package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func createOrg(tb testing.TB, token string) domain.OrganisationModel {
	statusCode, resp := do[domain.OrganisationModel](tb, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, token)

	if statusCode != http.StatusCreated {
		tb.Fatalf("create org failed: got status %d, error: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		tb.Fatalf("create org returned nil data")
	}

	return *resp.Data
}
