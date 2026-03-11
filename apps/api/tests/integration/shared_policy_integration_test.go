package integration

import (
	"testing"

	"github.com/Tokuchi61/Manga/apps/api/internal/shared/policy"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/settings"
	"github.com/stretchr/testify/require"
)

func TestRateLimitSurfacesCompatibleWithSettingsMetricGrammar(t *testing.T) {
	for _, surface := range policy.AllRateLimitSurfaces {
		require.True(t, settings.IsModuleSurfaceMetricKey(string(surface)))
	}
}

func TestOutboxTraceFieldsAlignWithOperationalTraceModel(t *testing.T) {
	require.Contains(t, policy.RequiredOutboxMessageFields, policy.OutboxMessageField(policy.TraceFieldRequestID))
	require.Contains(t, policy.RequiredOutboxMessageFields, policy.OutboxMessageField(policy.TraceFieldCorrelationID))
}
