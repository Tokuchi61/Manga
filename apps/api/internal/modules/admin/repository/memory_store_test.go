package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/admin/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestMemoryStoreRuntimeActionAndSnapshot(t *testing.T) {
	store := NewMemoryStore()
	now := time.Now().UTC()
	store.now = func() time.Time { return now }

	err := store.UpdateRuntimeConfig(context.Background(), entity.RuntimeConfig{MaintenanceEnabled: true, UpdatedAt: now})
	require.NoError(t, err)

	action := entity.AdminAction{
		ActionID:                   uuid.NewString(),
		ActionType:                 entity.ActionTypeSettingChanged,
		ActorUserID:                uuid.NewString(),
		TargetModule:               "site",
		TargetType:                 "setting",
		TargetID:                   "site.maintenance.enabled",
		RequestID:                  "req-admin-1",
		Reason:                     "maintenance",
		Result:                     entity.ActionResultSuccess,
		RiskLevel:                  entity.RiskLevelHigh,
		RequiresDoubleConfirmation: true,
		DoubleConfirmed:            true,
		ConfirmationToken:          "token-1",
		CreatedAt:                  now,
		UpdatedAt:                  now,
	}
	require.NoError(t, store.CreateAdminAction(context.Background(), action))
	require.NoError(t, store.PutActionDedup(context.Background(), "setting:req-admin-1", action))

	override := entity.OverrideRecord{
		OverrideID:   uuid.NewString(),
		ActionID:     action.ActionID,
		TargetModule: "moderation",
		TargetType:   "case",
		TargetID:     "case-1",
		Decision:     entity.OverrideDecisionFreeze,
		Reason:       "risk",
		RiskLevel:    entity.RiskLevelHigh,
		Active:       true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	require.NoError(t, store.CreateOverrideRecord(context.Background(), override))

	review := entity.UserReviewRecord{
		ReviewID:     uuid.NewString(),
		ActionID:     action.ActionID,
		TargetUserID: uuid.NewString(),
		Decision:     entity.UserReviewDecisionWarning,
		Reason:       "policy",
		RiskLevel:    entity.RiskLevelLow,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	require.NoError(t, store.CreateUserReviewRecord(context.Background(), review))

	session := entity.ImpersonationSession{
		SessionID:    uuid.NewString(),
		ActionID:     action.ActionID,
		ActorUserID:  action.ActorUserID,
		TargetUserID: review.TargetUserID,
		Reason:       "investigation",
		RiskLevel:    entity.RiskLevelCritical,
		Active:       true,
		StartedAt:    now,
		ExpiresAt:    now.Add(15 * time.Minute),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	require.NoError(t, store.CreateImpersonationSession(context.Background(), session))

	payload, err := store.Snapshot()
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	restored := NewMemoryStore()
	restored.now = func() time.Time { return now }
	require.NoError(t, restored.RestoreSnapshot(payload))

	cfg, err := restored.GetRuntimeConfig(context.Background())
	require.NoError(t, err)
	require.True(t, cfg.MaintenanceEnabled)

	actions, err := restored.ListAdminActions(context.Background(), "", "", 0, 0)
	require.NoError(t, err)
	require.Len(t, actions, 1)

	activeSession, err := restored.FindActiveImpersonation(context.Background(), action.ActorUserID, review.TargetUserID)
	require.NoError(t, err)
	require.Equal(t, session.SessionID, activeSession.SessionID)
}
