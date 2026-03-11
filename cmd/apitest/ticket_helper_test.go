package apitest_test

import (
	"net/http"
	"testing"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

func createTicket(tb testing.TB, projectID string, token string, title, ticketType, priority string) domain.TicketModel {
	statusCode, resp := do[domain.TicketModel](tb, "POST", "/tickets?projectId="+projectID, domain.TicketCreateModel{
		Title:    title,
		Type:     ticketType,
		Priority: priority,
	}, token)

	if statusCode != http.StatusCreated {
		tb.Fatalf("create ticket failed: got status %d, error: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		tb.Fatalf("create ticket returned nil data")
	}

	return *resp.Data
}

func getTicket(tb testing.TB, ticketID string, token string) domain.TicketModel {
	statusCode, resp := do[domain.TicketModel](tb, "GET", "/tickets/"+ticketID, nil, token)

	if statusCode != http.StatusOK {
		tb.Fatalf("get ticket failed: got status %d, error: %v", statusCode, resp.Error)
	}

	if resp.Data == nil {
		tb.Fatalf("get ticket returned nil data")
	}

	return *resp.Data
}

func randomTicketTitle() string {
	return "Ticket " + randomString(8)
}
