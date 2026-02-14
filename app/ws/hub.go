package ws

import (
	"encoding/json"
	"sync"

	"anon-skrzynka/app/model"
)

// Hub manages WebSocket clients and broadcasts new messages to relevant participants.
type Hub struct {
	mu      sync.RWMutex
	clients map[string]map[*Client]struct{}
}

// NewHub returns a new WebSocket hub.
func NewHub() *Hub {
	return &Hub{clients: make(map[string]map[*Client]struct{})}
}

// Register adds a client for the given user ID so it can receive messages where that user is sender or recipient.
func (h *Hub) Register(userID string, c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.clients[userID] == nil {
		h.clients[userID] = make(map[*Client]struct{})
	}
	h.clients[userID][c] = struct{}{}
}

// Unregister removes a client from the hub for the given user ID.
func (h *Hub) Unregister(userID string, c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if m := h.clients[userID]; m != nil {
		delete(m, c)
		if len(m) == 0 {
			delete(h.clients, userID)
		}
	}
}

// BroadcastMessage sends the message as JSON to all clients registered for its sender_id or recipient_id.
func (h *Hub) BroadcastMessage(msg *model.Message) {
	body, err := json.Marshal(msg)
	if err != nil {
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, id := range []string{msg.SenderID, msg.RecipientID} {
		for c := range h.clients[id] {
			select {
			case c.Send <- body:
			default:
			}
		}
	}
}
