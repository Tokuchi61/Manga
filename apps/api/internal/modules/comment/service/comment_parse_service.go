package service

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/comment/entity"
	"github.com/google/uuid"
)

var whitespacePattern = regexp.MustCompile(`\s+`)

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

func parseTargetType(raw string) (entity.TargetType, error) {
	value := normalizeValue(raw)
	switch entity.TargetType(value) {
	case entity.TargetTypeManga, entity.TargetTypeChapter:
		return entity.TargetType(value), nil
	default:
		return "", fmt.Errorf("%w: invalid target_type", ErrValidation)
	}
}

func parseSortBy(raw string, defaultSort string) string {
	value := normalizeValue(raw)
	switch value {
	case "newest", "oldest", "popular":
		return value
	default:
		if defaultSort == "" {
			return "newest"
		}
		return defaultSort
	}
}

func parseModerationStatus(raw string) (entity.ModerationStatus, error) {
	value := normalizeValue(raw)
	switch entity.ModerationStatus(value) {
	case entity.ModerationStatusVisible, entity.ModerationStatusHidden, entity.ModerationStatusFlagged:
		return entity.ModerationStatus(value), nil
	default:
		return "", fmt.Errorf("%w: invalid moderation_status", ErrValidation)
	}
}

func sanitizeContent(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", fmt.Errorf("%w: content cannot be empty", ErrValidation)
	}
	collapsed := whitespacePattern.ReplaceAllString(trimmed, " ")
	collapsed = strings.ReplaceAll(collapsed, "<", "&lt;")
	collapsed = strings.ReplaceAll(collapsed, ">", "&gt;")
	return collapsed, nil
}

func normalizeAttachments(items []string) []string {
	if len(items) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(items))
	result := make([]string, 0, len(items))
	for _, item := range items {
		value := strings.TrimSpace(item)
		if value == "" {
			continue
		}
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func spamRiskScore(content string, attachments []string) int {
	risk := 0
	lower := strings.ToLower(content)
	if len(content) > 1000 {
		risk += 10
	}
	if strings.Count(lower, "http://")+strings.Count(lower, "https://") > 2 {
		risk += 15
	}
	if len(attachments) > 3 {
		risk += 10
	}
	if hasExcessiveRepetition(lower) {
		risk += 20
	}
	if risk > 100 {
		risk = 100
	}
	if risk < 0 {
		risk = 0
	}
	return risk
}

func hasExcessiveRepetition(content string) bool {
	if len(content) < 12 {
		return false
	}
	repeated := 1
	var prev rune
	for idx, r := range content {
		if idx == 0 {
			prev = r
			continue
		}
		if r == prev {
			repeated++
			if repeated >= 8 {
				return true
			}
		} else {
			repeated = 1
			prev = r
		}
	}
	return false
}

func toSafeContent(comment entity.Comment) string {
	if comment.DeletedAt != nil {
		return "[deleted]"
	}
	return comment.SanitizedContent
}
