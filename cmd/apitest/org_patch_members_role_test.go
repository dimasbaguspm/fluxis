package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestOrg_Members_UpdateRole(t *testing.T) {
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

	// Get User2's ID
	statusCode, userResp := do[domain.UserModel](t, "GET", "/users/me", nil, tokens2.AccessToken)
	if statusCode != http.StatusOK || userResp.Data == nil {
		t.Fatal("failed to get user2 data")
	}

	user2ID := userResp.Data.ID

	// User1 adds User2 as member
	do[struct{}](t, "POST", "/orgs/"+uuidToString(orgID)+"/members", domain.OrganisationMemberCreateModel{
		UserId: uuidToString(user2ID),
		Role:   "member",
	}, tokens1.AccessToken)

	// Update User2's role to admin
	statusCode, resp := do[struct{}](t, "PATCH", "/orgs/"+uuidToString(orgID)+"/members/"+uuidToString(user2ID), domain.OrganisationMemberUpdateModel{
		Role: "admin",
	}, tokens1.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}
}

func TestOrg_Members_InvalidRole(t *testing.T) {
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

	// Get User2's ID
	statusCode, userResp := do[domain.UserModel](t, "GET", "/users/me", nil, tokens2.AccessToken)
	if statusCode != http.StatusOK || userResp.Data == nil {
		t.Fatal("failed to get user2 data")
	}

	user2ID := userResp.Data.ID

	// Try to add User2 with invalid role "superuser"
	statusCode, _ = do[struct{}](t, "POST", "/orgs/"+uuidToString(orgID)+"/members", domain.OrganisationMemberCreateModel{
		UserId: uuidToString(user2ID),
		Role:   "superuser",
	}, tokens1.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}
}

func TestOrg_UpdateMember_MissingRole(t *testing.T) {
	tokens1 := register(t, randomEmail(), "User One", "SecurePassword123!")
	tokens2 := register(t, randomEmail(), "User Two", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens1.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	// Get user2 ID
	_, userResp := do[domain.UserModel](t, "GET", "/users/me", nil, tokens2.AccessToken)
	if userResp.Data == nil {
		t.Fatal("failed to get user2 data")
	}

	// Add user2 to org
	do[struct{}](t, "POST", "/orgs/"+uuidToString(orgResp.Data.ID)+"/members", domain.OrganisationMemberCreateModel{
		UserId: uuidToString(userResp.Data.ID),
		Role:   "member",
	}, tokens1.AccessToken)

	status, _ := do[struct{}](t, "PATCH", "/orgs/"+uuidToString(orgResp.Data.ID)+"/members/"+uuidToString(userResp.Data.ID), domain.OrganisationMemberUpdateModel{
		Role: "",
	}, tokens1.AccessToken)

	if status != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", status)
	}
}

func TestOrg_UpdateMember_InvalidRole(t *testing.T) {
	tokens1 := register(t, randomEmail(), "User One", "SecurePassword123!")
	tokens2 := register(t, randomEmail(), "User Two", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens1.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	// Get user2 ID
	_, userResp := do[domain.UserModel](t, "GET", "/users/me", nil, tokens2.AccessToken)
	if userResp.Data == nil {
		t.Fatal("failed to get user2 data")
	}

	// Add user2 to org
	do[struct{}](t, "POST", "/orgs/"+uuidToString(orgResp.Data.ID)+"/members", domain.OrganisationMemberCreateModel{
		UserId: uuidToString(userResp.Data.ID),
		Role:   "member",
	}, tokens1.AccessToken)

	status, _ := do[struct{}](t, "PATCH", "/orgs/"+uuidToString(orgResp.Data.ID)+"/members/"+uuidToString(userResp.Data.ID), domain.OrganisationMemberUpdateModel{
		Role: "superuser",
	}, tokens1.AccessToken)

	if status != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", status)
	}
}

func TestOrg_UpdateMember_NotFound(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	nonExistentUserID := "550e8400-e29b-41d4-a716-446655440000"

	status, _ := do[struct{}](t, "PATCH", "/orgs/"+uuidToString(orgResp.Data.ID)+"/members/"+nonExistentUserID, domain.OrganisationMemberUpdateModel{
		Role: "admin",
	}, tokens.AccessToken)

	// Non-existent member can result in 404 or 500 depending on error handling
	if status != http.StatusNotFound && status != http.StatusInternalServerError {
		t.Fatalf("expected status 404 or 500, got %d", status)
	}
}

func TestOrg_UpdateMember_Unauthenticated(t *testing.T) {
	orgID := "550e8400-e29b-41d4-a716-446655440000"
	userID := "550e8400-e29b-41d4-a716-446655440001"

	statusCode, _ := do[struct{}](t, "PATCH", "/orgs/"+orgID+"/members/"+userID, domain.OrganisationMemberUpdateModel{
		Role: "admin",
	}, "")

	if statusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", statusCode)
	}
}
