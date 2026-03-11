package policy

// TechnicalChoice stores a canonical stack decision.
type TechnicalChoice struct {
	Responsibility string
	Canonical      string
}

var CanonicalTechnicalStack = []TechnicalChoice{
	{Responsibility: "http_router", Canonical: "chi"},
	{Responsibility: "config_loader", Canonical: "caarlos0/env"},
	{Responsibility: "structured_logging", Canonical: "zap"},
	{Responsibility: "uuid", Canonical: "google/uuid"},
	{Responsibility: "input_validation", Canonical: "go-playground/validator/v10"},
	{Responsibility: "migration", Canonical: "golang-migrate"},
	{Responsibility: "sql_access", Canonical: "pgx/v5"},
	{Responsibility: "connection_pool", Canonical: "pgxpool"},
	{Responsibility: "password_hashing", Canonical: "argon2id"},
	{Responsibility: "test_helpers", Canonical: "testify"},
}

const (
	CanonicalCacheBackend      = "redis"
	CanonicalAsyncBaseline     = "postgresql_jobs_and_transactional_outbox"
	CanonicalSearchEngine      = "postgresql_full_text_search"
	CanonicalMediaBaseline     = "shared_platform_media_infrastructure"
	CanonicalReportingBaseline = "projection_and_read_model_export_query_layer"
)

var BrokerTransitionCriteria = []string{
	"high_volume_fan_out_across_modules",
	"persistent_lag_for_independent_consumer_groups",
	"db_backed_jobs_operational_observability_or_retry_cost_too_high",
	"heavy_bidirectional_event_integration_with_external_systems_or_services",
}

var SearchProviderSwapCriteria = []string{
	"postgresql_fts_cannot_meet_quality_or_latency",
	"need_typo_tolerance_weighted_ranking_synonyms_or_facets_beyond_baseline",
	"separate_high_volume_index_ops_make_db_load_unsustainable",
}

var ReportingServiceSplitCriteria = []string{
	"summary_and_analytics_queries_cannot_share_same_db_load",
	"high_volume_history_and_event_flows_need_separate_storage",
	"separate_team_access_policy_or_data_retention_requirements",
}

var MediaModuleSplitCriteria = []string{
	"multiple_modules_share_the_same_upload_transform_delivery_pipeline",
	"media_lifecycle_requires_operations_independent_from_business_modules",
	"storage_scanning_variant_generation_and_signed_access_need_separate_service_boundary",
}
