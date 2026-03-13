package service

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/admin/dto"
)

func (s *AdminService) ListAuditTrail(ctx context.Context, request dto.ListAuditTrailRequest) (dto.ListAuditTrailResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ListAuditTrailResponse{}, err
	}

	actions, err := s.store.ListAdminActions(ctx, request.ActionType, request.RiskLevel, request.Limit, request.Offset)
	if err != nil {
		return dto.ListAuditTrailResponse{}, err
	}

	items := make([]dto.AuditActionResponse, 0, len(actions))
	for _, action := range actions {
		items = append(items, toAuditActionResponse(action))
	}

	return dto.ListAuditTrailResponse{Items: items, Count: len(items)}, nil
}

func (s *AdminService) GetDashboard(ctx context.Context) (dto.DashboardResponse, error) {
	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.DashboardResponse{}, err
	}

	actions, err := s.store.ListAdminActions(ctx, "", "", 10, 0)
	if err != nil {
		return dto.DashboardResponse{}, err
	}
	overrides, err := s.store.ListOverrideRecords(ctx, "", 0, 0)
	if err != nil {
		return dto.DashboardResponse{}, err
	}
	reviews, err := s.store.ListUserReviewRecords(ctx, "", 0, 0)
	if err != nil {
		return dto.DashboardResponse{}, err
	}
	activeImpersonation, err := s.store.ListImpersonationSessions(ctx, true, 0, 0)
	if err != nil {
		return dto.DashboardResponse{}, err
	}

	latestActions := make([]dto.AuditActionResponse, 0, len(actions))
	for _, action := range actions {
		latestActions = append(latestActions, toAuditActionResponse(action))
	}

	return dto.DashboardResponse{
		MaintenanceEnabled:  cfg.MaintenanceEnabled,
		TotalActions:        len(actions),
		TotalOverrides:      len(overrides),
		TotalUserReviews:    len(reviews),
		ActiveImpersonation: len(activeImpersonation),
		LatestActions:       latestActions,
		GeneratedAt:         s.now().UTC(),
	}, nil
}
