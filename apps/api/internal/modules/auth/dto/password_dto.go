package dto

// ForgotPasswordRequest starts password reset process.
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordRequest applies password reset by token.
type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required,min=20"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=128"`
}

// ChangePasswordRequest changes password with current password validation.
type ChangePasswordRequest struct {
	CredentialID string `json:"credential_id" validate:"required,uuid4"`
	OldPassword  string `json:"old_password" validate:"required,min=8,max=128"`
	NewPassword  string `json:"new_password" validate:"required,min=8,max=128"`
}
