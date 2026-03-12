package entity

import "time"

// CaseSource defines moderation intake origin.
type CaseSource string

const (
	CaseSourceSupportReport CaseSource = "support_report"
	CaseSourceManual        CaseSource = "manual"
)

// TargetType defines canonical moderation target set.
type TargetType string

const (
	TargetTypeManga   TargetType = "manga"
	TargetTypeChapter TargetType = "chapter"
	TargetTypeComment TargetType = "comment"
)

// CaseStatus defines moderation case lifecycle.
type CaseStatus string

const (
	CaseStatusNew       CaseStatus = "new"
	CaseStatusQueued    CaseStatus = "queued"
	CaseStatusAssigned  CaseStatus = "assigned"
	CaseStatusInReview  CaseStatus = "in_review"
	CaseStatusEscalated CaseStatus = "escalated"
	CaseStatusResolved  CaseStatus = "resolved"
	CaseStatusRejected  CaseStatus = "rejected"
	CaseStatusClosed    CaseStatus = "closed"
)

// AssignmentStatus defines moderator assignment lifecycle.
type AssignmentStatus string

const (
	AssignmentStatusUnassigned    AssignmentStatus = "unassigned"
	AssignmentStatusAssigned      AssignmentStatus = "assigned"
	AssignmentStatusHandoffPending AssignmentStatus = "handoff_pending"
	AssignmentStatusReleased      AssignmentStatus = "released"
)

// EscalationStatus defines escalation lifecycle.
type EscalationStatus string

const (
	EscalationStatusNotEscalated EscalationStatus = "not_escalated"
	EscalationStatusPendingAdmin EscalationStatus = "pending_admin"
	EscalationStatusEscalated    EscalationStatus = "escalated"
	EscalationStatusResolved     EscalationStatus = "resolved"
)

// ActionType defines supported moderation actions.
type ActionType string

const (
	ActionTypeHide           ActionType = "hide"
	ActionTypeUnhide         ActionType = "unhide"
	ActionTypeLock           ActionType = "lock"
	ActionTypeUnlock         ActionType = "unlock"
	ActionTypeWarning        ActionType = "warning"
	ActionTypeReviewComplete ActionType = "review_complete"
	ActionTypeEscalate       ActionType = "escalate"
)

// ActionResult defines action outcome summary.
type ActionResult string

const (
	ActionResultNone            ActionResult = "none"
	ActionResultContentHidden   ActionResult = "content_hidden"
	ActionResultContentRestored ActionResult = "content_restored"
	ActionResultWarningSent     ActionResult = "warning_sent"
	ActionResultNoAction        ActionResult = "no_action"
)

// ModeratorNote stores moderator internal timeline notes.
type ModeratorNote struct {
	ID           string
	CaseID        string
	AuthorUserID string
	Body          string
	InternalOnly  bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// CaseAction stores applied moderation action timeline.
type CaseAction struct {
	ID           string
	CaseID        string
	ActorUserID  string
	ActionType   ActionType
	ReasonCode   string
	Summary      string
	ActionResult ActionResult
	CreatedAt    time.Time
}

// Case is the owner aggregate for stage-11 moderation flows.
type Case struct {
	ID                      string
	Source                  CaseSource
	SourceRefID             *string
	RequestID               string
	CorrelationID           string
	TargetType              TargetType
	TargetID                string
	ReporterUserID          *string
	Status                  CaseStatus
	AssignmentStatus        AssignmentStatus
	AssignedModeratorUserID *string
	EscalationStatus        EscalationStatus
	EscalationReason        string
	EscalatedAt             *time.Time
	ActionResult            ActionResult
	LastActionAt            *time.Time
	Notes                   []ModeratorNote
	Actions                 []CaseAction
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

func (c Case) IsTerminal() bool {
	switch c.Status {
	case CaseStatusResolved, CaseStatusRejected, CaseStatusClosed:
		return true
	default:
		return false
	}
}
