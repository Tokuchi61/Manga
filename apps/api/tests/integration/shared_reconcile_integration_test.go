package integration

import (
	"testing"

	"github.com/Tokuchi61/Manga/apps/api/internal/shared/policy"
	"github.com/stretchr/testify/require"
)

func TestCanonicalReconcileFlowsAlignWithTransactionReferences(t *testing.T) {
	flowMap := make(map[string]policy.TransactionFlowReference, len(policy.TransactionFlowReferences))
	for _, flow := range policy.TransactionFlowReferences {
		flowMap[flow.Name] = flow
	}

	for _, reconcileFlow := range policy.CanonicalReconcileFlows {
		txFlow, ok := flowMap[reconcileFlow.Name]
		require.True(t, ok, "transaction reference missing for %s", reconcileFlow.Name)
		require.Equal(t, reconcileFlow.BoundaryType, txFlow.BoundaryType)
		require.NotEmpty(t, txFlow.RecoveryCompensation)
	}
}
