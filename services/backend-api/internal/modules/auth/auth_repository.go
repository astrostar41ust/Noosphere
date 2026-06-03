package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type AuthRepository interface {
	SaveRefreshToken(ctx context.Context, token *RefreshToken) error
	GetRefreshToken(ctx context.Context, tokenStr string) (*RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, tokenStr string) error
	DeleteUserRefreshTokens(ctx context.Context, userID uuid.UUID) error
}

type PostgresAuthRepository struct {
	db *sql.DB
}

func NewPostgresAuthRepository(db *sql.DB) *PostgresAuthRepository {
	return &PostgresAuthRepository{db: db}
}

func (r *PostgresAuthRepository) SaveRefreshToken(ctx context.Context, token *RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, token, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query, token.ID, token.UserID, token.Token, token.ExpiresAt, token.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to save refresh token to postgres: %w", err)
	}
	return nil
}

func (r *PostgresAuthRepository) GetRefreshToken(ctx context.Context, tokenStr string) (*RefreshToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at
		FROM refresh_tokens
		WHERE token = $1
	`
	var token RefreshToken
	err := r.db.QueryRowContext(ctx, query, tokenStr).Scan(&token.ID, &token.UserID, &token.Token, &token.ExpiresAt, &token.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to query refresh token: %w", err)
	}
	return &token, nil
}

func (r *PostgresAuthRepository) DeleteRefreshToken(ctx context.Context, tokenStr string) error {
	query := `
		DELETE FROM refresh_tokens
		WHERE token = $1
	`
	_, err := r.db.ExecContext(ctx, query, tokenStr)
	if err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}
	return nil
}

func (r *PostgresAuthRepository) DeleteUserRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	query := `
		DELETE FROM refresh_tokens
		WHERE user_id = $1
	`
	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user refresh tokens: %w", err)
	}
	return nil
}
