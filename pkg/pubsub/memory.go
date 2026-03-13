package pubsub

import (
	"context"
	"log/slog"
	"sync"
)

const subscriberBufSize = 64

type memoryBus struct {
	mu   sync.Mutex
	subs map[string][]chan Event
}

func New() Bus {
	slog.Info("[PubSub]: Initializing in-memory pub/sub bus")
	return &memoryBus{subs: make(map[string][]chan Event)}
}

func (b *memoryBus) Publish(_ context.Context, et EventType, payload map[string]string) error {
	ch := Channel(et)
	e := Event{Type: et, Payload: payload}

	b.mu.Lock()

	subscribers := make([]chan Event, len(b.subs[ch]))
	copy(subscribers, b.subs[ch])
	b.mu.Unlock()

	for _, sub := range subscribers {
		select {
		case sub <- e:
		default:
			slog.Warn("[PubSub]: subscriber channel full, dropping event",
				"channel", ch, "type", string(et))
		}
	}
	return nil
}

func (b *memoryBus) Subscribe(ctx context.Context, channel string, handler func(context.Context, Event) error) {
	ch := make(chan Event, subscriberBufSize)

	b.mu.Lock()
	b.subs[channel] = append(b.subs[channel], ch)
	b.mu.Unlock()

	defer func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		subs := b.subs[channel]
		for i, s := range subs {
			if s == ch {
				b.subs[channel] = append(subs[:i], subs[i+1:]...)
				break
			}
		}
	}()

	for {
		select {
		case e, ok := <-ch:
			if !ok {
				return
			}
			if err := handler(ctx, e); err != nil {
				slog.Error("[PubSub]: subscriber handler error",
					"channel", channel, "type", string(e.Type), "error", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (b *memoryBus) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, subs := range b.subs {
		for _, ch := range subs {
			close(ch)
		}
	}
	b.subs = make(map[string][]chan Event)
	return nil
}
