package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"noosphere/backend-api/internal/modules/user"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type JWTClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

type AuthService interface {
	Register(ctx context.Context, username, email, password string) (*user.User, error)
	Login(ctx context.Context, email, password string) (*user.User, error)
	GenerateTokens(ctx context.Context, u *user.User) (string, string, error)
	RefreshTokens(ctx context.Context, oldTokenStr string) (string, string, *user.User, error)
	RevokeToken(ctx context.Context, tokenStr string) error
}

type DefaultAuthService struct {
	authRepo           AuthRepository
	userRepo           user.UserRepository
	jwtSecret          []byte
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

func NewDefaultAuthService(authRepo AuthRepository, userRepo user.UserRepository, jwtSecret []byte, accessExpiry, refreshExpiry time.Duration) *DefaultAuthService {
	return &DefaultAuthService{
		authRepo:           authRepo,
		userRepo:           userRepo,
		jwtSecret:          jwtSecret,
		accessTokenExpiry:  accessExpiry,
		refreshTokenExpiry: refreshExpiry,
	}
}

func (s *DefaultAuthService) Register(ctx context.Context, username, email, password string) (*user.User, error) {
	existingEmail, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to complete registration node check: %w", err)
	}
	if existingEmail != nil {
		return nil, fmt.Errorf("a user profile is already registered with this email address")
	}


	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to validate password credentials: %w", err)
	}

	newUser := &user.User{
		ID:           uuid.New(),
		Username:     username,
		Email:        email,
		PasswordHash: string(passwordHash),
		CreatedAt:    time.Now(),
	}

	err = s.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("failed to compile client user node registration")
	}

	return newUser, nil
}

func (s *DefaultAuthService) Login(ctx context.Context, email, password string) (*user.User, error) {
	u, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("invalid email or password credentials")
		}
		return nil, fmt.Errorf("failed to lookup user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid email or password credentials")
	}

	return u, nil
}

func (s *DefaultAuthService) GenerateTokens(ctx context.Context, u *user.User) (string, string, error) {
	now := time.Now()
	accessClaims := &JWTClaims{
		UserID:   u.ID.String(),
		Username: u.Username,
		Email:    u.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   u.ID.String(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenStr, err := accessToken.SignedString(s.jwtSecret)
	if err != nil {
		return "", "", fmt.Errorf("failed to sign access token: %w", err)
	}

	refreshBytes := make([]byte, 32)
	if _, err := rand.Read(refreshBytes); err != nil {
		return "", "", fmt.Errorf("failed to generate secure refresh token entropy: %w", err)
	}
	refreshTokenStr := hex.EncodeToString(refreshBytes)

	refreshTokenEntity := &RefreshToken{
		ID:        uuid.New(),
		UserID:    u.ID,
		Token:     refreshTokenStr,
		ExpiresAt: now.Add(s.refreshTokenExpiry),
		CreatedAt: now,
	}

	err = s.authRepo.SaveRefreshToken(ctx, refreshTokenEntity)
	if err != nil {
		return "", "", fmt.Errorf("failed to persist refresh token to database: %w", err)
	}

	return accessTokenStr, refreshTokenStr, nil
}

func (s *DefaultAuthService) RefreshTokens(ctx context.Context, oldTokenStr string) (string, string, *user.User, error) {
	storedToken, err := s.authRepo.GetRefreshToken(ctx, oldTokenStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", nil, fmt.Errorf("refresh token not found or already revoked")
		}
		return "", "", nil, fmt.Errorf("failed to query refresh token: %w", err)
	}

	if time.Now().After(storedToken.ExpiresAt) {
		_ = s.authRepo.DeleteRefreshToken(ctx, oldTokenStr)
		return "", "", nil, fmt.Errorf("refresh token has expired")
	}

	u, err := s.userRepo.GetUserByID(ctx, storedToken.UserID)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to retrieve user associated with token: %w", err)
	}

	err = s.authRepo.DeleteRefreshToken(ctx, oldTokenStr)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to delete rotated token: %w", err)
	}

	accessTokenStr, newRefreshTokenStr, err := s.GenerateTokens(ctx, u)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to generate new token pair: %w", err)
	}

	return accessTokenStr, newRefreshTokenStr, u, nil
}

func (s *DefaultAuthService) RevokeToken(ctx context.Context, tokenStr string) error {
	return s.authRepo.DeleteRefreshToken(ctx, tokenStr)
}
