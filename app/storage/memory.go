package storage

import (
	"sync"
	"time"

	"anon-skrzynka/app/model"

	"github.com/google/uuid"
)

// Memory is an in-memory Repository implementation.
type Memory struct {
	mu   sync.RWMutex
	msgs []*model.Message
}

// NewMemory returns a new in-memory repository.
func NewMemory() *Memory {
	return &Memory{msgs: make([]*model.Message, 0)}
}

// Create stores a new message and sets ID and CreatedAt.
func (m *Memory) Create(msg *model.Message) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if msg.ID == "" {
		msg.ID = uuid.New().String()
	}
	if msg.CreatedAt.IsZero() {
		msg.CreatedAt = time.Now().UTC()
	}
	m.msgs = append(m.msgs, msg)
	return nil
}

// GetDialog returns all messages between the two participant IDs in either direction.
func (m *Memory) GetDialog(senderID, recipientID string) ([]*model.Message, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*model.Message, 0)
	for _, msg := range m.msgs {
		if (msg.SenderID == senderID && msg.RecipientID == recipientID) ||
			(msg.SenderID == recipientID && msg.RecipientID == senderID) {
			out = append(out, msg)
		}
	}
	return out, nil
}
