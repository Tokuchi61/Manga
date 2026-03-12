package policy

// ReportingLayer defines canonical reporting/read model layers.
type ReportingLayer string

const (
	ReportingLayerOperationalSummary ReportingLayer = "operational_summary"
	ReportingLayerAnalyticsAggregate ReportingLayer = "analytics_aggregate"
	ReportingLayerExportQuery        ReportingLayer = "export_query_layer"
)

var ReportingLayers = []ReportingLayer{
	ReportingLayerOperationalSummary,
	ReportingLayerAnalyticsAggregate,
	ReportingLayerExportQuery,
}

var ReportingRules = []string{
	"reporting_projections_must_not_replace_canonical_write_models",
	"export_surfaces_require_authorized_operation_flow_and_audit",
	"payment_ads_support_reporting_surfaces_follow_same_read_model_principles",
	"dashboard_metrics_must_not_be_treated_as_business_source_of_truth",
}
