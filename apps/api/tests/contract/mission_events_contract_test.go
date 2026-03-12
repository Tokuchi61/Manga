package contract_test

import (
	"testing"

	missionevents "github.com/Tokuchi61/Manga/apps/api/internal/modules/mission/events"
	"github.com/stretchr/testify/require"
)

func TestMissionEventConstants(t *testing.T) {
	require.Equal(t, "mission.progressed", missionevents.EventMissionProgressed)
	require.Equal(t, "mission.completed", missionevents.EventMissionCompleted)
	require.Equal(t, "mission.claim.requested", missionevents.EventMissionClaimRequested)
	require.Equal(t, "mission.reset", missionevents.EventMissionReset)
}
