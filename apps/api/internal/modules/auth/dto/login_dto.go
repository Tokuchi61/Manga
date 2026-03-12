package dto

import "time"

// LoginRequest represents login input.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=128"`
	Device   string `json:"device" validate:"max=128"`
	IP       string `json:"ip" validate:"max=64"`
}

// LoginResponse represents login output.
type LoginResponse struct {
	CredentialID         string    `json:"credential_id"`
	SessionID            string    `json:"session_id"`
	AccessToken          string    `json:"access_token"`
	RefreshToken         string    `json:"refresh_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}
