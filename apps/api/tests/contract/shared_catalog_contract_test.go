package contract

import (
	"testing"

	"github.com/Tokuchi61/Manga/apps/api/internal/shared/catalog"
	"github.com/stretchr/testify/require"
)

func TestSharedCatalogContractKeysRemainStable(t *testing.T) {
	require.Contains(t, catalog.AllAuditEventTypes, catalog.AuditEventSystemOps)
	require.Contains(t, catalog.AllNotificationCategories, catalog.NotificationSystemOps)
	require.Contains(t, catalog.AllPolicyEffects, catalog.PolicyEffectDenySoft)
	require.Contains(t, catalog.AllPurchaseSourceTypes, catalog.PurchaseSourceExternalProvider)
	require.Contains(t, catalog.AllRewardSourceTypes, catalog.RewardSourceReconciliationRepair)
	require.Contains(t, catalog.AllSupportStatuses, catalog.SupportStatusWaitingTeam)
	require.Contains(t, catalog.AllVisibilityStates, catalog.VisibilityStateArchived)
}
