package service

import (
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/admin/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/admin/entity"
)

func normalizeValue(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

func buildActionDedupKey(scope string, requestID string) string {
	return normalizeValue(scope) + ":" + normalizeValue(requestID)
}

func requiresDoubleConfirmation(riskLevel string) bool {
	switch normalizeValue(riskLevel) {
	case entity.RiskLevelHigh, entity.RiskLevelCritical:
		return true
	default:
		return false
	}
}

func requireDoubleConfirmation(riskLevel string, confirmed bool, token string) error {
	if !requiresDoubleConfirmation(riskLevel) {
		return nil
	}
	if !confirmed || strings.TrimSpace(token) == "" {
		return ErrDoubleConfirmationRequired
	}
	return nil
}

func toAuditActionResponse(action entity.AdminAction) dto.AuditActionResponse {
	return dto.AuditActionResponse{
		ActionID:                   action.ActionID,
		ActionType:                 action.ActionType,
		ActorUserID:                action.ActorUserID,
		TargetUserID:               action.TargetUserID,
		TargetModule:               action.TargetModule,
		TargetType:                 action.TargetType,
		TargetID:                   action.TargetID,
		RequestID:                  action.RequestID,
		CorrelationID:              action.CorrelationID,
		Reason:                     action.Reason,
		Result:                     action.Result,
		RiskLevel:                  action.RiskLevel,
		RequiresDoubleConfirmation: action.RequiresDoubleConfirmation,
		DoubleConfirmed:            action.DoubleConfirmed,
		Metadata:                   action.Metadata,
		CreatedAt:                  action.CreatedAt,
	}
}

func toOverrideResponse(status string, eventName string, actionID string, record entity.OverrideRecord) dto.OverrideResponse {
	return dto.OverrideResponse{
		Status:       status,
		ActionID:     actionID,
		OverrideID:   record.OverrideID,
		TargetModule: record.TargetModule,
		TargetType:   record.TargetType,
		TargetID:     record.TargetID,
		Decision:     record.Decision,
		RiskLevel:    record.RiskLevel,
		Active:       record.Active,
		Event:        eventName,
		ExpiresAt:    record.ExpiresAt,
		UpdatedAt:    record.UpdatedAt,
	}
}

func toUserReviewResponse(status string, eventName string, actionID string, record entity.UserReviewRecord) dto.UserReviewResponse {
	return dto.UserReviewResponse{
		Status:       status,
		ActionID:     actionID,
		ReviewID:     record.ReviewID,
		TargetUserID: record.TargetUserID,
		Decision:     record.Decision,
		RiskLevel:    record.RiskLevel,
		Event:        eventName,
		UpdatedAt:    record.UpdatedAt,
	}
}

func toImpersonationResponse(status string, actionID string, session entity.ImpersonationSession) dto.ImpersonationSessionResponse {
	return dto.ImpersonationSessionResponse{
		Status:       status,
		ActionID:     actionID,
		SessionID:    session.SessionID,
		ActorUserID:  session.ActorUserID,
		TargetUserID: session.TargetUserID,
		RiskLevel:    session.RiskLevel,
		Reason:       session.Reason,
		Active:       session.Active,
		StartedAt:    session.StartedAt,
		EndsAt:       session.EndedAt,
		ExpiresAt:    session.ExpiresAt,
		UpdatedAt:    session.UpdatedAt,
	}
}
