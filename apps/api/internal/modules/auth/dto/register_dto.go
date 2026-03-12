package dto

// RegisterRequest represents register input.
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=128"`
}

// RegisterResponse represents register output.
type RegisterResponse struct {
	CredentialID      string `json:"credential_id"`
	Email             string `json:"email"`
	EmailVerified     bool   `json:"email_verified"`
	VerificationToken string `json:"verification_token,omitempty"`
}
