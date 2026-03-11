package catalog

// SupportStatus defines canonical support case lifecycle values.
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

var AllSupportStatuses = []SupportStatus{
	SupportStatusOpen,
	SupportStatusTriaged,
	SupportStatusWaitingUser,
	SupportStatusWaitingTeam,
	SupportStatusResolved,
	SupportStatusRejected,
	SupportStatusClosed,
	SupportStatusSpam,
}

func IsValidSupportStatus(value SupportStatus) bool {
	switch value {
	case SupportStatusOpen,
		SupportStatusTriaged,
		SupportStatusWaitingUser,
		SupportStatusWaitingTeam,
		SupportStatusResolved,
		SupportStatusRejected,
		SupportStatusClosed,
		SupportStatusSpam:
		return true
	default:
		return false
	}
}

// SupportReplyVisibility defines canonical support reply visibility values.
type SupportReplyVisibility string

const (
	SupportReplyPublicToRequester SupportReplyVisibility = "public_to_requester"
	SupportReplyInternalOnly      SupportReplyVisibility = "internal_only"
)

var AllSupportReplyVisibility = []SupportReplyVisibility{
	SupportReplyPublicToRequester,
	SupportReplyInternalOnly,
}

func IsValidSupportReplyVisibility(value SupportReplyVisibility) bool {
	switch value {
	case SupportReplyPublicToRequester, SupportReplyInternalOnly:
		return true
	default:
		return false
	}
}
