package helper

import (
	"net/http"
)

// SSEHandler streams events to the client using Server-Sent Events (SSE).
// The eventWriter function should write events using the provided http.ResponseWriter and return when done.
func SSEHandler(eventWriter func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
			return
		}
		// Call the event writer
		eventWriter(w, r)
		// Ensure everything is sent
		flusher.Flush()
	}
}
