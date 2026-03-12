package contract_test

import (
	"testing"

	rpevents "github.com/Tokuchi61/Manga/apps/api/internal/modules/royalpass/events"
	"github.com/stretchr/testify/require"
)

func TestRoyalPassEventConstants(t *testing.T) {
	require.Equal(t, "royalpass.progressed", rpevents.EventRoyalPassProgressed)
	require.Equal(t, "royalpass.claim.requested", rpevents.EventRoyalPassClaimRequested)
	require.Equal(t, "royalpass.season.started", rpevents.EventRoyalPassSeasonStarted)
	require.Equal(t, "royalpass.premium.activated", rpevents.EventRoyalPassPremiumActivated)
}
