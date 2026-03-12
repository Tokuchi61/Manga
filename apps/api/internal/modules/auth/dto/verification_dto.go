package dto

// SendVerificationRequest starts or resends verification flow.
type SendVerificationRequest struct {
	CredentialID string `json:"credential_id" validate:"required,uuid4"`
}

// ConfirmVerificationRequest verifies account email.
type ConfirmVerificationRequest struct {
	Token string `json:"token" validate:"required,min=20"`
}
