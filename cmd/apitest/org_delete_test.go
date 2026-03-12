package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestOrg_Delete_Success(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	// Create an org
	statusCode, createResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Organization " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || createResp.Data == nil {
		t.Fatalf("failed to create org: status=%d, data=%v, error=%v", statusCode, createResp.Data, createResp.Error)
	}

	orgID := createResp.Data.ID

	// Delete the org
	statusCode, _ = do[struct{}](t, "DELETE", "/orgs/"+uuidToString(orgID), nil, tokens.AccessToken)

	if statusCode != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", statusCode)
	}
}
