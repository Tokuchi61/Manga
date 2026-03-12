package policy

// EventualConsistencyWindow defines accepted projection lag categories.
type EventualConsistencyWindow string

const (
	ConsistencyWindowShort  EventualConsistencyWindow = "short"
	ConsistencyWindowMedium EventualConsistencyWindow = "medium"
)

var SupportedConsistencyWindows = []EventualConsistencyWindow{
	ConsistencyWindowShort,
	ConsistencyWindowMedium,
}

// ProjectionImplementationRule captures canonical projection and read-model rules.
type ProjectionImplementationRule string

const (
	ProjectionRuleOwnerKeepsWriteModel          ProjectionImplementationRule = "owner_keeps_canonical_write_model"
	ProjectionRuleFeedByEventOrOutbox           ProjectionImplementationRule = "feed_by_event_outbox_or_projection_contract"
	ProjectionRuleDocumentConsistencyWindow     ProjectionImplementationRule = "consistency_window_must_be_documented"
	ProjectionRulePlanRebuildAndReplayFromStart ProjectionImplementationRule = "rebuild_and_replay_must_be_planned_from_start"
	ProjectionRuleObserveLagAndErrors           ProjectionImplementationRule = "lag_and_errors_must_be_observable"
	ProjectionRuleIdempotentRebuild             ProjectionImplementationRule = "rebuild_must_be_idempotent"
	ProjectionRuleReplayTraceFieldsRequired     ProjectionImplementationRule = "replay_payload_requires_schema_request_and_correlation"
	ProjectionRuleNoCrossModuleOwnerWrite       ProjectionImplementationRule = "no_direct_cross_module_owner_table_write_for_counters"
	ProjectionRuleAlignWithOutbox               ProjectionImplementationRule = "event_producers_align_with_transactional_outbox"
)

var ProjectionImplementationRules = []ProjectionImplementationRule{
	ProjectionRuleOwnerKeepsWriteModel,
	ProjectionRuleFeedByEventOrOutbox,
	ProjectionRuleDocumentConsistencyWindow,
	ProjectionRulePlanRebuildAndReplayFromStart,
	ProjectionRuleObserveLagAndErrors,
	ProjectionRuleIdempotentRebuild,
	ProjectionRuleReplayTraceFieldsRequired,
	ProjectionRuleNoCrossModuleOwnerWrite,
	ProjectionRuleAlignWithOutbox,
}

func IsSupportedConsistencyWindow(value EventualConsistencyWindow) bool {
	switch value {
	case ConsistencyWindowShort, ConsistencyWindowMedium:
		return true
	default:
		return false
	}
}
