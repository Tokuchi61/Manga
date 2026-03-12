package dto

// CreateFromSupportSignalRequest ingests support->notification signal.
type CreateFromSupportSignalRequest struct {
	SupportID     string `json:"support_id" validate:"required,uuid4"`
	Event         string `json:"event,omitempty" validate:"omitempty,max=64"`
	Channel       string `json:"channel,omitempty" validate:"omitempty,oneof=in_app email push"`
	Title         string `json:"title,omitempty" validate:"omitempty,max=180"`
	Body          string `json:"body,omitempty" validate:"omitempty,max=4000"`
	RequestID     string `json:"request_id,omitempty" validate:"omitempty,max=128"`
	CorrelationID string `json:"correlation_id,omitempty" validate:"omitempty,max=128"`
}

// CreateFromSupportSignalResponse returns upsert-like creation result.
type CreateFromSupportSignalResponse struct {
	NotificationID string `json:"notification_id"`
	Created        bool   `json:"created"`
	Category       string `json:"category"`
	State          string `json:"state"`
}
