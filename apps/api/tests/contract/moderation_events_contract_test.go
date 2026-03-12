package contract_test

import (
	"testing"

	moderationevents "github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/events"
	"github.com/stretchr/testify/require"
)

func TestModerationEventConstants(t *testing.T) {
	require.Equal(t, "moderation.case.created", moderationevents.EventModerationCaseCreated)
	require.Equal(t, "moderation.case.assigned", moderationevents.EventModerationCaseAssigned)
	require.Equal(t, "moderation.case.escalated", moderationevents.EventModerationCaseEscalated)
	require.Equal(t, "moderation.action.applied", moderationevents.EventModerationActionApplied)
}
