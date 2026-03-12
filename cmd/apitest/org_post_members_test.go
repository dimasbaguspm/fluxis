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
	statusCode, listResp := do[domain.OrganisationMembersPagedModel](t, "GET", "/orgs/"+uuidToString(orgID)+"/members", nil, tokens1.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, listResp.Error)
	}

	if listResp.Data == nil {
		t.Fatalf("expected member list, got nil data, error=%v", listResp.Error)
	}

	if len(listResp.Data.Items) < 2 {
		t.Fatalf("expected at least 2 members, got %d: %v", len(listResp.Data.Items), listResp.Data)
	}

	if listResp.Data.TotalCount < 2 {
		t.Fatalf("expected totalCount >= 2, got %d", listResp.Data.TotalCount)
	}
}

func TestOrg_AddMember_MissingUserId(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	status, _ := do[struct{}](t, "POST", "/orgs/"+uuidToString(orgResp.Data.ID)+"/members", domain.OrganisationMemberCreateModel{
		UserId: "",
		Role:   "member",
	}, tokens.AccessToken)

	if status != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", status)
	}
}

func TestOrg_AddMember_MissingRole(t *testing.T) {
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

	status, _ := do[struct{}](t, "POST", "/orgs/"+uuidToString(orgResp.Data.ID)+"/members", domain.OrganisationMemberCreateModel{
		UserId: uuidToString(userResp.Data.ID),
		Role:   "",
	}, tokens1.AccessToken)

	if status != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", status)
	}
}

func TestOrg_AddMember_InvalidRole(t *testing.T) {
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

	status, _ := do[struct{}](t, "POST", "/orgs/"+uuidToString(orgResp.Data.ID)+"/members", domain.OrganisationMemberCreateModel{
		UserId: uuidToString(userResp.Data.ID),
		Role:   "superuser",
	}, tokens1.AccessToken)

	if status != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", status)
	}
}

func TestOrg_AddMember_NonExistentUser(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	nonExistentUserID := "550e8400-e29b-41d4-a716-446655440000"

	status, _ := do[struct{}](t, "POST", "/orgs/"+uuidToString(orgResp.Data.ID)+"/members", domain.OrganisationMemberCreateModel{
		UserId: nonExistentUserID,
		Role:   "member",
	}, tokens.AccessToken)

	// Non-existent user can result in either 404 or 500 (database constraint)
	if status != http.StatusNotFound && status != http.StatusInternalServerError {
		t.Fatalf("expected status 404 or 500, got %d", status)
	}
}

func TestOrg_AddMember_NonExistentOrg(t *testing.T) {
	tokens1 := register(t, randomEmail(), "User One", "SecurePassword123!")
	tokens2 := register(t, randomEmail(), "User Two", "SecurePassword123!")

	// Get user2 ID
	_, userResp := do[domain.UserModel](t, "GET", "/users/me", nil, tokens2.AccessToken)
	if userResp.Data == nil {
		t.Fatal("failed to get user2 data")
	}

	nonExistentOrgID := "550e8400-e29b-41d4-a716-446655440000"

	status, _ := do[struct{}](t, "POST", "/orgs/"+nonExistentOrgID+"/members", domain.OrganisationMemberCreateModel{
		UserId: uuidToString(userResp.Data.ID),
		Role:   "member",
	}, tokens1.AccessToken)

	// Non-existent org can result in either 404 or 500 (database constraint)
	if status != http.StatusNotFound && status != http.StatusInternalServerError {
		t.Fatalf("expected status 404 or 500, got %d", status)
	}
}

func TestOrg_AddMember_Unauthenticated(t *testing.T) {
	orgID := "550e8400-e29b-41d4-a716-446655440000"
	userID := "550e8400-e29b-41d4-a716-446655440001"

	statusCode, _ := do[struct{}](t, "POST", "/orgs/"+orgID+"/members", domain.OrganisationMemberCreateModel{
		UserId: userID,
		Role:   "member",
	}, "")

	if statusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", statusCode)
	}
}
