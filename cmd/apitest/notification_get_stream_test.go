package apitest_test

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/dimasbaguspm/fluxis/pkg/domain"
)

// TestNotification_Stream_Unauthenticated returns 401 when no token provided
func TestNotification_Stream_Unauthenticated(t *testing.T) {
	req, err := http.NewRequest("GET", testServer.URL+"/notifications", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", resp.StatusCode)
	}
}

// TestNotification_Stream_Success establishes a connection and receives events
func TestNotification_Stream_Success(t *testing.T) {
	// Register and login user
	auth := register(t, randomEmail(), "Test User", "SecurePassword123!")
	token := auth.AccessToken
	org := createOrg(t, token)

	// Start streaming in goroutine
	eventsCh := make(chan map[string]any, 10)
	errCh := make(chan error, 1)

	go func() {
		req, err := http.NewRequest("GET", testServer.URL+"/notifications", nil)
		if err != nil {
			errCh <- err
			return
		}
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			errCh <- err
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			errCh <- fmt.Errorf("expected status 200, got %d", resp.StatusCode)
			return
		}

		// Check headers
		contentType := resp.Header.Get("Content-Type")
		if contentType != "text/event-stream" {
			errCh <- fmt.Errorf("expected Content-Type: text/event-stream, got %s", contentType)
			return
		}

		cacheControl := resp.Header.Get("Cache-Control")
		if cacheControl != "no-cache" {
			errCh <- fmt.Errorf("expected Cache-Control: no-cache, got %s", cacheControl)
			return
		}

		// Read SSE events
		scanner := bufio.NewScanner(resp.Body)
		eventType := ""
		for scanner.Scan() {
			line := scanner.Text()

			if strings.HasPrefix(line, "event: ") {
				eventType = strings.TrimPrefix(line, "event: ")
			} else if strings.HasPrefix(line, "data: ") {
				dataStr := strings.TrimPrefix(line, "data: ")

				var eventData map[string]any
				if err := json.Unmarshal([]byte(dataStr), &eventData); err != nil {
					errCh <- fmt.Errorf("failed to unmarshal event data: %v", err)
					return
				}

				eventData["_type"] = eventType
				eventsCh <- eventData
			}

			// Stop after a short time to avoid blocking
			if len(eventsCh) >= 2 {
				return
			}
		}

		if err := scanner.Err(); err != nil {
			errCh <- err
		}
	}()

	// Give stream time to connect
	time.Sleep(100 * time.Millisecond)

	// Trigger an org update to emit an event
	_, _ = do[domain.OrganisationModel](t, "PATCH", fmt.Sprintf("/orgs/%s", uuidToString(org.ID)), map[string]string{"name": "Updated Org"}, token)

	// Wait for events or timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	for {
		select {
		case evt := <-eventsCh:
			evtType, ok := evt["_type"].(string)
			if ok && evtType == "org.org.updated" {
				payload, ok := evt["payload"].(map[string]any)
				if ok && payload != nil {
					return
				}
			}
		case err := <-errCh:
			t.Fatalf("stream error: %v", err)
		case <-ctx.Done():
			t.Fatal("expected to receive org.org.updated event, got nothing")
		}
	}
}

// TestNotification_Stream_ReceivesMultipleEvents verifies multiple events are received
func TestNotification_Stream_ReceivesMultipleEvents(t *testing.T) {
	auth := register(t, randomEmail(), "Test User", "SecurePassword123!")
	token := auth.AccessToken
	org := createOrg(t, token)
	project := createProject(t, uuidToString(org.ID), token, "TEST", "Test Project", "private")

	eventsCh := make(chan map[string]any, 10)
	errCh := make(chan error, 1)
	stopCh := make(chan struct{})

	go func() {
		req, err := http.NewRequest("GET", testServer.URL+"/notifications", nil)
		if err != nil {
			errCh <- err
			return
		}
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			errCh <- err
			return
		}
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		eventType := ""
		for {
			select {
			case <-stopCh:
				return
			default:
				if !scanner.Scan() {
					return
				}

				line := scanner.Text()
				if strings.HasPrefix(line, "event: ") {
					eventType = strings.TrimPrefix(line, "event: ")
				} else if strings.HasPrefix(line, "data: ") {
					dataStr := strings.TrimPrefix(line, "data: ")

					var eventData map[string]any
					if err := json.Unmarshal([]byte(dataStr), &eventData); err != nil {
						continue
					}

					eventData["_type"] = eventType
					select {
					case eventsCh <- eventData:
					case <-stopCh:
						return
					}
				}
			}
		}
	}()

	time.Sleep(100 * time.Millisecond)

	// Trigger multiple events
	_, _ = do[domain.ProjectModel](t, "PATCH", fmt.Sprintf("/projects/%s", uuidToString(project.ID)), map[string]any{
		"name":        "Updated Project",
		"description": "New description",
	}, token)

	_, _ = do[domain.OrganisationModel](t, "PATCH", fmt.Sprintf("/orgs/%s", uuidToString(org.ID)), map[string]string{"name": "Updated Org 2"}, token)

	// Collect events
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	projectUpdates := 0
	orgUpdates := 0

	for {
		select {
		case evt := <-eventsCh:
			evtType, ok := evt["_type"].(string)
			if ok {
				if evtType == "project.project.updated" {
					projectUpdates++
				} else if evtType == "org.org.updated" {
					orgUpdates++
				}
			}

			if projectUpdates > 0 && orgUpdates > 0 {
				close(stopCh)
				return
			}
		case err := <-errCh:
			t.Fatalf("stream error: %v", err)
		case <-ctx.Done():
			if projectUpdates == 0 || orgUpdates == 0 {
				t.Fatalf("expected both project and org events, got %d project updates and %d org updates",
					projectUpdates, orgUpdates)
			}
			return
		}
	}
}

// TestNotification_Stream_EventFormatCorrect verifies event format has type and payload
func TestNotification_Stream_EventFormatCorrect(t *testing.T) {
	auth := register(t, randomEmail(), "Test User", "SecurePassword123!")
	token := auth.AccessToken
	org := createOrg(t, token)

	eventsCh := make(chan map[string]any, 5)
	errCh := make(chan error, 1)

	go func() {
		req, err := http.NewRequest("GET", testServer.URL+"/notifications", nil)
		if err != nil {
			errCh <- err
			return
		}
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			errCh <- err
			return
		}
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		eventType := ""
		for scanner.Scan() {
			line := scanner.Text()

			if strings.HasPrefix(line, "event: ") {
				eventType = strings.TrimPrefix(line, "event: ")
			} else if strings.HasPrefix(line, "data: ") {
				dataStr := strings.TrimPrefix(line, "data: ")

				var eventData map[string]any
				if err := json.Unmarshal([]byte(dataStr), &eventData); err != nil {
					errCh <- fmt.Errorf("failed to unmarshal event: %v", err)
					return
				}

				eventData["_type"] = eventType
				eventsCh <- eventData
			}

			if len(eventsCh) >= 1 {
				return
			}
		}
	}()

	time.Sleep(100 * time.Millisecond)

	_, _ = do[domain.OrganisationModel](t, "PATCH", fmt.Sprintf("/orgs/%s", uuidToString(org.ID)), map[string]string{"name": "Test Org"}, token)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	for {
		select {
		case evt := <-eventsCh:
			// Verify event has type field
			evtType, hasType := evt["_type"].(string)
			if !hasType || evtType == "" {
				t.Fatalf("event missing type field: %v", evt)
			}

			// Verify event has payload field
			payload, hasPayload := evt["payload"]
			if !hasPayload || payload == nil {
				t.Fatalf("event missing payload field: %v", evt)
			}

			// Verify payload is a map
			_, isMap := payload.(map[string]any)
			if !isMap {
				t.Fatalf("payload is not a map: %T", payload)
			}

			return
		case err := <-errCh:
			t.Fatalf("stream error: %v", err)
		case <-ctx.Done():
			t.Fatal("timeout waiting for event")
		}
	}
}
