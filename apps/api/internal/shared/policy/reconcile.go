package policy

// ReconcileRule captures required safeguards for multi-owner flows.
type ReconcileRule string

const (
	ReconcileRuleDocumentRecoveryBeforeImplementation ReconcileRule = "document_recovery_before_implementation"
	ReconcileRulePreserveTraceAndIdempotency          ReconcileRule = "preserve_request_correlation_and_idempotency_across_boundaries"
	ReconcileRuleNoCrossModuleOwnerWrite              ReconcileRule = "no_direct_write_to_another_module_owner_table"
	ReconcileRuleExposeManualReviewOnCriticalFlows    ReconcileRule = "critical_financial_or_entitlement_flows_require_manual_review_gate"
)

var ReconcileRules = []ReconcileRule{
	ReconcileRuleDocumentRecoveryBeforeImplementation,
	ReconcileRulePreserveTraceAndIdempotency,
	ReconcileRuleNoCrossModuleOwnerWrite,
	ReconcileRuleExposeManualReviewOnCriticalFlows,
}

// ReconcileFlowReference highlights flows where reconcile is expected.
type ReconcileFlowReference struct {
	Name                 string
	BoundaryType         string
	RecoveryCompensation string
}

var CanonicalReconcileFlows = []ReconcileFlowReference{
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
}
