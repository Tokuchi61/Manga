package dto

// SendVerificationRequest triggers verification mail for credential.
type SendVerificationRequest struct {
	CredentialID string `json:"-" validate:"required,uuid4"`
}

// ConfirmVerificationRequest confirms email verification token.
type ConfirmVerificationRequest struct {
	Token string `json:"token" validate:"required,min=32,max=512"`
}
