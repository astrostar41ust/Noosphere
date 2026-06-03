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
	aiClient AIClient
}

func NewDefaultChatService(repo ChatRepository, aiClient AIClient) *DefaultChatService {
	return &DefaultChatService{repo: repo, aiClient: aiClient}
}

func (s *DefaultChatService) SendMessage(ctx context.Context, sessionID uuid.UUID, role string, content string) (*Message, error) {
	if content == "" {
		return nil, fmt.Errorf("message content cannot be completely empty")
	}

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


	history, err := s.repo.GetSessionHistory(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to compile historical frame context for model: %w", err)
	}

	var plainHistory []Message
	for _, m := range history {
		plainHistory = append(plainHistory, *m)
	}

	aiText, err := s.aiClient.GenerateResponse(ctx, plainHistory)
	if err != nil {
		return nil, fmt.Errorf("ai engine generation cluster pipeline failure: %w", err)
	}

	assistantMsg := &Message{
		ID:        uuid.New(),
		SessionID: sessionID,
		Role:      "assistant",
		Content:   aiText,
		CreatedAt: time.Now(),
	}
	if err := s.repo.SaveMessage(ctx, assistantMsg); err != nil {
		return nil, fmt.Errorf("failed to securely store generated assistant message: %w", err)
	}

	return assistantMsg, nil
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