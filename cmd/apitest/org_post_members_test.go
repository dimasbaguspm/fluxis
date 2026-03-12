package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestOrg_Members_AddAndList(t *testing.T) {
	// Create two users
	tokens1 := register(t, randomEmail(), "User One", "SecurePassword123!")
	email2 := randomEmail()
	tokens2 := register(t, email2, "User Two", "SecurePassword123!")

	// User1 creates an org
	statusCode, createResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Organization " + randomString(8),
	}, tokens1.AccessToken)

	if statusCode != http.StatusCreated || createResp.Data == nil {
		t.Fatalf("failed to create org: status=%d, data=%v, error=%v", statusCode, createResp.Data, createResp.Error)
	}

	orgID := createResp.Data.ID

	// Get User2's ID from /users/me
	statusCode, userResp := do[domain.UserModel](t, "GET", "/users/me", nil, tokens2.AccessToken)
	if statusCode != http.StatusOK || userResp.Data == nil {
		t.Fatal("failed to get user2 data")
	}

	user2ID := userResp.Data.ID

	// User1 adds User2 to the org
	statusCode, resp := do[struct{}](t, "POST", "/orgs/"+uuidToString(orgID)+"/members", domain.OrganisationMemberCreateModel{
		UserId: uuidToString(user2ID),
		Role:   "member",
	}, tokens1.AccessToken)

	if statusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %v", statusCode, resp.Error)
	}

	// List org members
	statusCode, listResp := do[[]domain.OrganisationMemberModel](t, "GET", "/orgs/"+uuidToString(orgID)+"/members", nil, tokens1.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, listResp.Error)
	}

	if listResp.Data == nil {
		t.Fatalf("expected member list, got nil data, error=%v", listResp.Error)
	}

	if len(*listResp.Data) < 2 {
		t.Fatalf("expected at least 2 members, got %d: %v", len(*listResp.Data), listResp.Data)
	}
}
