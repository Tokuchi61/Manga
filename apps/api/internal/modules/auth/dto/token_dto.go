package dto

import "time"

// RefreshTokenRequest represents refresh token rotation input.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required,min=20"`
}

// RefreshTokenResponse represents refresh token rotation output.
type RefreshTokenResponse struct {
	SessionID            string    `json:"session_id"`
	AccessToken          string    `json:"access_token"`
	RefreshToken         string    `json:"refresh_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}
