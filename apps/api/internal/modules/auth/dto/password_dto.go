package dto

// ForgotPasswordRequest requests reset token by email.
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email,max=320"`
}

// ResetPasswordRequest applies password reset using token.
type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required,min=32,max=512"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=128"`
}

// ChangePasswordRequest changes password for authenticated credential.
type ChangePasswordRequest struct {
	CredentialID string `json:"-" validate:"required,uuid4"`
	OldPassword  string `json:"old_password" validate:"required,min=8,max=128"`
	NewPassword  string `json:"new_password" validate:"required,min=8,max=128"`
}
