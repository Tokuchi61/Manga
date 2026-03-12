package contract_test

import (
	"testing"
	"time"

	usercontract "github.com/Tokuchi61/Manga/apps/api/internal/modules/user/contract"
	"github.com/stretchr/testify/require"
)

func TestUserAccessContractShape(t *testing.T) {
	now := time.Date(2026, 3, 12, 10, 0, 0, 0, time.UTC)
	endAt := now.Add(30 * 24 * time.Hour)

	signal := usercontract.AccessSignal{
		UserID:                      "user-1",
		AccountState:                "active",
		ProfileVisibility:           "public",
		HistoryVisibilityPreference: "private",
		VIPActive:                   true,
		VIPFrozen:                   false,
		VIPEndsAt:                   &endAt,
		UpdatedAt:                   now,
	}

	require.Equal(t, "active", signal.AccountState)
	require.Equal(t, "public", signal.ProfileVisibility)
	require.Equal(t, "private", signal.HistoryVisibilityPreference)
	require.True(t, signal.VIPActive)
	require.NotNil(t, signal.VIPEndsAt)
}
