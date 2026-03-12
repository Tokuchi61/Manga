package contract_test

import (
	"testing"
	"time"

	supportcontract "github.com/Tokuchi61/Manga/apps/api/internal/modules/support/contract"
	"github.com/stretchr/testify/require"
)

func TestSupportNotificationSignalContractShape(t *testing.T) {
	now := time.Date(2026, 3, 12, 23, 31, 0, 0, time.UTC)
	signal := supportcontract.NotificationSignal{
		Event:           supportcontract.EventSupportCreated,
		SupportID:       "307cc4f9-30f3-4fd8-9544-1f8f6507dbb7",
		RequesterUserID: "dd30f1d6-2dc8-4ac4-a6f4-0e4070ed4b2b",
		Status:          "open",
		OccurredAt:      now,
		RequestID:       "req-stage10",
		CorrelationID:   "corr-stage10",
	}

	require.Equal(t, supportcontract.EventSupportCreated, signal.Event)
	require.Equal(t, "open", signal.Status)
	require.Equal(t, now, signal.OccurredAt)
	require.Equal(t, "req-stage10", signal.RequestID)
	require.Equal(t, "corr-stage10", signal.CorrelationID)
}

func TestSupportNotificationEventConstants(t *testing.T) {
	require.Equal(t, "support.created", supportcontract.EventSupportCreated)
	require.Equal(t, "support.replied", supportcontract.EventSupportReplied)
	require.Equal(t, "support.resolved", supportcontract.EventSupportResolved)
}
