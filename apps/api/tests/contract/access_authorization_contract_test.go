package contract_test

import (
	"testing"
	"time"

	accesscontract "github.com/Tokuchi61/Manga/apps/api/internal/modules/access/contract"
	authcontract "github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/contract"
	usercontract "github.com/Tokuchi61/Manga/apps/api/internal/modules/user/contract"
	"github.com/stretchr/testify/require"
)

func TestAccessAuthorizationContractShape(t *testing.T) {
	now := time.Date(2026, 3, 12, 11, 0, 0, 0, time.UTC)
	input := accesscontract.AuthorizationInput{
		UserID:              "8fb89911-2a20-4657-bf6e-53e5b7cbf56f",
		Permission:          accesscontract.PermissionHistoryTimelineReadOwn,
		FeatureKey:          "feature.history.timeline.enabled",
		ResourceOwnerUserID: "8fb89911-2a20-4657-bf6e-53e5b7cbf56f",
		Identity: authcontract.VerifiedIdentity{
			CredentialID:    "cred-1",
			SessionID:       "sess-1",
			EmailVerified:   true,
			AuthenticatedAt: now,
		},
		UserSignal: usercontract.AccessSignal{
			UserID:                      "8fb89911-2a20-4657-bf6e-53e5b7cbf56f",
			AccountState:                "active",
			ProfileVisibility:           "public",
			HistoryVisibilityPreference: "private",
			VIPActive:                   true,
			VIPFrozen:                   false,
			UpdatedAt:                   now,
		},
	}

	require.Equal(t, accesscontract.PermissionHistoryTimelineReadOwn, input.Permission)
	require.Equal(t, "active", input.UserSignal.AccountState)
	require.True(t, input.Identity.EmailVerified)
}

func TestAccessCanonicalPermissionExamples(t *testing.T) {
	require.Contains(t, accesscontract.CanonicalPermissions, accesscontract.PermissionHistoryContinueReadingOwn)
	require.Contains(t, accesscontract.CanonicalPermissions, accesscontract.PermissionHistoryTimelineReadOwn)
	require.Contains(t, accesscontract.CanonicalPermissions, accesscontract.PermissionHistoryLibraryReadOwn)
	require.Contains(t, accesscontract.CanonicalPermissions, accesscontract.PermissionHistoryBookmarkWriteOwn)
	require.Contains(t, accesscontract.CanonicalPermissions, accesscontract.PermissionHistoryLibraryReadPublic)
	require.Contains(t, accesscontract.CanonicalPermissions, accesscontract.PermissionMangaDiscoveryView)
	require.Contains(t, accesscontract.CanonicalPermissions, accesscontract.PermissionAdsView)
	require.Contains(t, accesscontract.CanonicalPermissions, accesscontract.PermissionShopItemPurchase)
	require.Contains(t, accesscontract.CanonicalPermissions, accesscontract.PermissionPaymentManaPurchase)
	require.Contains(t, accesscontract.CanonicalPermissions, accesscontract.PermissionPaymentTransactionReadOwn)
}
