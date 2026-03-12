package entity

import "time"

// TargetType defines canonical comment targets.
type TargetType string

const (
	TargetTypeManga   TargetType = "manga"
	TargetTypeChapter TargetType = "chapter"
)

// ModerationStatus defines comment visibility lifecycle.
type ModerationStatus string

const (
	ModerationStatusVisible ModerationStatus = "visible"
	ModerationStatusHidden  ModerationStatus = "hidden"
	ModerationStatusFlagged ModerationStatus = "flagged"
)

// Comment is the owner aggregate for stage-9 comment/thread flows.
type Comment struct {
	ID               string
	TargetType       TargetType
	TargetID         string
	AuthorUserID     string
	ParentCommentID  *string
	RootCommentID    *string
	Depth            int
	Content          string
	SanitizedContent string
	Attachments      []string
	Spoiler          bool
	Pinned           bool
	Locked           bool
	ModerationStatus ModerationStatus
	Shadowbanned     bool
	SpamRiskScore    int
	LikeCount        int64
	ReplyCount       int64
	EditCount        int
	EditedAt         *time.Time
	DeletedAt        *time.Time
	DeleteReason     string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
