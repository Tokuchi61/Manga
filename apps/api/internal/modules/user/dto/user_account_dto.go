package dto

// UpdateAccountStateRequest changes account lifecycle status.
type UpdateAccountStateRequest struct {
	UserID       string `json:"-" validate:"required,uuid4"`
	ActorUserID  string `json:"-" validate:"required,uuid4"`
	ActorScope   string `json:"-" validate:"required,oneof=self admin"`
	AccountState string `json:"account_state" validate:"required,oneof=active deactivated banned"`
	Reason       string `json:"reason,omitempty" validate:"omitempty,max=120"`
}

// AccountStateResponse returns current account status.
type AccountStateResponse struct {
	Status       string `json:"status"`
	AccountState string `json:"account_state"`
}
