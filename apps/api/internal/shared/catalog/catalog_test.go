package catalog

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuditEventTypesMatchCanonicalList(t *testing.T) {
	require.Equal(t, []AuditEventType{
		"security_auth",
		"access_policy",
		"admin_action",
		"moderation_action",
		"support_case",
		"payment_financial",
		"inventory_change",
		"shop_purchase",
		"user_state",
		"notification_ops",
		"ads_ops",
		"system_ops",
	}, AllAuditEventTypes)
	require.True(t, IsValidAuditEventType("system_ops"))
	require.False(t, IsValidAuditEventType("unknown"))
}

func TestModerationStatusesMatchCanonicalLists(t *testing.T) {
	require.Equal(t, []ModerationCaseStatus{
		"new",
		"queued",
		"assigned",
		"in_review",
		"escalated",
		"resolved",
		"rejected",
		"closed",
	}, AllModerationCaseStatuses)

	require.Equal(t, []ModerationAssignmentStatus{
		"unassigned",
		"assigned",
		"handoff_pending",
		"released",
	}, AllModerationAssignmentStatuses)

	require.Equal(t, []ModerationActionResult{
		"none",
		"content_hidden",
		"content_restored",
		"warning_sent",
		"no_action",
	}, AllModerationActionResults)
}

func TestNotificationCategoriesMatchCanonicalList(t *testing.T) {
	require.Equal(t, []NotificationCategory{
		"account_security",
		"social",
		"comment",
		"support",
		"moderation",
		"mission",
		"royalpass",
		"shop",
		"payment",
		"system_ops",
	}, AllNotificationCategories)
}

func TestPolicyEffectsMatchCanonicalList(t *testing.T) {
	require.Equal(t, []PolicyEffect{
		"allow",
		"deny",
		"deny_soft",
		"require_auth",
		"require_role",
		"require_entitlement",
		"read_only",
		"write_off",
		"mask",
		"needs_review",
	}, AllPolicyEffects)
}

func TestPurchaseAndRewardSourceTypesMatchCanonicalLists(t *testing.T) {
	require.Equal(t, []PurchaseSourceType{
		"catalog_purchase",
		"premium_activation",
		"mana_wallet",
		"external_provider",
		"recovery_replay",
		"admin_issue",
		"gift_code",
	}, AllPurchaseSourceTypes)

	require.Equal(t, []RewardSourceType{
		"mission",
		"royalpass",
		"shop",
		"admin_grant",
		"compensation",
		"seasonal_event",
		"referral",
		"reconciliation_repair",
	}, AllRewardSourceTypes)
}

func TestSupportTargetAndVisibilityListsMatchCanonicalValues(t *testing.T) {
	require.Equal(t, []SupportStatus{
		"open",
		"triaged",
		"waiting_user",
		"waiting_team",
		"resolved",
		"rejected",
		"closed",
		"spam",
	}, AllSupportStatuses)

	require.Equal(t, []SupportReplyVisibility{
		"public_to_requester",
		"internal_only",
	}, AllSupportReplyVisibility)

	require.Equal(t, []TargetType{"manga", "chapter", "comment", "social"}, AllTargetTypes)
	require.Equal(t, []VisibilityState{"public", "limited", "private", "hidden", "removed", "archived"}, AllVisibilityStates)
}
