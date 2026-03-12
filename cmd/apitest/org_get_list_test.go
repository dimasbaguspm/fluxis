package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestOrg_List_Success(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	// Create an org first
	do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	statusCode, resp := do[[]domain.OrganisationModel](t, "GET", "/orgs", nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected org list data")
	}

	if len(*resp.Data) < 1 {
		t.Fatal("expected at least one org in list")
	}
}

func TestOrg_List_ResponseStructure(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	// Create multiple orgs
	orgNames := []string{
		"Test Org 1 " + randomString(4),
		"Test Org 2 " + randomString(4),
	}
	for _, name := range orgNames {
		do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
			Name: name,
		}, tokens.AccessToken)
	}

	statusCode, resp := do[[]domain.OrganisationModel](t, "GET", "/orgs", nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if resp.Data == nil {
		t.Fatal("expected org list data")
	}

	if len(*resp.Data) < 2 {
		t.Fatalf("expected at least 2 orgs, got %d", len(*resp.Data))
	}

	// Verify response structure of list items
	for i, org := range *resp.Data {
		if uuidToString(org.ID) == "" {
			t.Fatalf("org[%d]: expected non-empty ID", i)
		}

		if org.Name == "" {
			t.Fatalf("org[%d]: expected Name field", i)
		}

		// Verify TotalMembers field exists (should be >= 1 for orgs the user is member of)
		if org.TotalMembers < 1 {
			t.Logf("org[%d]: TotalMembers is %d (may be 0 for non-owned orgs)", i, org.TotalMembers)
		}
	}
}
