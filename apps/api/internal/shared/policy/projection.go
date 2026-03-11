package policy

// ProjectionDefinition stores canonical projection references.
type ProjectionDefinition struct {
	Name                string
	CanonicalWriteModel string
	EventSources        []string
	ConsumerSurface     string
	ConsistencyWindow   string
	RebuildPath         string
	ReplaySupported     bool
}

var CanonicalProjections = []ProjectionDefinition{
	{
		Name:                "manga.comment_count",
		CanonicalWriteModel: "comment",
		EventSources:        []string{"comment.created", "comment.deleted", "comment.moderated"},
		ConsumerSurface:     "manga_detail_and_listing",
		ConsistencyWindow:   "short",
		RebuildPath:         "recount_from_comment_table_then_incremental_catch_up",
		ReplaySupported:     true,
	},
	{
		Name:                "manga.engagement_summary",
		CanonicalWriteModel: "history_comment",
		EventSources:        []string{"history.checkpoint", "engagement.events"},
		ConsumerSurface:     "discovery_and_admin_summary",
		ConsistencyWindow:   "medium",
		RebuildPath:         "daily_batch_rebuild_then_targeted_repair",
		ReplaySupported:     true,
	},
	{
		Name:                "history.continue_reading_projection",
		CanonicalWriteModel: "history",
		EventSources:        []string{"history.checkpoint", "history.finish"},
		ConsumerSurface:     "continue_reading",
		ConsistencyWindow:   "short",
		RebuildPath:         "recompute_from_last_checkpoint",
		ReplaySupported:     true,
	},
	{
		Name:                "notification.unread_counter",
		CanonicalWriteModel: "notification",
		EventSources:        []string{"notification.created", "notification.read"},
		ConsumerSurface:     "inbox_badge_and_header_counters",
		ConsistencyWindow:   "short",
		RebuildPath:         "recount_from_unread_records",
		ReplaySupported:     true,
	},
	{
		Name:                "support.queue_summary",
		CanonicalWriteModel: "support",
		EventSources:        []string{"support.create", "support.status_change", "support.assignee_change"},
		ConsumerSurface:     "support_operations_panel",
		ConsistencyWindow:   "short",
		RebuildPath:         "status_group_regroup",
		ReplaySupported:     true,
	},
	{
		Name:                "moderation.queue_summary",
		CanonicalWriteModel: "moderation",
		EventSources:        []string{"moderation.case_create", "moderation.assignment", "moderation.resolution"},
		ConsumerSurface:     "moderation_panel",
		ConsistencyWindow:   "short",
		RebuildPath:         "case_status_regroup",
		ReplaySupported:     true,
	},
	{
		Name:                "ads.impression_aggregate",
		CanonicalWriteModel: "ads",
		EventSources:        []string{"ads.impression.accepted", "ads.click.accepted"},
		ConsumerSurface:     "reporting_and_dashboard",
		ConsistencyWindow:   "medium",
		RebuildPath:         "batch_aggregation_job",
		ReplaySupported:     true,
	},
	{
		Name:                "mission.progress_projection",
		CanonicalWriteModel: "mission",
		EventSources:        []string{"mission.progress", "mission.claim", "mission.reset"},
		ConsumerSurface:     "mission_list_and_progress_summary",
		ConsistencyWindow:   "short",
		RebuildPath:         "objective_based_recompute",
		ReplaySupported:     true,
	},
	{
		Name:                "royalpass.tier_progress_snapshot",
		CanonicalWriteModel: "royalpass",
		EventSources:        []string{"royalpass.progress", "royalpass.claim"},
		ConsumerSurface:     "season_overview_and_tier_views",
		ConsistencyWindow:   "short",
		RebuildPath:         "tier_based_recompute",
		ReplaySupported:     true,
	},
}
