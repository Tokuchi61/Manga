package dto

// UpdateCommentRequest updates editable comment fields.
type UpdateCommentRequest struct {
	CommentID   string    `json:"-" validate:"required,uuid4"`
	ActorUserID string    `json:"-" validate:"required,uuid4"`
	Content     *string   `json:"content,omitempty" validate:"omitempty,min=1,max=5000"`
	Attachments *[]string `json:"attachments,omitempty" validate:"omitempty,max=10,dive,url,max=1500"`
	Spoiler     *bool     `json:"spoiler,omitempty"`
}

// DeleteCommentRequest soft-deletes comment content.
type DeleteCommentRequest struct {
	CommentID   string `json:"-" validate:"required,uuid4"`
	ActorUserID string `json:"-" validate:"required,uuid4"`
	Reason      string `json:"reason,omitempty" validate:"omitempty,max=500"`
}

// RestoreCommentRequest restores a soft-deleted comment.
type RestoreCommentRequest struct {
	CommentID   string `json:"-" validate:"required,uuid4"`
	ActorUserID string `json:"-" validate:"required,uuid4"`
}

// UpdateModerationRequest updates moderation/spoiler/pin/lock fields.
type UpdateModerationRequest struct {
	CommentID        string  `json:"-" validate:"required,uuid4"`
	ActorUserID      string  `json:"-" validate:"required,uuid4"`
	ModerationStatus *string `json:"moderation_status,omitempty" validate:"omitempty,oneof=visible hidden flagged"`
	Pinned           *bool   `json:"pinned,omitempty"`
	Locked           *bool   `json:"locked,omitempty"`
	Shadowbanned     *bool   `json:"shadowbanned,omitempty"`
	Spoiler          *bool   `json:"spoiler,omitempty"`
}

// OperationResponse is generic operation result payload.
type OperationResponse struct {
	Status string `json:"status"`
}
