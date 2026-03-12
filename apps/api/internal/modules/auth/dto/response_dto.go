package dto

// OperationResponse is a generic success payload.
type OperationResponse struct {
	Status string `json:"status"`
}

// TokenDispatchResponse is returned for email verification and reset requests.
type TokenDispatchResponse struct {
	Status string `json:"status"`
	Token  string `json:"token,omitempty"`
}
