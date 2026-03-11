package catalog

// Moderation case lifecycle states.
type ModerationCaseStatus string

const (
	ModerationCaseNew       ModerationCaseStatus = "new"
	ModerationCaseQueued    ModerationCaseStatus = "queued"
	ModerationCaseAssigned  ModerationCaseStatus = "assigned"
	ModerationCaseInReview  ModerationCaseStatus = "in_review"
	ModerationCaseEscalated ModerationCaseStatus = "escalated"
	ModerationCaseResolved  ModerationCaseStatus = "resolved"
	ModerationCaseRejected  ModerationCaseStatus = "rejected"
	ModerationCaseClosed    ModerationCaseStatus = "closed"
)

var AllModerationCaseStatuses = []ModerationCaseStatus{
	ModerationCaseNew,
	ModerationCaseQueued,
	ModerationCaseAssigned,
	ModerationCaseInReview,
	ModerationCaseEscalated,
	ModerationCaseResolved,
	ModerationCaseRejected,
	ModerationCaseClosed,
}

func IsValidModerationCaseStatus(value ModerationCaseStatus) bool {
	switch value {
	case ModerationCaseNew,
		ModerationCaseQueued,
		ModerationCaseAssigned,
		ModerationCaseInReview,
		ModerationCaseEscalated,
		ModerationCaseResolved,
		ModerationCaseRejected,
		ModerationCaseClosed:
		return true
	default:
		return false
	}
}

// Moderation assignment states.
type ModerationAssignmentStatus string

const (
	ModerationAssignmentUnassigned ModerationAssignmentStatus = "unassigned"
	ModerationAssignmentAssigned   ModerationAssignmentStatus = "assigned"
	ModerationAssignmentHandoff    ModerationAssignmentStatus = "handoff_pending"
	ModerationAssignmentReleased   ModerationAssignmentStatus = "released"
)

var AllModerationAssignmentStatuses = []ModerationAssignmentStatus{
	ModerationAssignmentUnassigned,
	ModerationAssignmentAssigned,
	ModerationAssignmentHandoff,
	ModerationAssignmentReleased,
}

func IsValidModerationAssignmentStatus(value ModerationAssignmentStatus) bool {
	switch value {
	case ModerationAssignmentUnassigned,
		ModerationAssignmentAssigned,
		ModerationAssignmentHandoff,
		ModerationAssignmentReleased:
		return true
	default:
		return false
	}
}

// Moderation decision outcomes.
type ModerationActionResult string

const (
	ModerationActionNone           ModerationActionResult = "none"
	ModerationActionContentHidden  ModerationActionResult = "content_hidden"
	ModerationActionContentRestore ModerationActionResult = "content_restored"
	ModerationActionWarningSent    ModerationActionResult = "warning_sent"
	ModerationActionNoAction       ModerationActionResult = "no_action"
)

var AllModerationActionResults = []ModerationActionResult{
	ModerationActionNone,
	ModerationActionContentHidden,
	ModerationActionContentRestore,
	ModerationActionWarningSent,
	ModerationActionNoAction,
}

func IsValidModerationActionResult(value ModerationActionResult) bool {
	switch value {
	case ModerationActionNone,
		ModerationActionContentHidden,
		ModerationActionContentRestore,
		ModerationActionWarningSent,
		ModerationActionNoAction:
		return true
	default:
		return false
	}
}
