package storage

import "anon-skrzynka/app/model"

// Repository persists and retrieves messages. Implementations must be safe for concurrent use.
type Repository interface {
	// Create stores a new message and returns it with ID and CreatedAt set.
	Create(msg *model.Message) error
	// GetDialog returns all messages between the two given participant UUIDs (both directions).
	GetDialog(senderID, recipientID string) ([]*model.Message, error)
}
