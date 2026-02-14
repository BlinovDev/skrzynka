package http

import "net/http"

// Router returns the application HTTP router with all routes registered.
func Router(h *Handlers) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /messages", h.CreateMessage)
	mux.HandleFunc("GET /messages", h.GetDialog)
	return mux
}
