package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestOrg_Update_Success(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	// Create an org
	statusCode, createResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Original Name " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || createResp.Data == nil {
		t.Fatalf("failed to create org: status=%d, data=%v, error=%v", statusCode, createResp.Data, createResp.Error)
	}

	orgID := createResp.Data.ID

	// Update the org
	updatedName := "Updated Name " + randomString(8)
	statusCode, resp := do[domain.OrganisationModel](t, "PATCH", "/orgs/"+uuidToString(orgID), domain.OrganisationUpdateModel{
		Name: updatedName,
	}, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil || resp.Data.Name != updatedName {
		t.Fatalf("expected name '%s', got %v", updatedName, resp.Data)
	}
}

func TestOrg_Update_MissingName(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, createResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Original Name " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || createResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := createResp.Data.ID

	status, _ := do[domain.OrganisationModel](t, "PATCH", "/orgs/"+uuidToString(orgID), domain.OrganisationUpdateModel{
		Name: "",
	}, tokens.AccessToken)

	if status != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", status)
	}
}

func TestOrg_Update_NotFound(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	nonExistentID := "550e8400-e29b-41d4-a716-446655440000"

	statusCode, _ := do[domain.OrganisationModel](t, "PATCH", "/orgs/"+nonExistentID, domain.OrganisationUpdateModel{
		Name: "Updated Name",
	}, tokens.AccessToken)

	if statusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", statusCode)
	}
}

func TestOrg_Update_InvalidUUID(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, _ := do[domain.OrganisationModel](t, "PATCH", "/orgs/not-a-uuid", domain.OrganisationUpdateModel{
		Name: "Updated Name",
	}, tokens.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}
}

func TestOrg_Update_Unauthenticated(t *testing.T) {
	nonExistentID := "550e8400-e29b-41d4-a716-446655440000"

	statusCode, _ := do[domain.OrganisationModel](t, "PATCH", "/orgs/"+nonExistentID, domain.OrganisationUpdateModel{
		Name: "Updated Name",
	}, "")

	if statusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", statusCode)
	}
}
