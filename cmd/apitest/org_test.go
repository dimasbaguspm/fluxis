package apitest_test

import (
	"net/http"
	"testing"
)

func TestOrgs_Create_Success(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, resp := do[orgModel](t, "POST", "/orgs", map[string]string{
		"name": "Test Organization",
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected org data")
	}

	if resp.Data.Name != "Test Organization" {
		t.Fatalf("expected name 'Test Organization', got %s", resp.Data.Name)
	}

	if resp.Data.ID == "" {
		t.Fatal("expected non-empty ID")
	}

	if resp.Data.TotalMembers != 1 {
		t.Fatalf("expected totalMembers=1, got %d", resp.Data.TotalMembers)
	}
}

func TestOrgs_Create_Unauthenticated(t *testing.T) {
	statusCode, _ := do[orgModel](t, "POST", "/orgs", map[string]string{
		"name": "Test Organization",
	}, "")

	if statusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", statusCode)
	}
}

func TestOrgs_Create_MissingName(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, _ := do[orgModel](t, "POST", "/orgs", map[string]interface{}{}, tokens.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}
}

func TestOrgs_List(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	// Create an org first
	do[orgModel](t, "POST", "/orgs", map[string]string{
		"name": "Test Org 1",
	}, tokens.AccessToken)

	statusCode, resp := do[[]orgModel](t, "GET", "/orgs", nil, tokens.AccessToken)

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
	statusCode, createResp := do[orgModel](t, "POST", "/orgs", map[string]string{
		"name": "Test Organization",
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || createResp.Data == nil {
		t.Fatal("failed to create org")
	}

	orgID := createResp.Data.ID

	// Get the org
	statusCode, resp := do[orgModel](t, "GET", "/orgs/"+orgID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil || resp.Data.ID != orgID {
		t.Fatalf("expected org ID %s, got %v", orgID, resp.Data)
	}
}

func TestOrgs_GetByID_NotFound(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	nonExistentID := "550e8400-e29b-41d4-a716-446655440000"

	statusCode, _ := do[orgModel](t, "GET", "/orgs/"+nonExistentID, nil, tokens.AccessToken)

	if statusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", statusCode)
	}
}

func TestOrgs_GetByID_InvalidUUID(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, _ := do[orgModel](t, "GET", "/orgs/not-a-uuid", nil, tokens.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}
}

func TestOrgs_Update_Success(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	// Create an org
	statusCode, createResp := do[orgModel](t, "POST", "/orgs", map[string]string{
		"name": "Original Name",
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || createResp.Data == nil {
		t.Fatal("failed to create org")
	}

	orgID := createResp.Data.ID

	// Update the org
	statusCode, resp := do[orgModel](t, "PATCH", "/orgs/"+orgID, map[string]string{
		"name": "Updated Name",
	}, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil || resp.Data.Name != "Updated Name" {
		t.Fatalf("expected updated name, got %v", resp.Data)
	}
}

func TestOrgs_Delete_Success(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	// Create an org
	statusCode, createResp := do[orgModel](t, "POST", "/orgs", map[string]string{
		"name": "Test Organization",
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || createResp.Data == nil {
		t.Fatal("failed to create org")
	}

	orgID := createResp.Data.ID

	// Delete the org
	statusCode, resp := do[struct{}](t, "DELETE", "/orgs/"+orgID, nil, tokens.AccessToken)

	if statusCode != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d: %v", statusCode, resp.Error)
	}
}

func TestOrgs_Members_AddAndList(t *testing.T) {
	// Create two users
	tokens1 := register(t, randomEmail(), "User One", "SecurePassword123!")
	email2 := randomEmail()
	tokens2 := register(t, email2, "User Two", "SecurePassword123!")

	// User1 creates an org
	statusCode, createResp := do[orgModel](t, "POST", "/orgs", map[string]string{
		"name": "Test Organization",
	}, tokens1.AccessToken)

	if statusCode != http.StatusCreated || createResp.Data == nil {
		t.Fatal("failed to create org")
	}

	orgID := createResp.Data.ID

	// Get User2's ID from /users/me
	statusCode, userResp := do[userModel](t, "GET", "/users/me", nil, tokens2.AccessToken)
	if statusCode != http.StatusOK || userResp.Data == nil {
		t.Fatal("failed to get user2 data")
	}

	user2ID := userResp.Data.ID

	// User1 adds User2 to the org
	statusCode, resp := do[struct{}](t, "POST", "/orgs/"+orgID+"/members", map[string]string{
		"userId": user2ID,
		"role":   "member",
	}, tokens1.AccessToken)

	if statusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %v", statusCode, resp.Error)
	}

	// List org members
	statusCode, listResp := do[[]orgMemberModel](t, "GET", "/orgs/"+orgID+"/members", nil, tokens1.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, listResp.Error)
	}

	if listResp.Data == nil || len(*listResp.Data) < 2 {
		t.Fatalf("expected at least 2 members, got %v", listResp.Data)
	}
}

func TestOrgs_Members_UpdateRole(t *testing.T) {
	// Create two users
	tokens1 := register(t, randomEmail(), "User One", "SecurePassword123!")
	email2 := randomEmail()
	tokens2 := register(t, email2, "User Two", "SecurePassword123!")

	// User1 creates an org
	statusCode, createResp := do[orgModel](t, "POST", "/orgs", map[string]string{
		"name": "Test Organization",
	}, tokens1.AccessToken)

	if statusCode != http.StatusCreated || createResp.Data == nil {
		t.Fatal("failed to create org")
	}

	orgID := createResp.Data.ID

	// Get User2's ID
	statusCode, userResp := do[userModel](t, "GET", "/users/me", nil, tokens2.AccessToken)
	if statusCode != http.StatusOK || userResp.Data == nil {
		t.Fatal("failed to get user2 data")
	}

	user2ID := userResp.Data.ID

	// User1 adds User2 as member
	do[struct{}](t, "POST", "/orgs/"+orgID+"/members", map[string]string{
		"userId": user2ID,
		"role":   "member",
	}, tokens1.AccessToken)

	// Update User2's role to admin
	statusCode, resp := do[struct{}](t, "PATCH", "/orgs/"+orgID+"/members/"+user2ID, map[string]string{
		"role": "admin",
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
	statusCode, createResp := do[orgModel](t, "POST", "/orgs", map[string]string{
		"name": "Test Organization",
	}, tokens1.AccessToken)

	if statusCode != http.StatusCreated || createResp.Data == nil {
		t.Fatal("failed to create org")
	}

	orgID := createResp.Data.ID

	// Get User2's ID
	statusCode, userResp := do[userModel](t, "GET", "/users/me", nil, tokens2.AccessToken)
	if statusCode != http.StatusOK || userResp.Data == nil {
		t.Fatal("failed to get user2 data")
	}

	user2ID := userResp.Data.ID

	// User1 adds User2
	do[struct{}](t, "POST", "/orgs/"+orgID+"/members", map[string]string{
		"userId": user2ID,
		"role":   "member",
	}, tokens1.AccessToken)

	// User1 removes User2
	statusCode, resp := do[struct{}](t, "DELETE", "/orgs/"+orgID+"/members/"+user2ID, nil, tokens1.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}
}

func TestOrgs_Members_InvalidRole(t *testing.T) {
	// Create two users
	tokens1 := register(t, randomEmail(), "User One", "SecurePassword123!")
	email2 := randomEmail()
	tokens2 := register(t, email2, "User Two", "SecurePassword123!")

	// User1 creates an org
	statusCode, createResp := do[orgModel](t, "POST", "/orgs", map[string]string{
		"name": "Test Organization",
	}, tokens1.AccessToken)

	if statusCode != http.StatusCreated || createResp.Data == nil {
		t.Fatal("failed to create org")
	}

	orgID := createResp.Data.ID

	// Get User2's ID
	statusCode, userResp := do[userModel](t, "GET", "/users/me", nil, tokens2.AccessToken)
	if statusCode != http.StatusOK || userResp.Data == nil {
		t.Fatal("failed to get user2 data")
	}

	user2ID := userResp.Data.ID

	// Try to add User2 with invalid role "superuser"
	statusCode, _ = do[struct{}](t, "POST", "/orgs/"+orgID+"/members", map[string]string{
		"userId": user2ID,
		"role":   "superuser",
	}, tokens1.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}
}
