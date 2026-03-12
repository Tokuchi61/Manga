package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/mission/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/mission/entity"
)

func normalizeValue(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

func adjustedNow(now time.Time, resetHour int) time.Time {
	if resetHour < 0 || resetHour > 23 {
		resetHour = 0
	}
	return now.UTC().Add(-time.Duration(resetHour) * time.Hour)
}

func buildPeriodKey(category string, now time.Time, resetHour int) string {
	adjusted := adjustedNow(now, resetHour)
	switch normalizeValue(category) {
	case "weekly":
		year, week := adjusted.ISOWeek()
		return fmt.Sprintf("%04d-W%02d", year, week)
	case "monthly":
		return adjusted.Format("2006-01")
	case "event":
		return "event"
	case "level":
		return "lifetime"
	default:
		return adjusted.Format("2006-01-02")
	}
}

func isMissionActive(definition entity.MissionDefinition, now time.Time) bool {
	if !definition.Active {
		return false
	}
	utcNow := now.UTC()
	if definition.StartsAt != nil && utcNow.Before(definition.StartsAt.UTC()) {
		return false
	}
	if definition.EndsAt != nil && utcNow.After(definition.EndsAt.UTC()) {
		return false
	}
	return true
}

func resolveMissionStatus(definition entity.MissionDefinition, progress *entity.MissionProgress, now time.Time) string {
	if !isMissionActive(definition, now) {
		return "expired"
	}
	if progress == nil {
		return "active"
	}
	if progress.Claimed {
		return "claimed"
	}
	if progress.Completed {
		return "completed"
	}
	return "active"
}

func toProgressItemResponse(definition entity.MissionDefinition, progress entity.MissionProgress, now time.Time, resetHour int) dto.MissionProgressItemResponse {
	periodKey := progress.PeriodKey
	if strings.TrimSpace(periodKey) == "" {
		periodKey = buildPeriodKey(definition.Category, now, resetHour)
	}

	status := resolveMissionStatus(definition, &progress, now)
	if strings.TrimSpace(progress.UserID) == "" {
		status = resolveMissionStatus(definition, nil, now)
	}

	return dto.MissionProgressItemResponse{
		MissionID:      definition.MissionID,
		Category:       definition.Category,
		Title:          definition.Title,
		ObjectiveType:  definition.ObjectiveType,
		TargetCount:    definition.TargetCount,
		ProgressCount:  progress.ProgressCount,
		RewardItemID:   definition.RewardItemID,
		RewardQuantity: definition.RewardQuantity,
		Status:         status,
		PeriodKey:      periodKey,
		CompletedAt:    progress.CompletedAt,
		ClaimedAt:      progress.ClaimedAt,
		UpdatedAt:      progress.UpdatedAt,
	}
}

func toMissionDetailResponse(definition entity.MissionDefinition, progress entity.MissionProgress, now time.Time, resetHour int) dto.MissionDetailResponse {
	item := toProgressItemResponse(definition, progress, now, resetHour)
	return dto.MissionDetailResponse{
		MissionID:      item.MissionID,
		Category:       item.Category,
		Title:          item.Title,
		ObjectiveType:  item.ObjectiveType,
		TargetCount:    item.TargetCount,
		ProgressCount:  item.ProgressCount,
		RewardItemID:   item.RewardItemID,
		RewardQuantity: item.RewardQuantity,
		Status:         item.Status,
		PeriodKey:      item.PeriodKey,
		StartsAt:       definition.StartsAt,
		EndsAt:         definition.EndsAt,
		CompletedAt:    item.CompletedAt,
		ClaimedAt:      item.ClaimedAt,
		CreatedAt:      definition.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
	}
}

func toDefinitionResponse(definition entity.MissionDefinition) dto.MissionDefinitionResponse {
	return dto.MissionDefinitionResponse{
		MissionID:      definition.MissionID,
		Category:       definition.Category,
		Title:          definition.Title,
		ObjectiveType:  definition.ObjectiveType,
		TargetCount:    definition.TargetCount,
		RewardItemID:   definition.RewardItemID,
		RewardQuantity: definition.RewardQuantity,
		Active:         definition.Active,
		StartsAt:       definition.StartsAt,
		EndsAt:         definition.EndsAt,
		UpdatedAt:      definition.UpdatedAt,
	}
}

func buildProgressDedupKey(userID string, missionID string, periodKey string, requestID string) string {
	return normalizeValue(userID) + ":" + normalizeValue(missionID) + ":" + normalizeValue(periodKey) + ":" + normalizeValue(requestID)
}

func buildClaimDedupKey(userID string, missionID string, periodKey string, requestID string) string {
	return normalizeValue(userID) + ":" + normalizeValue(missionID) + ":" + normalizeValue(periodKey) + ":" + normalizeValue(requestID)
}

func applyOffsetLimit[T any](items []T, offset int, limit int) []T {
	if offset < 0 {
		offset = 0
	}
	if offset >= len(items) {
		return []T{}
	}
	if limit <= 0 {
		limit = 50
	}
	end := offset + limit
	if end > len(items) {
		end = len(items)
	}
	return append([]T(nil), items[offset:end]...)
}
