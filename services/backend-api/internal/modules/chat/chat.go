package chat

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID `json:"id"`
	SessionID uuid.UUID `json:"session_id"`
	Role      string    `json:"role"` // "user" or "assistant"
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
