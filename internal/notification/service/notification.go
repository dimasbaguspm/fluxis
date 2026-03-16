package service

import (
	"context"
	"sync"

	"github.com/dimasbaguspm/fluxis/pkg/pubsub"
)

// Stream subscribes to all domain event channels and forwards events to the send function.
// It launches goroutines to listen on each channel and funnels events through a single
// goroutine to ensure send is never called concurrently.
func (s *Service) Stream(ctx context.Context, send func(pubsub.Event) error) error {
	eventCh := make(chan pubsub.Event, 10)
	closeCh := make(chan struct{})
	var wg sync.WaitGroup

	channels := []pubsub.EventType{
		pubsub.Ticket,
		pubsub.Sprint,
		pubsub.Board,
		pubsub.Org,
		pubsub.Project,
		pubsub.User,
	}

	for _, eventType := range channels {
		wg.Add(1)
		go func(et pubsub.EventType) {
			defer wg.Done()
			forwardFn := func(ctx context.Context, e pubsub.Event) error {
				select {
				case eventCh <- e:
				case <-ctx.Done():
					return ctx.Err()
				case <-closeCh:
					return context.Canceled
				}
				return nil
			}
			s.Bus.Subscribe(ctx, pubsub.Channel(et), forwardFn)
		}(eventType)
	}

	for {
		select {
		case e := <-eventCh:
			if err := send(e); err != nil {
				close(closeCh)
				wg.Wait()
				return err
			}
		case <-ctx.Done():
			close(closeCh)
			wg.Wait()
			return ctx.Err()
		}
	}
}
