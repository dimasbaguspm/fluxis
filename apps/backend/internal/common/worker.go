package common

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type Trigger struct {
	Resource string
	ID       string
	Action   string
	Meta     map[string]interface{}
}

type Handler func(t Trigger)

type Worker struct {
	ch       chan Trigger
	stop     chan struct{}
	wg       sync.WaitGroup
	itv      time.Duration
	handler  Handler
	stopping int32
	ctx      context.Context
}

func NewWorker(ctx context.Context, handler Handler) *Worker {
	if handler == nil {
		panic("handler cannot be nil")
	}

	w := &Worker{
		ch:      make(chan Trigger, 1024),
		stop:    make(chan struct{}),
		itv:     10 * time.Second,
		handler: handler,
	}
	w.wg.Add(1)
	go w.run()
	return w
}

// Enqueue adds a trigger to the worker queue.
// Returns immediately; trigger may be dropped if worker is stopping.
func (w *Worker) Enqueue(t Trigger) {
	// worker is shutting down; drop trigger
	if atomic.LoadInt32(&w.stopping) == 1 {
		return
	}
	select {
	case w.ch <- t:
	default:
		// drop trigger if queue full
	}
}

// Stop gracefully shuts down the worker, draining remaining triggers.
func (w *Worker) Stop() {
	if !atomic.CompareAndSwapInt32(&w.stopping, 0, 1) {
		return
	}
	close(w.stop)
	w.wg.Wait()
}

func (w *Worker) run() {
	defer w.wg.Done()

	ticker := time.NewTicker(w.itv)
	defer ticker.Stop()

	pending := make(map[string]Trigger)

	drain := func() {
		if len(pending) == 0 {
			return
		}

		for _, t := range pending {
			w.handler(t)
		}
		pending = make(map[string]Trigger)
	}

	for {
		select {
		case <-w.stop:
			// received the stop request and drain remaining pending queue
			for {
				select {
				case t := <-w.ch:
					key := t.Resource + ":" + t.ID
					pending[key] = t
				default:
					drain()
					return
				}
			}
		case t := <-w.ch:
			// de-duplicate by resource+id
			key := t.Resource + ":" + t.ID
			pending[key] = t
		case <-ticker.C:
			drain()
		}
	}
}
