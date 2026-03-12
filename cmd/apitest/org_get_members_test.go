package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestOrg_ListMembers_Success(t *testing.T) {
	// Create two users
	tokens1 := register(t, randomEmail(), "User One", "SecurePassword123!")
	email2 := randomEmail()
	tokens2 := register(t, email2, "User Two", "SecurePassword123!")

	// User1 creates an org
	statusCode, createResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Organization " + randomString(8),
	}, tokens1.AccessToken)

	if statusCode != http.StatusCreated || createResp.Data == nil {
		t.Fatalf("failed to create org: status=%d", statusCode)
	}

	orgID := createResp.Data.ID

	// Get User2's ID
	statusCode, userResp := do[domain.UserModel](t, "GET", "/users/me", nil, tokens2.AccessToken)
	if statusCode != http.StatusOK || userResp.Data == nil {
		t.Fatal("failed to get user2 data")
	}

	user2ID := userResp.Data.ID

	// User1 adds User2 to the org
	do[struct{}](t, "POST", "/orgs/"+uuidToString(orgID)+"/members", domain.OrganisationMemberCreateModel{
		UserId: uuidToString(user2ID),
		Role:   "member",
	}, tokens1.AccessToken)

	// List org members
	statusCode, resp := do[domain.OrganisationMembersPagedModel](t, "GET", "/orgs/"+uuidToString(orgID)+"/members", nil, tokens1.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected member list data")
	}

	if len(resp.Data.Items) < 2 {
		t.Fatalf("expected at least 2 members, got %d", len(resp.Data.Items))
	}

	if resp.Data.TotalCount < 2 {
		t.Fatalf("expected totalCount >= 2, got %d", resp.Data.TotalCount)
	}

	if resp.Data.PageNumber != 1 {
		t.Fatalf("expected page 1, got %d", resp.Data.PageNumber)
	}

	if resp.Data.PageSize != 25 {
		t.Fatalf("expected pageSize 25, got %d", resp.Data.PageSize)
	}
}

func TestOrg_ListMembers_EmptyOrg(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	// Create an org
	statusCode, createResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || createResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := createResp.Data.ID

	// List members - should return at least the creator
	statusCode, resp := do[domain.OrganisationMembersPagedModel](t, "GET", "/orgs/"+uuidToString(orgID)+"/members", nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if resp.Data == nil {
		t.Fatal("expected member list data")
	}

	if len(resp.Data.Items) < 1 {
		t.Fatal("expected at least one member (creator)")
	}

	if resp.Data.TotalCount < 1 {
		t.Fatal("expected totalCount >= 1")
	}
}

func TestOrg_ListMembers_NotFound(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	nonExistentID := "550e8400-e29b-41d4-a716-446655440000"

	statusCode, resp := do[domain.OrganisationMembersPagedModel](t, "GET", "/orgs/"+nonExistentID+"/members", nil, tokens.AccessToken)

	// API returns empty list for non-existent org instead of 404
	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if resp.Data == nil || len(resp.Data.Items) != 0 {
		t.Fatalf("expected empty member list for non-existent org")
	}

	if resp.Data.TotalCount != 0 {
		t.Fatalf("expected totalCount 0, got %d", resp.Data.TotalCount)
	}
}

func TestOrg_ListMembers_InvalidUUID(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, _ := do[domain.OrganisationMembersPagedModel](t, "GET", "/orgs/not-a-uuid/members", nil, tokens.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}
}

func TestOrg_ListMembers_Unauthenticated(t *testing.T) {
	orgID := "550e8400-e29b-41d4-a716-446655440000"

	statusCode, _ := do[domain.OrganisationMembersPagedModel](t, "GET", "/orgs/"+orgID+"/members", nil, "")

	if statusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", statusCode)
	}
}

func TestOrg_ListMembers_WithPagination(t *testing.T) {
	// Create user and org
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, createResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || createResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := createResp.Data.ID

	// List members with custom pageSize
	statusCode, resp := do[domain.OrganisationMembersPagedModel](t, "GET", "/orgs/"+uuidToString(orgID)+"/members?pageSize=10&pageNumber=1", nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if resp.Data == nil {
		t.Fatal("expected member list data")
	}

	if resp.Data.PageNumber != 1 {
		t.Fatalf("expected page 1, got %d", resp.Data.PageNumber)
	}

	if resp.Data.PageSize != 10 {
		t.Fatalf("expected pageSize 10, got %d", resp.Data.PageSize)
	}
}

func TestOrg_ListMembers_WithEmailFilter(t *testing.T) {
	// Create two users
	tokens1 := register(t, randomEmail(), "User One", "SecurePassword123!")
	email2 := randomEmail()
	tokens2 := register(t, email2, "User Two", "SecurePassword123!")

	// User1 creates an org
	statusCode, createResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Organization " + randomString(8),
	}, tokens1.AccessToken)

	if statusCode != http.StatusCreated || createResp.Data == nil {
		t.Fatalf("failed to create org: status=%d", statusCode)
	}

	orgID := createResp.Data.ID

	// Get User2's ID
	statusCode, userResp := do[domain.UserModel](t, "GET", "/users/me", nil, tokens2.AccessToken)
	if statusCode != http.StatusOK || userResp.Data == nil {
		t.Fatal("failed to get user2 data")
	}

	user2ID := userResp.Data.ID

	// User1 adds User2 to the org
	do[struct{}](t, "POST", "/orgs/"+uuidToString(orgID)+"/members", domain.OrganisationMemberCreateModel{
		UserId: uuidToString(user2ID),
		Role:   "member",
	}, tokens1.AccessToken)

	// List members filtering by email
	statusCode, resp := do[domain.OrganisationMembersPagedModel](t, "GET", "/orgs/"+uuidToString(orgID)+"/members?email="+email2, nil, tokens1.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if resp.Data == nil {
		t.Fatal("expected member list data")
	}

	if len(resp.Data.Items) != 1 {
		t.Fatalf("expected 1 member matching filter, got %d", len(resp.Data.Items))
	}

	if resp.Data.Items[0].Email != email2 {
		t.Fatalf("expected email %s, got %s", email2, resp.Data.Items[0].Email)
	}
}
