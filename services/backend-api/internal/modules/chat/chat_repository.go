package chat

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type ChatRepository interface {
	SaveMessage(ctx context.Context, msg *Message) error
	GetSessionHistory(ctx context.Context, sessionID uuid.UUID) ([]*Message, error)
}

type PostgresChatRepository struct {
	db *sql.DB
}

func NewPostgresChatRepository(db *sql.DB) *PostgresChatRepository {
	return &PostgresChatRepository{db: db}
}

func (r *PostgresChatRepository) SaveMessage(ctx context.Context, msg *Message) error {
	query := `
		INSERT INTO chat_messages (id, session_id, role, content, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query, msg.ID, msg.SessionID, msg.Role, msg.Content, msg.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to persist message to postgres: %w", err)
	}
	return nil
}


func (r *PostgresChatRepository) GetSessionHistory(ctx context.Context, sessionID uuid.UUID) ([]*Message, error) {
	query := `
		SELECT id, session_id, role, content, created_at 
		FROM chat_messages 
		WHERE session_id = $1 
		ORDER BY created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to query session log: %w", err)
	}
	defer rows.Close()

	var history []*Message

	for rows.Next() {
		var msg Message
		err := rows.Scan(&msg.ID, &msg.SessionID, &msg.Role, &msg.Content, &msg.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan database row: %w", err)
		}
		history = append(history, &msg)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row data stream corruption detected: %w", err)
	}

	return history, nil
}