package policy

// TransactionFlowReference stores canonical boundary examples.
type TransactionFlowReference struct {
	Name                 string
	BoundaryType         string
	RecoveryCompensation string
}

var TransactionSelectionRules = []string{
	"single_owner_single_db_aggregate_flows_should_use_single_transaction",
	"external_provider_queue_or_cross_owner_flows_must_use_multi_step_idempotent_orchestration",
	"avoid_db_write_plus_immediate_publish_use_transactional_outbox",
	"recovery_or_compensation_must_be_documented_before_implementation",
	"cross_boundary_hops_must_preserve_request_id_correlation_id_and_idempotency_key",
	"financial_or_user_entitlement_flows_must_expose_manual_review_gate",
}

var TransactionFlowReferences = []TransactionFlowReference{
	{
		Name:                 "auth_login_and_session_create",
		BoundaryType:         "single_transaction",
		RecoveryCompensation: "session_revoke_and_security_audit",
	},
	{
		Name:                 "shop_purchase_to_payment_callback_to_inventory_grant",
		BoundaryType:         "multi_step_event_driven",
		RecoveryCompensation: "reconcile_duplicate_guard_and_grant_retry",
	},
	{
		Name:                 "mission_complete_to_reward_grant",
		BoundaryType:         "coordinated_idempotent",
		RecoveryCompensation: "claim_replay_and_grant_dedup",
	},
	{
		Name:                 "support_report_to_moderation_case",
		BoundaryType:         "sync_or_async_policy_based",
		RecoveryCompensation: "linked_case_reference_retry",
	},
	{
		Name:                 "notification_create_to_channel_delivery",
		BoundaryType:         "write_plus_async_delivery",
		RecoveryCompensation: "backoff_suppression_and_dead_letter",
	},
	{
		Name:                 "royalpass_premium_activation",
		BoundaryType:         "multi_step_event_driven",
		RecoveryCompensation: "activation_reference_reconcile_and_replay",
	},
}
