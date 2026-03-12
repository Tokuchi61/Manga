package dto

import "time"

// ListInboxRequest defines requester-owned inbox listing controls.
type ListInboxRequest struct {
	UserID     string `json:"-" validate:"required,uuid4"`
	Category   string `json:"-" validate:"omitempty,max=64"`
	State      string `json:"-" validate:"omitempty,oneof=created delivered failed read"`
	UnreadOnly bool   `json:"-"`
	Limit      int    `json:"-" validate:"omitempty,min=1,max=200"`
	Offset     int    `json:"-" validate:"omitempty,min=0"`
	SortBy     string `json:"-" validate:"omitempty,oneof=newest oldest"`
}

// GetNotificationDetailRequest resolves notification detail.
type GetNotificationDetailRequest struct {
	UserID         string `json:"-" validate:"required,uuid4"`
	NotificationID string `json:"-" validate:"required,uuid4"`
}

// MarkReadRequest marks notification as read.
type MarkReadRequest struct {
	UserID         string `json:"-" validate:"required,uuid4"`
	NotificationID string `json:"-" validate:"required,uuid4"`
}

// NotificationListItemResponse is inbox list item payload.
type NotificationListItemResponse struct {
	NotificationID string     `json:"notification_id"`
	Category       string     `json:"category"`
	Channel        string     `json:"channel"`
	State          string     `json:"state"`
	Title          string     `json:"title"`
	Body           string     `json:"body"`
	SourceEvent    string     `json:"source_event"`
	SourceRefID    string     `json:"source_ref_id,omitempty"`
	ReadAt         *time.Time `json:"read_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// ListInboxResponse wraps inbox list payload.
type ListInboxResponse struct {
	Items      []NotificationListItemResponse `json:"items"`
	Count      int                            `json:"count"`
	UnreadCount int                           `json:"unread_count"`
}

// NotificationDetailResponse is detail payload.
type NotificationDetailResponse struct {
	NotificationID      string     `json:"notification_id"`
	UserID              string     `json:"user_id"`
	Category            string     `json:"category"`
	Channel             string     `json:"channel"`
	TemplateKey         string     `json:"template_key"`
	Title               string     `json:"title"`
	Body                string     `json:"body"`
	State               string     `json:"state"`
	SourceEvent         string     `json:"source_event"`
	SourceRefID         string     `json:"source_ref_id,omitempty"`
	RequestID           string     `json:"request_id,omitempty"`
	CorrelationID       string     `json:"correlation_id,omitempty"`
	DeliveryAttemptCount int       `json:"delivery_attempt_count"`
	LastFailureReason   string     `json:"last_failure_reason,omitempty"`
	ReadAt              *time.Time `json:"read_at,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

// OperationResponse is generic operation result payload.
type OperationResponse struct {
	Status string `json:"status"`
}
