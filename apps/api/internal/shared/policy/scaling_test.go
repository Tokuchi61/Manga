package policy

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestModuleInventoryAndDomainGroupPolicies(t *testing.T) {
	require.Equal(t, []ModuleStatus{"planned", "active", "deprecated", "archived"}, AllModuleStatuses)
	require.Equal(t, []string{"identity", "content", "community", "operations", "engagement", "commerce", "gameplay"}, SuggestedDomainGroups)
	require.Equal(t, []string{"name", "domain_group_optional", "description", "status", "doc_path"}, ModuleInventoryRequiredFields)

	require.True(t, IsValidModuleStatus("active"))
	require.False(t, IsValidModuleStatus("unknown"))
	require.True(t, IsSuggestedDomainGroup("content"))
	require.False(t, IsSuggestedDomainGroup("unknown"))

	require.Equal(t, "apps/api/internal/modules/auth", BuildModuleRootPath("auth", ""))
	require.Equal(t, "apps/api/internal/modules/content/manga", BuildModuleRootPath("manga", "content"))

	record := ModuleInventoryRecord{
		Name:        "manga",
		DomainGroup: "content",
		Description: "Manga owner module",
		Status:      ModuleStatusActive,
		DocPath:     "docs/modules/content/manga.md",
	}
	require.True(t, record.HasRequiredFields())
	require.True(t, record.IsDocPathUnderDocsModules())

	record.DocPath = "docs/modules.md"
	require.True(t, record.IsDocPathUnderDocsModules())

	record.DocPath = "docs/shared.md"
	require.False(t, record.IsDocPathUnderDocsModules())
}

func TestReadModelAndProjectionRules(t *testing.T) {
	require.Equal(t, []EventualConsistencyWindow{"short", "medium"}, SupportedConsistencyWindows)
	require.True(t, IsSupportedConsistencyWindow("short"))
	require.False(t, IsSupportedConsistencyWindow("long"))

	require.Equal(t, []ProjectionImplementationRule{
		"owner_keeps_canonical_write_model",
		"feed_by_event_outbox_or_projection_contract",
		"consistency_window_must_be_documented",
		"rebuild_and_replay_must_be_planned_from_start",
		"lag_and_errors_must_be_observable",
		"rebuild_must_be_idempotent",
		"replay_payload_requires_schema_request_and_correlation",
		"no_direct_cross_module_owner_table_write_for_counters",
		"event_producers_align_with_transactional_outbox",
	}, ProjectionImplementationRules)
}

func TestReportingAndMaintenanceRules(t *testing.T) {
	require.Equal(t, []ReportingLayer{"operational_summary", "analytics_aggregate", "export_query_layer"}, ReportingLayers)
	require.Len(t, ReportingRules, 4)

	require.Equal(t, []MaintenanceRefactorRule{
		"no_unnecessary_refactor",
		"changes_must_be_logical_small_and_revertible",
		"update_docs_when_architecture_or_boundaries_change",
		"avoid_cross_module_tight_coupling",
		"keep_owner_boundaries_clear",
	}, MaintenanceRefactorRules)
	require.Equal(t, []string{
		"stage_tests_must_pass",
		"docker_build_and_run_must_pass",
		"versioning_must_be_applied_before_push",
		"changelog_and_upgrade_docs_must_be_updated",
	}, StageThreeOperationalChecklist)
}

func TestReconcileRulesAndFlowReferences(t *testing.T) {
	require.Equal(t, []ReconcileRule{
		"document_recovery_before_implementation",
		"preserve_request_correlation_and_idempotency_across_boundaries",
		"no_direct_write_to_another_module_owner_table",
		"critical_financial_or_entitlement_flows_require_manual_review_gate",
	}, ReconcileRules)

	require.Equal(t, []ReconcileFlowReference{
		{
			Name:                 "shop_purchase_to_payment_callback_to_inventory_grant",
			BoundaryType:         "multi_step_event_driven",
			RecoveryCompensation: "reconcile_duplicate_guard_and_grant_retry",
		},
		{
			Name:                 "royalpass_premium_activation",
			BoundaryType:         "multi_step_event_driven",
			RecoveryCompensation: "activation_reference_reconcile_and_replay",
		},
	}, CanonicalReconcileFlows)
}
