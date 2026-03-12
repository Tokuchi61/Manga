package contract_test

import (
	"testing"
	"time"

	supportcontract "github.com/Tokuchi61/Manga/apps/api/internal/modules/support/contract"
	"github.com/stretchr/testify/require"
)

func TestSupportModerationHandoffReferenceContractShape(t *testing.T) {
	now := time.Date(2026, 3, 12, 23, 30, 0, 0, time.UTC)
	reference := supportcontract.ModerationHandoffReference{
		SupportID:     "307cc4f9-30f3-4fd8-9544-1f8f6507dbb7",
		SupportKind:   "report",
		TargetType:    "comment",
		TargetID:      "a6d5c0f7-ccfb-4958-ad8c-f5deb145f376",
		ReasonCode:    "abuse",
		RequestedAt:   now,
		RequestID:     "req-stage10",
		CorrelationID: "corr-stage10",
	}

	require.Equal(t, "report", reference.SupportKind)
	require.Equal(t, "comment", reference.TargetType)
	require.Equal(t, now, reference.RequestedAt)
	require.Equal(t, "req-stage10", reference.RequestID)
	require.Equal(t, "corr-stage10", reference.CorrelationID)
}

func TestSupportModerationEventConstant(t *testing.T) {
	require.Equal(t, "support.moderation_handoff_requested", supportcontract.EventSupportModerationHandoffRequested)
}
