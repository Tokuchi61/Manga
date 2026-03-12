package entity

import "time"

// SupportKind defines support intake surface.
type SupportKind string

const (
	SupportKindCommunication SupportKind = "communication"
	SupportKindTicket        SupportKind = "ticket"
	SupportKindReport        SupportKind = "report"
)

// SupportStatus defines support lifecycle states.
type SupportStatus string

const (
	SupportStatusOpen        SupportStatus = "open"
	SupportStatusTriaged     SupportStatus = "triaged"
	SupportStatusWaitingUser SupportStatus = "waiting_user"
	SupportStatusWaitingTeam SupportStatus = "waiting_team"
	SupportStatusResolved    SupportStatus = "resolved"
	SupportStatusRejected    SupportStatus = "rejected"
	SupportStatusClosed      SupportStatus = "closed"
	SupportStatusSpam        SupportStatus = "spam"
)

// SupportPriority defines support queue priority.
type SupportPriority string

const (
	SupportPriorityLow    SupportPriority = "low"
	SupportPriorityNormal SupportPriority = "normal"
	SupportPriorityHigh   SupportPriority = "high"
	SupportPriorityUrgent SupportPriority = "urgent"
)

// SupportTargetType defines canonical support report targets.
type SupportTargetType string

const (
	SupportTargetTypeManga   SupportTargetType = "manga"
	SupportTargetTypeChapter SupportTargetType = "chapter"
	SupportTargetTypeComment SupportTargetType = "comment"
)

// ReplyVisibility defines requester-visible/internal support reply mode.
type ReplyVisibility string

const (
	ReplyVisibilityPublicToRequester ReplyVisibility = "public_to_requester"
	ReplyVisibilityInternalOnly      ReplyVisibility = "internal_only"
)

// SupportReply is support case reply/message record.
type SupportReply struct {
	ID            string
	SupportID     string
	AuthorUserID  string
	Message       string
	SanitizedBody string
	Visibility    ReplyVisibility
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// SupportCase is the owner aggregate for stage-10 support flows.
type SupportCase struct {
	ID                           string
	RequesterUserID              string
	Kind                         SupportKind
	Category                     string
	Priority                     SupportPriority
	ReasonCode                   string
	ReasonText                   string
	TargetType                   *SupportTargetType
	TargetID                     *string
	Status                       SupportStatus
	DuplicateOfSupportID         *string
	RequestID                    string
	SpamRiskScore                int
	Attachments                  []string
	Replies                      []SupportReply
	ResolutionNote               string
	AssigneeUserID               *string
	ReviewedByUserID             *string
	ResolvedAt                   *time.Time
	ClosedAt                     *time.Time
	ModerationHandoffRequestedAt *time.Time
	LinkedModerationCaseID       *string
	CreatedAt                    time.Time
	UpdatedAt                    time.Time
}
