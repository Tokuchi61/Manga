package policy

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOperationalPolicyCanonicalValues(t *testing.T) {
	require.Equal(t, []RateLimitSurface{
		"auth.login",
		"comment.write",
		"support.intake",
		"social.messaging",
		"payment.callback",
		"ads.click_intake",
	}, AllRateLimitSurfaces)
	require.Equal(t, "request_id", TraceFieldRequestID)
	require.Equal(t, "correlation_id", TraceFieldCorrelationID)
	require.Equal(t, "provider_event_reference", CallbackAdditionalTraceReference)
}

func TestOutboxPolicyCanonicalValues(t *testing.T) {
	require.Equal(t, []OutboxComponent{
		"transactional_record",
		"background_publisher",
		"retry_backoff",
		"dead_letter",
		"observability_metrics",
	}, RequiredOutboxComponents)

	require.Equal(t, []OutboxMessageField{
		"event_id",
		"schema_version",
		"request_id",
		"correlation_id",
		"causation_id",
	}, RequiredOutboxMessageFields)

	require.Equal(t, []string{"payment", "inventory", "mission", "royalpass", "notification", "support", "moderation", "history"}, OutboxPriorityModules)
}

func TestProjectionCatalogHasExpectedEntries(t *testing.T) {
	require.Len(t, CanonicalProjections, 9)
	first := CanonicalProjections[0]
	require.Equal(t, "manga.comment_count", first.Name)
	require.Equal(t, "comment", first.CanonicalWriteModel)
	require.True(t, first.ReplaySupported)
	last := CanonicalProjections[len(CanonicalProjections)-1]
	require.Equal(t, "royalpass.tier_progress_snapshot", last.Name)
}

func TestTechnicalAndStrategyDecisionsAreStable(t *testing.T) {
	require.Equal(t, []TechnicalChoice{
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
	}, CanonicalTechnicalStack)

	require.Equal(t, "redis", CanonicalCacheBackend)
	require.Equal(t, "postgresql_jobs_and_transactional_outbox", CanonicalAsyncBaseline)
	require.Equal(t, "postgresql_full_text_search", CanonicalSearchEngine)
	require.Equal(t, "shared_platform_media_infrastructure", CanonicalMediaBaseline)
	require.Equal(t, "projection_and_read_model_export_query_layer", CanonicalReportingBaseline)

	require.Len(t, BrokerTransitionCriteria, 4)
	require.Len(t, SearchProviderSwapCriteria, 3)
	require.Len(t, ReportingServiceSplitCriteria, 3)
	require.Len(t, MediaModuleSplitCriteria, 3)
}

func TestTransactionRulesAndReferences(t *testing.T) {
	require.Len(t, TransactionSelectionRules, 6)
	require.Len(t, TransactionFlowReferences, 6)
	require.Equal(t, "auth_login_and_session_create", TransactionFlowReferences[0].Name)
	require.Equal(t, "royalpass_premium_activation", TransactionFlowReferences[5].Name)
}
