package auth

import (
	"time"

	"noosphere/backend-api/internal/modules/user"
)

type RegisterRequest struct {
	Username        string `json:"username"        validate:"required,min=3,max=20"`
	Email           string `json:"email"           validate:"required,email"`
	Password        string `json:"password"        validate:"required,min=6"`
	ConfirmPassword string `json:"confirmPassword" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

type TokenResponse struct {
	User        UserResponse `json:"user"`
	AccessToken string       `json:"access_token"`
}

func MapUserToResponse(u *user.User) UserResponse {
	return UserResponse{
		ID:        u.ID.String(),
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.UTC().Format(time.RFC3339),
	}
}
