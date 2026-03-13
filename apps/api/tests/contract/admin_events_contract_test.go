package contract_test

import (
	"testing"

	adminevents "github.com/Tokuchi61/Manga/apps/api/internal/modules/admin/events"
	"github.com/stretchr/testify/require"
)

func TestAdminEventConstants(t *testing.T) {
	require.Equal(t, "admin.setting.changed", adminevents.EventAdminSettingChanged)
	require.Equal(t, "admin.override.applied", adminevents.EventAdminOverrideApplied)
	require.Equal(t, "admin.user.reviewed", adminevents.EventAdminUserReviewed)
}
