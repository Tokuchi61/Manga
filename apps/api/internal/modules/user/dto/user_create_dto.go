package dto

// CreateUserRequest creates a user account linked to a verified auth credential.
type CreateUserRequest struct {
	CredentialID string `json:"credential_id" validate:"required,uuid4"`
	Username     string `json:"username" validate:"required,min=3,max=32"`
	DisplayName  string `json:"display_name,omitempty" validate:"omitempty,max=64"`
	Bio          string `json:"bio,omitempty" validate:"omitempty,max=280"`
}

// CreateUserResponse returns canonical user account identity.
type CreateUserResponse struct {
	UserID       string `json:"user_id"`
	CredentialID string `json:"credential_id"`
	Username     string `json:"username"`
	AccountState string `json:"account_state"`
}
