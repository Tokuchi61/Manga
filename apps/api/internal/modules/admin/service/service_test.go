package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/admin/dto"
	adminrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/admin/repository"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestAdminServiceMaintenanceOverrideReviewAndAuditFlow(t *testing.T) {
	store := adminrepository.NewMemoryStore()
	svc := New(store, validation.New())
	now := time.Date(2026, 3, 13, 10, 0, 0, 0, time.UTC)
	svc.now = func() time.Time { return now }

	adminID := uuid.NewString()
	targetUserID := uuid.NewString()

	_, err := svc.UpdateMaintenanceState(context.Background(), adminID, dto.UpdateMaintenanceStateRequest{
		RequestID: "req-maint-1",
		Enabled:   true,
		Reason:    "maintenance",
		RiskLevel: "high",
	})
	require.True(t, errors.Is(err, ErrDoubleConfirmationRequired))

	maintenanceRes, err := svc.UpdateMaintenanceState(context.Background(), adminID, dto.UpdateMaintenanceStateRequest{
		RequestID:         "req-maint-1",
		Enabled:           true,
		Reason:            "maintenance",
		RiskLevel:         "high",
		DoubleConfirmed:   true,
		ConfirmationToken: "confirm-1",
	})
	require.NoError(t, err)
	require.Equal(t, "accepted", maintenanceRes.Status)

	maintenanceIdempotentRes, err := svc.UpdateMaintenanceState(context.Background(), adminID, dto.UpdateMaintenanceStateRequest{
		RequestID:         "req-maint-1",
		Enabled:           true,
		Reason:            "maintenance",
		RiskLevel:         "high",
		DoubleConfirmed:   true,
		ConfirmationToken: "confirm-1",
	})
	require.NoError(t, err)
	require.Equal(t, "idempotent", maintenanceIdempotentRes.Status)

	overrideRes, err := svc.ApplyOverride(context.Background(), adminID, dto.ApplyOverrideRequest{
		RequestID:         "req-override-1",
		TargetModule:      "moderation",
		TargetType:        "case",
		TargetID:          "case-1",
		Decision:          "freeze",
		Reason:            "high-risk",
		RiskLevel:         "critical",
		DoubleConfirmed:   true,
		ConfirmationToken: "confirm-override",
	})
	require.NoError(t, err)
	require.Equal(t, "accepted", overrideRes.Status)

	reviewRes, err := svc.ReviewUser(context.Background(), adminID, dto.ReviewUserRequest{
		RequestID:    "req-review-1",
		TargetUserID: targetUserID,
		Decision:     "warning",
		Reason:       "policy",
		RiskLevel:    "low",
	})
	require.NoError(t, err)
	require.Equal(t, "accepted", reviewRes.Status)

	auditRes, err := svc.ListAuditTrail(context.Background(), dto.ListAuditTrailRequest{})
	require.NoError(t, err)
	require.GreaterOrEqual(t, auditRes.Count, 3)

	dashboardRes, err := svc.GetDashboard(context.Background())
	require.NoError(t, err)
	require.True(t, dashboardRes.MaintenanceEnabled)
	require.Equal(t, 1, dashboardRes.TotalOverrides)
	require.Equal(t, 1, dashboardRes.TotalUserReviews)
}

func TestAdminServiceImpersonationLifecycle(t *testing.T) {
	store := adminrepository.NewMemoryStore()
	svc := New(store, validation.New())
	now := time.Date(2026, 3, 13, 11, 0, 0, 0, time.UTC)
	svc.now = func() time.Time { return now }

	adminID := uuid.NewString()
	targetUserID := uuid.NewString()

	_, err := svc.StartImpersonation(context.Background(), adminID, dto.StartImpersonationRequest{
		RequestID:    "req-start-1",
		TargetUserID: targetUserID,
		Reason:       "investigation",
		RiskLevel:    "high",
	})
	require.True(t, errors.Is(err, ErrDoubleConfirmationRequired))

	startRes, err := svc.StartImpersonation(context.Background(), adminID, dto.StartImpersonationRequest{
		RequestID:         "req-start-1",
		TargetUserID:      targetUserID,
		Reason:            "investigation",
		RiskLevel:         "high",
		DurationMinutes:   20,
		DoubleConfirmed:   true,
		ConfirmationToken: "confirm-start",
	})
	require.NoError(t, err)
	require.Equal(t, "started", startRes.Status)

	_, err = svc.StartImpersonation(context.Background(), adminID, dto.StartImpersonationRequest{
		RequestID:         "req-start-2",
		TargetUserID:      targetUserID,
		Reason:            "duplicate",
		RiskLevel:         "high",
		DurationMinutes:   10,
		DoubleConfirmed:   true,
		ConfirmationToken: "confirm-start-2",
	})
	require.True(t, errors.Is(err, ErrConflict))

	listRes, err := svc.ListImpersonationSessions(context.Background(), dto.ListImpersonationSessionsRequest{ActiveOnly: true})
	require.NoError(t, err)
	require.Equal(t, 1, listRes.Count)

	stopRes, err := svc.StopImpersonation(context.Background(), adminID, dto.StopImpersonationRequest{
		RequestID:         "req-stop-1",
		SessionID:         startRes.SessionID,
		Reason:            "completed",
		RiskLevel:         "high",
		DoubleConfirmed:   true,
		ConfirmationToken: "confirm-stop",
	})
	require.NoError(t, err)
	require.Equal(t, "stopped", stopRes.Status)

	_, err = svc.StopImpersonation(context.Background(), adminID, dto.StopImpersonationRequest{
		RequestID:         "req-stop-2",
		SessionID:         startRes.SessionID,
		Reason:            "already-stopped",
		RiskLevel:         "high",
		DoubleConfirmed:   true,
		ConfirmationToken: "confirm-stop-2",
	})
	require.True(t, errors.Is(err, ErrConflict))
}
