package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestOrgs_Create_Success(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	orgName := "Test Organization " + randomString(8)
	statusCode, resp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: orgName,
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected org data")
	}

	if resp.Data.Name != orgName {
		t.Fatalf("expected name '%s', got %s", orgName, resp.Data.Name)
	}

	if uuidToString(resp.Data.ID) == "" {
		t.Fatal("expected non-empty ID")
	}

	if resp.Data.TotalMembers != 1 {
		t.Fatalf("expected totalMembers=1, got %d", resp.Data.TotalMembers)
	}
}

func TestOrgs_Create_Unauthenticated(t *testing.T) {
	statusCode, _ := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Organization " + randomString(8),
	}, "")

	if statusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", statusCode)
	}
}

func TestOrgs_Create_MissingName(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, _ := do[domain.OrganisationModel](t, "POST", "/orgs", map[string]interface{}{}, tokens.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}
}

func TestOrgs_List(t *testing.T) {
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

func TestOrgs_GetByID_Success(t *testing.T) {
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

func TestOrgs_GetByID_NotFound(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	nonExistentID := "550e8400-e29b-41d4-a716-446655440000"

	statusCode, _ := do[domain.OrganisationModel](t, "GET", "/orgs/"+nonExistentID, nil, tokens.AccessToken)

	if statusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", statusCode)
	}
}

func TestOrgs_GetByID_InvalidUUID(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, _ := do[domain.OrganisationModel](t, "GET", "/orgs/not-a-uuid", nil, tokens.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}
}

func TestOrgs_Update_Success(t *testing.T) {
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

func TestOrgs_Delete_Success(t *testing.T) {
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

func TestOrgs_Members_AddAndList(t *testing.T) {
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

func TestOrgs_Members_UpdateRole(t *testing.T) {
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

func TestOrgs_Members_Delete(t *testing.T) {
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

	// User1 adds User2
	do[struct{}](t, "POST", "/orgs/"+uuidToString(orgID)+"/members", domain.OrganisationMemberCreateModel{
		UserId: uuidToString(user2ID),
		Role:   "member",
	}, tokens1.AccessToken)

	// User1 removes User2
	statusCode, _ = do[struct{}](t, "DELETE", "/orgs/"+uuidToString(orgID)+"/members/"+uuidToString(user2ID), nil, tokens1.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}
}

func TestOrgs_Members_InvalidRole(t *testing.T) {
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
