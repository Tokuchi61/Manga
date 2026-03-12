package contract

import "time"

const EventSupportModerationHandoffRequested = "support.moderation_handoff_requested"

// ModerationHandoffReference exposes support -> moderation linked case contract.
type ModerationHandoffReference struct {
	SupportID     string
	SupportKind   string
	TargetType    string
	TargetID      string
	ReasonCode    string
	RequestedAt   time.Time
	RequestID     string
	CorrelationID string
}
