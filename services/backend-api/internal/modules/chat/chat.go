package chat

import "time"

type Message struct {
	ID        string    `json:"id"`
	SessionID string    `json:"session_id"`
	Role      string    `json:"role"` // "user" or "assistant"
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type Repository interface {
	SaveMessage(msg *Message) error
	GetHistory(sessionID string) ([]Message, error)
}

type Service interface {
	SendMessage(sessionID string, userPrompt string) (*Message, error)
	GetChatSessionHistory(sessionID string) ([]Message, error)
}