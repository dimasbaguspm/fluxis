package handler

import (
	"encoding/json"
	"net/http"

	"github.com/dimasbaguspm/fluxis/pkg/pubsub"
)

// Stream godoc
//
//	@Summary		Stream domain events
//	@Description	Opens a Server-Sent Events stream that broadcasts all domain events (ticket, sprint, board, org, project, user) in real-time
//	@Tags			notification
//	@Produce		text/event-stream
//	@Success		200	{object}	string	"Event stream with events in format: event: <type>\\ndata: <JSON>\\n\\n"
//	@Failure		401	{object}	httpx.ErrBlock
//	@Failure		500	{object}	httpx.ErrBlock
//	@Security		BearerAuth
//	@Router			/notifications [get]
func (h *Handler) Stream(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	sendFn := func(e pubsub.Event) error {
		eventData := map[string]any{
			"type":    string(e.Type),
			"payload": e.Payload,
		}
		data, err := json.Marshal(eventData)
		if err != nil {
			return err
		}

		_, err = w.Write([]byte("event: " + string(e.Type) + "\n"))
		if err != nil {
			return err
		}

		_, err = w.Write([]byte("data: " + string(data) + "\n\n"))
		if err != nil {
			return err
		}

		flusher.Flush()
		return nil
	}

	_ = h.svc.Stream(r.Context(), sendFn)
}
