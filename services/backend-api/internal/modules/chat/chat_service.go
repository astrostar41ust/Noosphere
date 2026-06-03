package chat

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)


type ChatService interface {
	SendMessage(ctx context.Context, sessionID uuid.UUID, role string, content string) (*Message, error)
	GetChatHistory(ctx context.Context, sessionID uuid.UUID) ([]*Message, error)
}

type DefaultChatService struct {
	repo ChatRepository
}

func NewDefaultChatService(repo ChatRepository) *DefaultChatService {
	return &DefaultChatService{repo: repo}
}

func (s *DefaultChatService) SendMessage(ctx context.Context, sessionID uuid.UUID, role string, content string) (*Message, error) {

	msg := &Message{
		ID:        uuid.New(),
		SessionID: sessionID,
		Role:      role,
		Content:   content,
		CreatedAt: time.Now(),
	}

	err := s.repo.SaveMessage(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("service failed to save message to DB: %w", err)
	}

	return msg, nil
}

func (s *DefaultChatService) GetChatHistory(ctx context.Context, sessionID uuid.UUID) ([]*Message, error) {
	if sessionID == uuid.Nil {
		return nil, fmt.Errorf("invalid or empty session identifier provided")
	}

	history, err := s.repo.GetSessionHistory(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("service failed to gather history logs: %w", err)
	}

	return history, nil
}