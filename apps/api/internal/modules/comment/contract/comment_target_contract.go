package contract

import "time"

const (
	EventCommentCreated   = "comment.created"
	EventCommentEdited    = "comment.edited"
	EventCommentDeleted   = "comment.deleted"
	EventCommentModerated = "comment.moderated"
)

// TargetRelation exposes comment target relation to support/moderation surfaces.
type TargetRelation struct {
	CommentID        string
	TargetType       string
	TargetID         string
	ParentCommentID  *string
	RootCommentID    *string
	ModerationStatus string
	Deleted          bool
	UpdatedAt        time.Time
}

// ModerationSignal is comment moderation/event signal payload.
type ModerationSignal struct {
	Event         string
	CommentID     string
	TargetType    string
	TargetID      string
	OccurredAt    time.Time
	RequestID     string
	CorrelationID string
}
