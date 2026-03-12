package service

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/entity"
	"github.com/google/uuid"
)

var slugCleanupPattern = regexp.MustCompile(`[^a-z0-9]+`)

func normalizeValue(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

func normalizeSlug(raw string) string {
	value := normalizeValue(raw)
	value = slugCleanupPattern.ReplaceAllString(value, "-")
	value = strings.Trim(value, "-")
	return value
}

func ensureSlug(slug string, title string, displayNumber string) (string, error) {
	resolved := normalizeSlug(slug)
	if resolved == "" {
		resolved = normalizeSlug(title)
	}
	if resolved == "" {
		resolved = normalizeSlug(displayNumber)
	}
	if resolved == "" {
		return "", fmt.Errorf("%w: invalid slug", ErrValidation)
	}
	return resolved, nil
}

func parseID(raw string, fieldName string) (string, error) {
	parsed, err := uuid.Parse(strings.TrimSpace(raw))
	if err != nil {
		return "", fmt.Errorf("%w: invalid %s", ErrValidation, fieldName)
	}
	return parsed.String(), nil
}

func toPublishState(raw string, defaultState entity.PublishState) (entity.PublishState, error) {
	value := normalizeValue(raw)
	if value == "" {
		return defaultState, nil
	}
	switch entity.PublishState(value) {
	case entity.PublishStateDraft, entity.PublishStateScheduled, entity.PublishStatePublished, entity.PublishStateArchived, entity.PublishStateUnpublished:
		return entity.PublishState(value), nil
	default:
		return "", fmt.Errorf("%w: invalid publish_state", ErrValidation)
	}
}

func toReadAccessLevel(raw string, defaultLevel entity.ReadAccessLevel) (entity.ReadAccessLevel, error) {
	value := normalizeValue(raw)
	if value == "" {
		return defaultLevel, nil
	}
	switch entity.ReadAccessLevel(value) {
	case entity.ReadAccessGuest, entity.ReadAccessAuthenticated, entity.ReadAccessVIP:
		return entity.ReadAccessLevel(value), nil
	default:
		return "", fmt.Errorf("%w: invalid read_access_level", ErrValidation)
	}
}

func toEarlyAccessLevel(raw string, enabled bool, vipOnly bool) (entity.EarlyAccessLevel, error) {
	value := normalizeValue(raw)
	if value == "" {
		if enabled {
			return entity.EarlyAccessVIP, nil
		}
		return entity.EarlyAccessNone, nil
	}
	switch entity.EarlyAccessLevel(value) {
	case entity.EarlyAccessNone, entity.EarlyAccessVIP:
		level := entity.EarlyAccessLevel(value)
		if level == entity.EarlyAccessVIP && !enabled {
			return "", fmt.Errorf("%w: early_access_level requires early_access_enabled=true", ErrValidation)
		}
		if vipOnly && enabled {
			return "", fmt.Errorf("%w: vip_only chapter cannot use early_access_enabled", ErrValidation)
		}
		return level, nil
	default:
		return "", fmt.Errorf("%w: invalid early_access_level", ErrValidation)
	}
}

func toMediaHealthStatus(raw string) (entity.MediaHealthStatus, error) {
	value := normalizeValue(raw)
	switch entity.MediaHealthStatus(value) {
	case entity.MediaHealthHealthy, entity.MediaHealthDegraded, entity.MediaHealthBroken:
		return entity.MediaHealthStatus(value), nil
	default:
		return "", fmt.Errorf("%w: invalid media_health_status", ErrValidation)
	}
}

func toIntegrityStatus(raw string) (entity.IntegrityStatus, error) {
	value := normalizeValue(raw)
	switch entity.IntegrityStatus(value) {
	case entity.IntegrityUnknown, entity.IntegrityPassed, entity.IntegrityFailed:
		return entity.IntegrityStatus(value), nil
	default:
		return "", fmt.Errorf("%w: invalid integrity_status", ErrValidation)
	}
}

func parseRFC3339Ptr(raw *string, fieldName string) (*time.Time, error) {
	if raw == nil {
		return nil, nil
	}
	value := strings.TrimSpace(*raw)
	if value == "" {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid %s", ErrValidation, fieldName)
	}
	resolved := parsed.UTC()
	return &resolved, nil
}

func pagePreviewCount(pageCount int, requested int, enabled bool) int {
	if !enabled {
		return 0
	}
	if requested <= 0 {
		requested = 3
	}
	if requested > pageCount {
		return pageCount
	}
	return requested
}
