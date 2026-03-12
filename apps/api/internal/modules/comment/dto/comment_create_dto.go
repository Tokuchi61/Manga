package dto

// CreateCommentRequest creates a new root/reply comment.
type CreateCommentRequest struct {
	TargetType      string   `json:"target_type" validate:"required,oneof=manga chapter"`
	TargetID        string   `json:"target_id" validate:"required,uuid4"`
	AuthorUserID    string   `json:"-" validate:"required,uuid4"`
	ParentCommentID *string  `json:"parent_comment_id,omitempty" validate:"omitempty,uuid4"`
	Content         string   `json:"content" validate:"required,min=1,max=5000"`
	Attachments     []string `json:"attachments,omitempty" validate:"omitempty,max=10,dive,url,max=1500"`
	Spoiler         bool     `json:"spoiler"`
}

// CreateCommentResponse returns stable creation fields.
type CreateCommentResponse struct {
	CommentID       string  `json:"comment_id"`
	TargetType      string  `json:"target_type"`
	TargetID        string  `json:"target_id"`
	ParentCommentID *string `json:"parent_comment_id,omitempty"`
	RootCommentID   *string `json:"root_comment_id,omitempty"`
	Depth           int     `json:"depth"`
	ModerationState string  `json:"moderation_status"`
}
