package settings

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAudienceAndScopeVocabulariesMatchCanonicalValues(t *testing.T) {
	require.Equal(t, []Audience{"all", "guest", "authenticated", "authenticated_non_vip", "vip"}, AllAudiences)
	require.Equal(t, []ScopeKind{"site", "module", "feature", "resource/context"}, AllScopeKinds)
}

func TestBehaviorAndResponseVocabulariesMatchCanonicalValues(t *testing.T) {
	require.Equal(t, []DisabledBehavior{
		"visibility_off",
		"read_only",
		"write_off",
		"intake_pause",
		"read_only_intake",
		"attachment_off",
		"preview_off",
		"benefit_pause",
	}, AllDisabledBehaviors)

	require.Equal(t, []ErrorResponsePolicy{
		"not_found",
		"forbidden",
		"rate_limited",
		"validation_error",
		"temporarily_unavailable",
	}, AllErrorResponsePolicies)

	require.Equal(t, []EntitlementImpactPolicy{"none", "freeze_on_system_disable"}, AllEntitlementImpactPolicies)
}

func TestInterpretationOrderAndKillSwitchLevels(t *testing.T) {
	require.Equal(t, []string{"global", "module", "surface", "audience", "entitlement", "action", "rate_limit"}, AccessInterpretationOrder)
	require.Equal(t, []KillSwitchLevel{
		"global",
		"module-level",
		"surface-level",
		"action-level",
		"intake-only",
		"write-only",
		"external-integration-only",
	}, AllKillSwitchLevels)
}

func TestSettingsRecordFieldsMatchCanonicalSchema(t *testing.T) {
	require.Equal(t, []string{
		"key",
		"description",
		"category",
		"owner_module",
		"consumer_layer",
		"value_type",
		"default_value",
		"allowed_range_or_enum",
		"scope_kind",
		"scope_selector",
		"audience_kind",
		"audience_selector",
		"sensitive",
		"apply_mode",
		"cache_strategy",
		"schedule_support",
		"audit_required",
		"affected_surfaces",
		"disabled_behavior",
		"error_response_policy",
		"entitlement_impact_policy",
		"status",
		"notes",
	}, SettingsRecordFields)
}

func TestKeyGrammarValidators(t *testing.T) {
	require.True(t, IsFeatureToggleKey("feature.comment.write.enabled"))
	require.True(t, IsFeatureToggleKey("feature.history.bookmark_write.enabled"))
	require.False(t, IsFeatureToggleKey("comment.write.enabled"))

	require.True(t, IsModuleSurfaceMetricKey("mission.daily.reset_hour_utc"))
	require.True(t, IsModuleSurfaceMetricKey("payment.callback.intake.paused"))
	require.False(t, IsModuleSurfaceMetricKey("site.global.kill_switch"))

	require.True(t, IsSiteMetricKey("site.access.global_kill_switch"))
	require.False(t, IsSiteMetricKey("feature.site.access.enabled"))
}

func TestSettingRecordValidationHelpers(t *testing.T) {
	record := SettingRecord{
		Key:           "feature.comment.write.enabled",
		OwnerModule:   "comment",
		ConsumerLayer: "access",
		Status:        "planned",
	}
	require.True(t, record.HasRequiredIdentity())
	require.True(t, record.IsStatusKnown())

	record.Status = "invalid"
	require.False(t, record.IsStatusKnown())
}
