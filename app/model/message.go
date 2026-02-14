package model

import "time"

// Message is the canonical message entity. Payload is opaque to the server.
type Message struct {
	ID          string    `json:"id"`
	SenderID    string    `json:"sender_id"`
	RecipientID string    `json:"recipient_id"`
	Payload     string    `json:"payload"`
	CreatedAt   time.Time `json:"created_at"`
}
