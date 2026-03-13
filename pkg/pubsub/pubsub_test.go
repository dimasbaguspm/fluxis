package pubsub_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/dimasbaguspm/fluxis/pkg/pubsub"
)

func TestChannel_DerivesCorrectChannel(t *testing.T) {
	tests := []struct {
		eventType pubsub.EventType
		expected  string
	}{
		// Org events
		{pubsub.OrgCreated, "events:org"},
		{pubsub.OrgMemberAdded, "events:org"},

		// Project events
		{pubsub.ProjectCreated, "events:project"},
		{pubsub.ProjectVisibilityUpdated, "events:project"},

		// Sprint events
		{pubsub.SprintCreated, "events:sprint"},
		{pubsub.SprintCompleted, "events:sprint"},

		// Board events
		{pubsub.BoardCreated, "events:board"},
		{pubsub.BoardColumnCreated, "events:board"},

		// Ticket events
		{pubsub.TicketCreated, "events:ticket"},
		{pubsub.TicketMovedToBoard, "events:ticket"},
	}

	for _, tt := range tests {
		t.Run(string(tt.eventType), func(t *testing.T) {
			got := pubsub.Channel(tt.eventType)
			if got != tt.expected {
				t.Errorf("Channel(%s) = %s, want %s", tt.eventType, got, tt.expected)
			}
		})
	}
}

func TestMemory_PublishSubscribe_ReceivesEvent(t *testing.T) {
	ctx := context.Background()
	bus := pubsub.New()
	defer bus.Close()

	received := &sync.WaitGroup{}
	received.Add(1)

	var gotEvent pubsub.Event
	handler := func(_ context.Context, e pubsub.Event) error {
		gotEvent = e
		received.Done()
		return nil
	}

	channel := "events:ticket"
	go bus.Subscribe(ctx, channel, handler)

	// Small delay to allow goroutine to register.
	time.Sleep(10 * time.Millisecond)

	payload := map[string]string{"id": "123", "title": "test"}
	bus.Publish(ctx, pubsub.TicketCreated, payload)

	done := make(chan struct{})
	go func() {
		received.Wait()
		close(done)
	}()

	select {
	case <-done:
		if gotEvent.Type != pubsub.TicketCreated {
			t.Errorf("Event type = %s, want %s", gotEvent.Type, pubsub.TicketCreated)
		}
		if gotEvent.Payload["id"] != "123" {
			t.Errorf("Event payload[id] = %s, want 123", gotEvent.Payload["id"])
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("handler not called within 500ms")
	}
}

func TestMemory_Publish_BestEffort(t *testing.T) {
	ctx := context.Background()
	bus := pubsub.New()
	defer bus.Close()

	// Publish with no subscribers should return nil error.
	err := bus.Publish(ctx, pubsub.TicketCreated, map[string]string{"id": "1"})
	if err != nil {
		t.Errorf("Publish with no subscribers returned error: %v", err)
	}
}

func TestMemory_Subscribe_ExitsOnContextCancel(t *testing.T) {
	bus := pubsub.New()
	defer bus.Close()

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})

	go func() {
		bus.Subscribe(ctx, "events:ticket", func(_ context.Context, _ pubsub.Event) error {
			return nil
		})
		close(done)
	}()

	// Small delay to allow goroutine to register.
	time.Sleep(10 * time.Millisecond)

	cancel()

	select {
	case <-done:
		// Expected: Subscribe exited after context cancellation.
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Subscribe did not exit within 500ms after context cancel")
	}
}

func TestMemory_Publish_FanOut(t *testing.T) {
	ctx := context.Background()
	bus := pubsub.New()
	defer bus.Close()

	wg := &sync.WaitGroup{}
	wg.Add(2)

	channel := "events:ticket"

	// First subscriber.
	go bus.Subscribe(ctx, channel, func(_ context.Context, e pubsub.Event) error {
		if e.Type != pubsub.TicketCreated {
			t.Errorf("subscriber1: Event type = %s, want %s", e.Type, pubsub.TicketCreated)
		}
		wg.Done()
		// Keep consuming to allow second subscriber to also receive.
		<-make(chan pubsub.Event)
		return nil
	})

	// Second subscriber.
	go bus.Subscribe(ctx, channel, func(_ context.Context, e pubsub.Event) error {
		if e.Type != pubsub.TicketCreated {
			t.Errorf("subscriber2: Event type = %s, want %s", e.Type, pubsub.TicketCreated)
		}
		wg.Done()
		// Keep consuming to allow proper cleanup.
		<-make(chan pubsub.Event)
		return nil
	})

	// Small delay to allow goroutines to register.
	time.Sleep(10 * time.Millisecond)

	bus.Publish(ctx, pubsub.TicketCreated, map[string]string{"id": "1"})

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Both subscribers received the event.
	case <-time.After(500 * time.Millisecond):
		t.Fatal("not all subscribers received event within 500ms")
	}
}

func TestMemory_Subscribe_IsolatedByChannel(t *testing.T) {
	ctx := context.Background()
	bus := pubsub.New()
	defer bus.Close()

	ticketReceived := false
	boardReceived := false

	ticketDone := make(chan struct{})
	boardDone := make(chan struct{})

	// Subscribe to ticket channel.
	go func() {
		count := 0
		bus.Subscribe(ctx, "events:ticket", func(_ context.Context, _ pubsub.Event) error {
			ticketReceived = true
			count++
			if count == 1 {
				close(ticketDone)
			}
			// Keep consuming.
			<-make(chan pubsub.Event)
			return nil
		})
	}()

	// Subscribe to board channel.
	go func() {
		count := 0
		bus.Subscribe(ctx, "events:board", func(_ context.Context, _ pubsub.Event) error {
			boardReceived = true
			count++
			if count == 1 {
				close(boardDone)
			}
			// Keep consuming.
			<-make(chan pubsub.Event)
			return nil
		})
	}()

	// Small delay to allow goroutines to register.
	time.Sleep(10 * time.Millisecond)

	// Publish to ticket channel only.
	bus.Publish(ctx, pubsub.TicketCreated, map[string]string{"id": "1"})

	select {
	case <-ticketDone:
		// Ticket subscriber received the event.
	case <-time.After(500 * time.Millisecond):
		t.Fatal("ticket subscriber did not receive event within 500ms")
	}

	// Verify board subscriber did not receive the event.
	select {
	case <-boardDone:
		t.Fatal("board subscriber should not have received ticket event")
	case <-time.After(100 * time.Millisecond):
		// Expected: board subscriber did not receive ticket event.
	}

	if !ticketReceived {
		t.Error("ticket subscriber did not receive event")
	}
	if boardReceived {
		t.Error("board subscriber should not have received event")
	}
}
