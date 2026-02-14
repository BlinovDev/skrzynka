package http

import (
	"encoding/json"
	"net/http"

	"anon-skrzynka/app/model"
	"anon-skrzynka/app/storage"
)

// OnMessageCreated is called after a message is successfully created (e.g. to broadcast via WebSocket).
type OnMessageCreated func(*model.Message)

// CreateMessageRequest is the JSON body for creating a message (sender_id, recipient_id, payload required).
type CreateMessageRequest struct {
	SenderID    string `json:"sender_id"`
	RecipientID string `json:"recipient_id"`
	Payload     string `json:"payload"`
}

// Handlers holds dependencies for HTTP handlers.
type Handlers struct {
	Repo             storage.Repository
	OnMessageCreated OnMessageCreated
}

// NewHandlers returns Handlers with the given repository.
func NewHandlers(repo storage.Repository) *Handlers {
	return &Handlers{Repo: repo}
}

// CreateMessage handles POST /messages: creates a message and returns it.
func (h *Handlers) CreateMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req CreateMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if req.SenderID == "" || req.RecipientID == "" || req.Payload == "" {
		http.Error(w, "sender_id, recipient_id and payload are required", http.StatusBadRequest)
		return
	}
	msg := &model.Message{
		SenderID:    req.SenderID,
		RecipientID: req.RecipientID,
		Payload:     req.Payload,
	}
	if err := h.Repo.Create(msg); err != nil {
		http.Error(w, "create failed", http.StatusInternalServerError)
		return
	}
	if h.OnMessageCreated != nil {
		h.OnMessageCreated(msg)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(msg)
}

// GetDialog handles GET /messages?sender_id=&recipient_id= and returns all messages between the two.
func (h *Handlers) GetDialog(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	senderID := r.URL.Query().Get("sender_id")
	recipientID := r.URL.Query().Get("recipient_id")
	if senderID == "" || recipientID == "" {
		http.Error(w, "sender_id and recipient_id are required", http.StatusBadRequest)
		return
	}
	msgs, err := h.Repo.GetDialog(senderID, recipientID)
	if err != nil {
		http.Error(w, "get dialog failed", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(msgs)
}
