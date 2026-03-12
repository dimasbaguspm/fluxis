package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func TestTicket_List_Success(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := uuidToString(orgResp.Data.ID)
	project := createProject(t, orgID, tokens.AccessToken, randomProjectKey(), "Test Project "+randomString(8), "private")
	projectID := uuidToString(project.ID)

	// Create multiple tickets
	createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "story", "medium")
	createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "task", "low")
	createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "epic", "high")

	// List tickets with pagination
	statusCode, resp := do[domain.TicketsPagedModel](t, "GET", "/tickets?projectId="+projectID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		t.Fatal("expected tickets data")
	}

	if len(resp.Data.Items) != 3 {
		t.Fatalf("expected 3 tickets, got %d", len(resp.Data.Items))
	}

	if resp.Data.TotalCount != 3 {
		t.Fatalf("expected total count 3, got %d", resp.Data.TotalCount)
	}

	if resp.Data.PageNumber != 1 {
		t.Fatalf("expected page number 1, got %d", resp.Data.PageNumber)
	}
}

func TestTicket_List_EmptyProject(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := uuidToString(orgResp.Data.ID)
	project := createProject(t, orgID, tokens.AccessToken, randomProjectKey(), "Test Project", "private")
	projectID := uuidToString(project.ID)

	statusCode, resp := do[domain.TicketsPagedModel](t, "GET", "/tickets?projectId="+projectID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if resp.Data == nil {
		t.Fatal("expected tickets data")
	}

	if len(resp.Data.Items) != 0 {
		t.Fatalf("expected 0 tickets, got %d", len(resp.Data.Items))
	}

	if resp.Data.TotalCount != 0 {
		t.Fatalf("expected total count 0, got %d", resp.Data.TotalCount)
	}
}

func TestTicket_List_WithSprint(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := uuidToString(orgResp.Data.ID)
	project := createProject(t, orgID, tokens.AccessToken, randomProjectKey(), "Test Project", "private")
	projectID := uuidToString(project.ID)
	sprint := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprintID := uuidToString(sprint.ID)

	// Create tickets with and without sprint
	createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "story", "medium")
	t1 := createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "task", "low")
	t2 := createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "bug", "high")

	// Move tickets to sprint
	ticketID1 := uuidToString(t1.ID)
	ticketID2 := uuidToString(t2.ID)
	do[domain.TicketModel](t, "PATCH", "/tickets/"+ticketID1+"/move-to-sprint?sprintId="+sprintID, nil, tokens.AccessToken)
	do[domain.TicketModel](t, "PATCH", "/tickets/"+ticketID2+"/move-to-sprint?sprintId="+sprintID, nil, tokens.AccessToken)

	// List tickets for sprint
	statusCode, resp := do[domain.TicketsPagedModel](t, "GET", "/tickets?projectId="+projectID+"&sprintId="+sprintID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if resp.Data == nil {
		t.Fatal("expected tickets data")
	}

	if len(resp.Data.Items) != 2 {
		t.Fatalf("expected 2 tickets in sprint, got %d", len(resp.Data.Items))
	}

	if resp.Data.TotalCount != 2 {
		t.Fatalf("expected total count 2, got %d", resp.Data.TotalCount)
	}
}

func TestTicket_List_WithPagination(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := uuidToString(orgResp.Data.ID)
	project := createProject(t, orgID, tokens.AccessToken, randomProjectKey(), "Test Project", "private")
	projectID := uuidToString(project.ID)

	// Create 5 tickets
	for i := 0; i < 5; i++ {
		createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "story", "medium")
	}

	// Get first page with size 2
	statusCode, resp := do[domain.TicketsPagedModel](t, "GET", "/tickets?projectId="+projectID+"&pageNumber=1&pageSize=2", nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if len(resp.Data.Items) != 2 {
		t.Fatalf("expected 2 tickets on page 1, got %d", len(resp.Data.Items))
	}

	if resp.Data.TotalCount != 5 {
		t.Fatalf("expected total count 5, got %d", resp.Data.TotalCount)
	}

	if resp.Data.TotalPages != 3 {
		t.Fatalf("expected 3 total pages, got %d", resp.Data.TotalPages)
	}

	// Get second page
	statusCode, resp2 := do[domain.TicketsPagedModel](t, "GET", "/tickets?projectId="+projectID+"&pageNumber=2&pageSize=2", nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if len(resp2.Data.Items) != 2 {
		t.Fatalf("expected 2 tickets on page 2, got %d", len(resp2.Data.Items))
	}

	if resp2.Data.PageNumber != 2 {
		t.Fatalf("expected page number 2, got %d", resp2.Data.PageNumber)
	}

	// Get third page (should have 1 ticket)
	statusCode, resp3 := do[domain.TicketsPagedModel](t, "GET", "/tickets?projectId="+projectID+"&pageNumber=3&pageSize=2", nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if len(resp3.Data.Items) != 1 {
		t.Fatalf("expected 1 ticket on page 3, got %d", len(resp3.Data.Items))
	}
}

func TestTicket_List_InvalidProjectId(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, _ := do[domain.TicketsPagedModel](t, "GET", "/tickets?projectId=invalid-id", nil, tokens.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}
}

func TestTicket_List_MissingProjectId(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, _ := do[domain.TicketsPagedModel](t, "GET", "/tickets", nil, tokens.AccessToken)

	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", statusCode)
	}
}

func TestTicket_List_Unauthenticated(t *testing.T) {
	projectID := "550e8400-e29b-41d4-a716-446655440000"
	statusCode, _ := do[domain.TicketsPagedModel](t, "GET", "/tickets?projectId="+projectID, nil, "")

	if statusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", statusCode)
	}
}

// ===== Array Filtering Tests =====

func TestTicket_List_FilterByMultipleTicketIds(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := uuidToString(orgResp.Data.ID)
	project := createProject(t, orgID, tokens.AccessToken, randomProjectKey(), "Test Project", "private")
	projectID := uuidToString(project.ID)

	// Create 5 tickets
	tickets := make([]domain.TicketModel, 5)
	for i := 0; i < 5; i++ {
		tickets[i] = createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "story", "medium")
	}

	// Filter by first 2 tickets
	ticketID1 := uuidToString(tickets[0].ID)
	ticketID2 := uuidToString(tickets[1].ID)
	statusCode, resp := do[domain.TicketsPagedModel](t, "GET", "/tickets?projectId="+projectID+"&id="+ticketID1+"&id="+ticketID2, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if len(resp.Data.Items) != 2 {
		t.Fatalf("expected 2 tickets, got %d", len(resp.Data.Items))
	}

	if resp.Data.TotalCount != 2 {
		t.Fatalf("expected total count 2, got %d", resp.Data.TotalCount)
	}
}

func TestTicket_List_FilterByMultipleSprints(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := uuidToString(orgResp.Data.ID)
	project := createProject(t, orgID, tokens.AccessToken, randomProjectKey(), "Test Project", "private")
	projectID := uuidToString(project.ID)

	// Create 2 sprints
	sprint1 := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprint2 := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprint1ID := uuidToString(sprint1.ID)
	sprint2ID := uuidToString(sprint2.ID)

	// Create tickets in sprint1
	t1 := createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "story", "medium")
	t2 := createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "task", "low")
	do[domain.TicketModel](t, "PATCH", "/tickets/"+uuidToString(t1.ID)+"/move-to-sprint?sprintId="+sprint1ID, nil, tokens.AccessToken)
	do[domain.TicketModel](t, "PATCH", "/tickets/"+uuidToString(t2.ID)+"/move-to-sprint?sprintId="+sprint1ID, nil, tokens.AccessToken)

	// Create tickets in sprint2
	t3 := createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "bug", "high")
	t4 := createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "epic", "medium")
	do[domain.TicketModel](t, "PATCH", "/tickets/"+uuidToString(t3.ID)+"/move-to-sprint?sprintId="+sprint2ID, nil, tokens.AccessToken)
	do[domain.TicketModel](t, "PATCH", "/tickets/"+uuidToString(t4.ID)+"/move-to-sprint?sprintId="+sprint2ID, nil, tokens.AccessToken)

	// Filter by both sprints
	statusCode, resp := do[domain.TicketsPagedModel](t, "GET", "/tickets?projectId="+projectID+"&sprintId="+sprint1ID+"&sprintId="+sprint2ID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if len(resp.Data.Items) != 4 {
		t.Fatalf("expected 4 tickets (2 per sprint), got %d", len(resp.Data.Items))
	}

	if resp.Data.TotalCount != 4 {
		t.Fatalf("expected total count 4, got %d", resp.Data.TotalCount)
	}
}

func TestTicket_List_CombinedFilters(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := uuidToString(orgResp.Data.ID)
	project := createProject(t, orgID, tokens.AccessToken, randomProjectKey(), "Test Project", "private")
	projectID := uuidToString(project.ID)
	sprint := createSprint(t, projectID, tokens.AccessToken, randomSprintName())
	sprintID := uuidToString(sprint.ID)

	// Create tickets
	t1 := createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "story", "medium")
	t2 := createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "task", "low")
	createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "bug", "high") // not in sprint

	// Move first 2 to sprint
	do[domain.TicketModel](t, "PATCH", "/tickets/"+uuidToString(t1.ID)+"/move-to-sprint?sprintId="+sprintID, nil, tokens.AccessToken)
	do[domain.TicketModel](t, "PATCH", "/tickets/"+uuidToString(t2.ID)+"/move-to-sprint?sprintId="+sprintID, nil, tokens.AccessToken)

	// Filter: sprint AND specific ticket IDs
	statusCode, resp := do[domain.TicketsPagedModel](t, "GET", "/tickets?projectId="+projectID+"&sprintId="+sprintID+"&id="+uuidToString(t1.ID), nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if len(resp.Data.Items) != 1 {
		t.Fatalf("expected 1 ticket (intersection of filters), got %d", len(resp.Data.Items))
	}
}

func TestTicket_List_FilterByMultipleTicketIds_WithPagination(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := uuidToString(orgResp.Data.ID)
	project := createProject(t, orgID, tokens.AccessToken, randomProjectKey(), "Test Project", "private")
	projectID := uuidToString(project.ID)

	// Create 6 tickets
	var ticketIDs []string
	for i := 0; i < 6; i++ {
		ticket := createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "story", "medium")
		ticketIDs = append(ticketIDs, uuidToString(ticket.ID))
	}

	// Filter by first 4 tickets with pagination (pageSize=2)
	query := "/tickets?projectId=" + projectID + "&pageSize=2"
	for i := 0; i < 4; i++ {
		query += "&id=" + ticketIDs[i]
	}

	statusCode, resp := do[domain.TicketsPagedModel](t, "GET", query, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if len(resp.Data.Items) != 2 {
		t.Fatalf("expected 2 tickets on page 1, got %d", len(resp.Data.Items))
	}

	if resp.Data.TotalCount != 4 {
		t.Fatalf("expected total count 4, got %d", resp.Data.TotalCount)
	}

	if resp.Data.TotalPages != 2 {
		t.Fatalf("expected 2 total pages, got %d", resp.Data.TotalPages)
	}
}

func TestTicket_List_NonExistentTicketId_InArray(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := uuidToString(orgResp.Data.ID)
	project := createProject(t, orgID, tokens.AccessToken, randomProjectKey(), "Test Project", "private")
	projectID := uuidToString(project.ID)

	// Create 1 ticket
	ticket := createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "story", "medium")
	ticketID := uuidToString(ticket.ID)

	// Non-existent ticket ID
	nonExistentID := "550e8400-e29b-41d4-a716-446655440000"

	// Filter with both existing and non-existing
	statusCode, resp := do[domain.TicketsPagedModel](t, "GET", "/tickets?projectId="+projectID+"&id="+ticketID+"&id="+nonExistentID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	// Should return only the existing ticket
	if len(resp.Data.Items) != 1 {
		t.Fatalf("expected 1 ticket, got %d", len(resp.Data.Items))
	}
}

func TestTicket_List_EmptyIdArray(t *testing.T) {
	tokens := register(t, randomEmail(), "Test User", "SecurePassword123!")

	statusCode, orgResp := do[domain.OrganisationModel](t, "POST", "/orgs", domain.OrganisationCreateModel{
		Name: "Test Org " + randomString(8),
	}, tokens.AccessToken)

	if statusCode != http.StatusCreated || orgResp.Data == nil {
		t.Fatalf("failed to create org")
	}

	orgID := uuidToString(orgResp.Data.ID)
	project := createProject(t, orgID, tokens.AccessToken, randomProjectKey(), "Test Project", "private")
	projectID := uuidToString(project.ID)

	// Create tickets
	for i := 0; i < 3; i++ {
		createTicket(t, projectID, tokens.AccessToken, randomTicketTitle(), "story", "medium")
	}

	// No id filter - should return all tickets for projectId
	statusCode, resp := do[domain.TicketsPagedModel](t, "GET", "/tickets?projectId="+projectID, nil, tokens.AccessToken)

	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}

	if len(resp.Data.Items) != 3 {
		t.Fatalf("expected 3 tickets (no id filter), got %d", len(resp.Data.Items))
	}
}
