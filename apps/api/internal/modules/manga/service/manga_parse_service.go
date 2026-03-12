package service

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/entity"
	"github.com/google/uuid"
)

var slugCleanupPattern = regexp.MustCompile(`[^a-z0-9]+`)

func parseID(raw string, fieldName string) (string, error) {
	parsed, err := uuid.Parse(strings.TrimSpace(raw))
	if err != nil {
		return "", fmt.Errorf("%w: invalid %s", ErrValidation, fieldName)
	}
	return parsed.String(), nil
}

func normalizeSlug(raw string) string {
	value := strings.ToLower(strings.TrimSpace(raw))
	value = slugCleanupPattern.ReplaceAllString(value, "-")
	value = strings.Trim(value, "-")
	return value
}

func ensureSlug(slug string, title string) (string, error) {
	resolved := normalizeSlug(slug)
	if resolved == "" {
		resolved = normalizeSlug(title)
	}
	if resolved == "" {
		return "", fmt.Errorf("%w: invalid slug", ErrValidation)
	}
	return resolved, nil
}

func normalizeList(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	result := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		normalized := strings.ToLower(strings.TrimSpace(value))
		if normalized == "" {
			continue
		}
		if _, exists := seen[normalized]; exists {
			continue
		}
		seen[normalized] = struct{}{}
		result = append(result, normalized)
	}
	return result
}

func normalizeFreeText(raw string) string {
	return strings.TrimSpace(raw)
}

func toPublishState(raw string, defaultState entity.PublishState) (entity.PublishState, error) {
	value := strings.ToLower(strings.TrimSpace(raw))
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

func toVisibility(raw string) (entity.VisibilityState, error) {
	value := strings.ToLower(strings.TrimSpace(raw))
	switch entity.VisibilityState(value) {
	case entity.VisibilityPublic, entity.VisibilityHidden:
		return entity.VisibilityState(value), nil
	default:
		return "", fmt.Errorf("%w: invalid visibility", ErrValidation)
	}
}

func toReadAccessLevel(raw string) (entity.ReadAccessLevel, error) {
	value := strings.ToLower(strings.TrimSpace(raw))
	if value == "" {
		return entity.ReadAccessAuthenticated, nil
	}
	switch entity.ReadAccessLevel(value) {
	case entity.ReadAccessGuest, entity.ReadAccessAuthenticated:
		return entity.ReadAccessLevel(value), nil
	default:
		return "", fmt.Errorf("%w: invalid default_read_access_level", ErrValidation)
	}
}

func toEarlyAccessLevel(raw string, enabled bool) (entity.EarlyAccessLevel, error) {
	value := strings.ToLower(strings.TrimSpace(raw))
	if value == "" {
		if enabled {
			return entity.EarlyAccessVIP, nil
		}
		return entity.EarlyAccessNone, nil
	}
	switch entity.EarlyAccessLevel(value) {
	case entity.EarlyAccessNone, entity.EarlyAccessVIP:
		if !enabled && entity.EarlyAccessLevel(value) != entity.EarlyAccessNone {
			return "", fmt.Errorf("%w: early access level requires enabled=true", ErrValidation)
		}
		return entity.EarlyAccessLevel(value), nil
	default:
		return "", fmt.Errorf("%w: invalid default_early_access_level", ErrValidation)
	}
}

func parseRFC3339Ptr(raw *string, fieldName string) (*time.Time, error) {
	if raw == nil {
		return nil, nil
	}
	trimmed := strings.TrimSpace(*raw)
	if trimmed == "" {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, trimmed)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid %s", ErrValidation, fieldName)
	}
	value := parsed.UTC()
	return &value, nil
}

func boolValue(ptr *bool, fallback bool) bool {
	if ptr == nil {
		return fallback
	}
	return *ptr
}

func normalizeValue(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}
