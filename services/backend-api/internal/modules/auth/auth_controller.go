package auth

import (
	"net/http"
	"time"

	"noosphere/backend-api/internal/pkg/httpio"

	"github.com/go-playground/validator/v10"
)

type AuthController struct {
	service  AuthService
	validate *validator.Validate
}

func NewAuthController(service AuthService) *AuthController {
	return &AuthController{
		service:  service,
		validate: validator.New(),
	}
}

func (c *AuthController) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := httpio.DecodeAndValidate(w, r, &req, c.validate); err != nil {
		return
	}

	if req.Password != req.ConfirmPassword {
		httpio.RespondWithError(w, http.StatusBadRequest, "Passphrase verification nodes do not match")
		return
	}

	user, err := c.service.Register(r.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		httpio.RespondWithError(w, http.StatusConflict, err.Error())
		return
	}

	accessToken, refreshToken, err := c.service.GenerateTokens(r.Context(), user)
	if err != nil {
		httpio.RespondWithError(w, http.StatusInternalServerError, "Token generation failed")
		return
	}

	c.setRefreshTokenCookie(w, refreshToken, 168*time.Hour)

	response := TokenResponse{
		User:        MapUserToResponse(user),
		AccessToken: accessToken,
	}

	httpio.RespondWithJSON(w, http.StatusCreated, response)
}

func (c *AuthController) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := httpio.DecodeAndValidate(w, r, &req, c.validate); err != nil {
		return
	}

	user, err := c.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		httpio.RespondWithError(w, http.StatusUnauthorized, "Authentication credentials validation rejected")
		return
	}

	accessToken, refreshToken, err := c.service.GenerateTokens(r.Context(), user)
	if err != nil {
		httpio.RespondWithError(w, http.StatusInternalServerError, "Token generation failed")
		return
	}

	c.setRefreshTokenCookie(w, refreshToken, 168*time.Hour)

	response := TokenResponse{
		User:        MapUserToResponse(user),
		AccessToken: accessToken,
	}

	httpio.RespondWithJSON(w, http.StatusOK, response)
}

func (c *AuthController) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("__Secure-refresh_token")
	if err != nil {
		httpio.RespondWithError(w, http.StatusUnauthorized, "Missing session verification refresh token")
		return
	}

	accessToken, newRefreshToken, user, err := c.service.RefreshTokens(r.Context(), cookie.Value)
	if err != nil {
		httpio.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	c.setRefreshTokenCookie(w, newRefreshToken, 168*time.Hour)

	response := TokenResponse{
		User:        MapUserToResponse(user),
		AccessToken: accessToken,
	}

	httpio.RespondWithJSON(w, http.StatusOK, response)
}

func (c *AuthController) HandleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("__Secure-refresh_token")
	if err == nil {
		_ = c.service.RevokeToken(r.Context(), cookie.Value)
	}

	c.setRefreshTokenCookie(w, "", 0)
	w.WriteHeader(http.StatusNoContent)
}

func (c *AuthController) setRefreshTokenCookie(w http.ResponseWriter, token string, expiry time.Duration) {
	expires := time.Now().Add(expiry)
	if token == "" {
		expires = time.Unix(0, 0)
	}

	cookie := &http.Cookie{
		Name:     "__Secure-refresh_token",
		Value:    token,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
}
