package service

import (
	"fmt"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/notification/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/notification/entity"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/catalog"
	"github.com/google/uuid"
)

func normalizeValue(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

func parseID(raw string, fieldName string) (string, error) {
	parsed, err := uuid.Parse(strings.TrimSpace(raw))
	if err != nil {
		return "", fmt.Errorf("%w: invalid %s", ErrValidation, fieldName)
	}
	return parsed.String(), nil
}

func parseSortBy(raw string, fallback string) string {
	value := normalizeValue(raw)
	switch value {
	case "newest", "oldest":
		return value
	default:
		if fallback == "" {
			return "newest"
		}
		return fallback
	}
}

func parseOptionalCategory(raw string) (catalog.NotificationCategory, error) {
	value := normalizeValue(raw)
	if value == "" {
		return "", nil
	}
	category := catalog.NotificationCategory(value)
	if !catalog.IsValidNotificationCategory(category) {
		return "", fmt.Errorf("%w: invalid category", ErrValidation)
	}
	return category, nil
}

func parseOptionalState(raw string) (entity.NotificationState, error) {
	value := normalizeValue(raw)
	if value == "" {
		return "", nil
	}
	switch entity.NotificationState(value) {
	case entity.NotificationStateCreated,
		entity.NotificationStateDelivered,
		entity.NotificationStateFailed,
		entity.NotificationStateRead:
		return entity.NotificationState(value), nil
	default:
		return "", fmt.Errorf("%w: invalid state", ErrValidation)
	}
}

func parseChannel(raw string, fallback entity.DeliveryChannel) (entity.DeliveryChannel, error) {
	value := normalizeValue(raw)
	if value == "" {
		if fallback == "" {
			return entity.DeliveryChannelInApp, nil
		}
		return fallback, nil
	}
	switch entity.DeliveryChannel(value) {
	case entity.DeliveryChannelInApp,
		entity.DeliveryChannelEmail,
		entity.DeliveryChannelPush:
		return entity.DeliveryChannel(value), nil
	default:
		return "", fmt.Errorf("%w: invalid channel", ErrValidation)
	}
}

func sanitizeText(raw string, fieldName string, maxLen int) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", fmt.Errorf("%w: %s cannot be empty", ErrValidation, fieldName)
	}
	if maxLen > 0 && len(trimmed) > maxLen {
		return "", fmt.Errorf("%w: %s too long", ErrValidation, fieldName)
	}
	trimmed = strings.ReplaceAll(trimmed, "<", "&lt;")
	trimmed = strings.ReplaceAll(trimmed, ">", "&gt;")
	return trimmed, nil
}

func parseMutedCategories(raw []string) ([]catalog.NotificationCategory, error) {
	if len(raw) == 0 {
		return nil, nil
	}

	seen := make(map[catalog.NotificationCategory]struct{}, len(raw))
	result := make([]catalog.NotificationCategory, 0, len(raw))
	for _, value := range raw {
		category, err := parseOptionalCategory(value)
		if err != nil {
			return nil, err
		}
		if category == "" {
			continue
		}
		if _, exists := seen[category]; exists {
			continue
		}
		seen[category] = struct{}{}
		result = append(result, category)
	}

	if len(result) == 0 {
		return nil, nil
	}
	return result, nil
}

func hasMutedCategory(categories []catalog.NotificationCategory, category catalog.NotificationCategory) bool {
	for _, item := range categories {
		if item == category {
			return true
		}
	}
	return false
}

func mapListItem(notification entity.Notification) dto.NotificationListItemResponse {
	return dto.NotificationListItemResponse{
		NotificationID: notification.ID,
		Category:       string(notification.Category),
		Channel:        string(notification.Channel),
		State:          string(notification.State),
		Title:          notification.Title,
		Body:           notification.Body,
		SourceEvent:    notification.SourceEvent,
		SourceRefID:    notification.SourceRefID,
		ReadAt:         notification.ReadAt,
		CreatedAt:      notification.CreatedAt,
		UpdatedAt:      notification.UpdatedAt,
	}
}

func mapDetail(notification entity.Notification) dto.NotificationDetailResponse {
	return dto.NotificationDetailResponse{
		NotificationID:       notification.ID,
		UserID:               notification.UserID,
		Category:             string(notification.Category),
		Channel:              string(notification.Channel),
		TemplateKey:          notification.TemplateKey,
		Title:                notification.Title,
		Body:                 notification.Body,
		State:                string(notification.State),
		SourceEvent:          notification.SourceEvent,
		SourceRefID:          notification.SourceRefID,
		RequestID:            notification.RequestID,
		CorrelationID:        notification.CorrelationID,
		DeliveryAttemptCount: notification.DeliveryAttemptCount,
		LastFailureReason:    notification.LastFailureReason,
		ReadAt:               notification.ReadAt,
		CreatedAt:            notification.CreatedAt,
		UpdatedAt:            notification.UpdatedAt,
	}
}

func mapPreference(preference entity.Preference) dto.PreferenceResponse {
	mutedCategories := make([]string, 0, len(preference.MutedCategories))
	for _, category := range preference.MutedCategories {
		mutedCategories = append(mutedCategories, string(category))
	}

	return dto.PreferenceResponse{
		UserID:            preference.UserID,
		MutedCategories:   mutedCategories,
		QuietHoursEnabled: preference.QuietHoursEnabled,
		QuietHoursStart:   preference.QuietHoursStart,
		QuietHoursEnd:     preference.QuietHoursEnd,
		InAppEnabled:      preference.InAppEnabled,
		EmailEnabled:      preference.EmailEnabled,
		PushEnabled:       preference.PushEnabled,
		DigestEnabled:     preference.DigestEnabled,
		UpdatedAt:         preference.UpdatedAt,
	}
}

func mapRuntimeConfig(cfg entity.RuntimeConfig) dto.RuntimeConfigResponse {
	categoryEnabled := make(map[string]bool, len(cfg.CategoryEnabled))
	for category, enabled := range cfg.CategoryEnabled {
		categoryEnabled[string(category)] = enabled
	}

	channelEnabled := make(map[string]bool, len(cfg.ChannelEnabled))
	for channel, enabled := range cfg.ChannelEnabled {
		channelEnabled[string(channel)] = enabled
	}

	return dto.RuntimeConfigResponse{
		CategoryEnabled: categoryEnabled,
		ChannelEnabled:  channelEnabled,
		DigestEnabled:   cfg.DigestEnabled,
		DeliveryPaused:  cfg.DeliveryPaused,
		UpdatedAt:       cfg.UpdatedAt,
	}
}
