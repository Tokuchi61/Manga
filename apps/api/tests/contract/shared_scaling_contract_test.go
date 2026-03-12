package contract

import (
	"testing"

	"github.com/Tokuchi61/Manga/apps/api/internal/shared/policy"
	"github.com/stretchr/testify/require"
)

func TestStageThreeScalingContractValues(t *testing.T) {
	require.Contains(t, policy.AllModuleStatuses, policy.ModuleStatusActive)
	require.Contains(t, policy.SuggestedDomainGroups, "operations")
	require.Contains(t, policy.ProjectionImplementationRules, policy.ProjectionRuleIdempotentRebuild)
	require.Contains(t, policy.ReportingLayers, policy.ReportingLayerExportQuery)
	require.Contains(t, policy.ReconcileRules, policy.ReconcileRulePreserveTraceAndIdempotency)
	require.Contains(t, policy.MaintenanceRefactorRules, policy.MaintenanceRuleNoUnnecessaryRefactor)
}
