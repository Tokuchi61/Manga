package dto

import "time"

// OpenThreadRequest opens or resolves thread with target user.
type OpenThreadRequest struct {
	ActorUserID  string `json:"-" validate:"required,uuid4"`
	TargetUserID string `json:"target_user_id" validate:"required,uuid4"`
}

// SendMessageRequest sends message into thread.
type SendMessageRequest struct {
	ActorUserID   string `json:"-" validate:"required,uuid4"`
	ThreadID      string `json:"-" validate:"required,uuid4"`
	Body          string `json:"body" validate:"required,max=4000"`
	RequestID     string `json:"request_id,omitempty" validate:"omitempty,max=128"`
	CorrelationID string `json:"correlation_id,omitempty" validate:"omitempty,max=128"`
}

// MarkThreadReadRequest marks own unread state as read.
type MarkThreadReadRequest struct {
	ActorUserID string `json:"-" validate:"required,uuid4"`
	ThreadID    string `json:"-" validate:"required,uuid4"`
}

// ListThreadsRequest resolves own thread list.
type ListThreadsRequest struct {
	ActorUserID string `json:"-" validate:"required,uuid4"`
}

// ListThreadMessagesRequest resolves own thread messages.
type ListThreadMessagesRequest struct {
	ActorUserID string `json:"-" validate:"required,uuid4"`
	ThreadID    string `json:"-" validate:"required,uuid4"`
	Limit       int    `json:"-" validate:"omitempty,min=1,max=200"`
	Offset      int    `json:"-" validate:"omitempty,min=0"`
	SortBy      string `json:"-" validate:"omitempty,oneof=newest oldest"`
}

// ThreadItemResponse is thread list payload.
type ThreadItemResponse struct {
	ThreadID      string    `json:"thread_id"`
	PeerUserID    string    `json:"peer_user_id"`
	LastMessageID string    `json:"last_message_id,omitempty"`
	UnreadCount   int       `json:"unread_count"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// MessageItemResponse is message payload.
type MessageItemResponse struct {
	MessageID     string    `json:"message_id"`
	ThreadID      string    `json:"thread_id"`
	SenderUserID  string    `json:"sender_user_id"`
	Body          string    `json:"body"`
	RequestID     string    `json:"request_id,omitempty"`
	CorrelationID string    `json:"correlation_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// OpenThreadResponse returns thread open operation result.
type OpenThreadResponse struct {
	ThreadID   string `json:"thread_id"`
	PeerUserID string `json:"peer_user_id"`
	Created    bool   `json:"created"`
}

// SendMessageResponse returns send operation result.
type SendMessageResponse struct {
	MessageID string `json:"message_id"`
	ThreadID  string `json:"thread_id"`
	Created   bool   `json:"created"`
}

// ListThreadsResponse wraps thread list payload.
type ListThreadsResponse struct {
	Items []ThreadItemResponse `json:"items"`
	Count int                  `json:"count"`
}

// ListThreadMessagesResponse wraps message list payload.
type ListThreadMessagesResponse struct {
	Items []MessageItemResponse `json:"items"`
	Count int                   `json:"count"`
}

// OperationResponse is generic social operation result payload.
type OperationResponse struct {
	Status string `json:"status"`
}
