package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestOrg_GetByID_Success(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	// Create an org
	statusCode, createResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Organization " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || createResp.Data == nil {
		t.Fatalf("failed to create org: status=%d, data=%v, error=%v", statusCode, createResp.Data, createResp.Error)
	}

	orgID := createResp.Data.ID

	// Get the org
	statusCode, resp := do[domain.OrganisationModel](t, "GET", "/orgs/"+uuidToString(orgID), nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil || resp.Data.ID != orgID {
		t.Fatalf("expected org ID %s, got %v", uuidToString(orgID), resp.Data)
	}
}

func TestOrg_GetByID_NotFound(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	nonExistentID := "550e8400-e29b-41d4-a716-446655440000"

	statusCode, _ := do[domain.OrganisationModel](t, "GET", "/orgs/"+nonExistentID, nil, tokens.AccessToken)

	if statusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", statusCode)
	}
}

func TestOrg_GetByID_InvalidUUID(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, _ := do[domain.OrganisationModel](t, "GET", "/orgs/not-a-uuid", nil, tokens.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}
}

func TestOrg_GetByID_ResponseStructure(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	orgName := "Test Org " + randomString(8)
	statusCode, createResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: orgName,
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || createResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := createResp.Data.ID

	// Get the org and verify response structure
	statusCode, resp := do[domain.OrganisationModel](t, "GET", "/orgs/"+uuidToString(orgID), nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if resp.Data == nil {
		t.Fatal("expected org data in response")
	}

	// Verify all expected fields are present
	if uuidToString(resp.Data.ID) == "" {
		t.Fatal("expected non-empty ID in response")
	}

	if resp.Data.Name == "" {
		t.Fatal("expected Name field in response")
	}

	if resp.Data.Name != orgName {
		t.Fatalf("expected name '%s', got '%s'", orgName, resp.Data.Name)
	}

	if resp.Data.TotalMembers == 0 {
		t.Fatal("expected TotalMembers field in response")
	}
}
