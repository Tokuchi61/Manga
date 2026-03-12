package service

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/support/entity"
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

func parsePriority(raw string, fallback entity.SupportPriority) (entity.SupportPriority, error) {
	value := normalizeValue(raw)
	if value == "" {
		if fallback == "" {
			return entity.SupportPriorityNormal, nil
		}
		return fallback, nil
	}

	switch entity.SupportPriority(value) {
	case entity.SupportPriorityLow,
		entity.SupportPriorityNormal,
		entity.SupportPriorityHigh,
		entity.SupportPriorityUrgent:
		return entity.SupportPriority(value), nil
	default:
		return "", fmt.Errorf("%w: invalid priority", ErrValidation)
	}
}

func parseStatus(raw string) (entity.SupportStatus, error) {
	value := normalizeValue(raw)
	switch entity.SupportStatus(value) {
	case entity.SupportStatusOpen,
		entity.SupportStatusTriaged,
		entity.SupportStatusWaitingUser,
		entity.SupportStatusWaitingTeam,
		entity.SupportStatusResolved,
		entity.SupportStatusRejected,
		entity.SupportStatusClosed,
		entity.SupportStatusSpam:
		return entity.SupportStatus(value), nil
	default:
		return "", fmt.Errorf("%w: invalid status", ErrValidation)
	}
}

func parseTargetType(raw string) (entity.SupportTargetType, error) {
	value := normalizeValue(raw)
	switch entity.SupportTargetType(value) {
	case entity.SupportTargetTypeManga,
		entity.SupportTargetTypeChapter,
		entity.SupportTargetTypeComment:
		return entity.SupportTargetType(value), nil
	default:
		return "", fmt.Errorf("%w: invalid target_type", ErrValidation)
	}
}

func parseReplyVisibility(raw string) (entity.ReplyVisibility, error) {
	value := normalizeValue(raw)
	switch entity.ReplyVisibility(value) {
	case entity.ReplyVisibilityPublicToRequester, entity.ReplyVisibilityInternalOnly:
		return entity.ReplyVisibility(value), nil
	default:
		return "", fmt.Errorf("%w: invalid visibility", ErrValidation)
	}
}

func parseSortBy(raw string, fallback string) string {
	value := normalizeValue(raw)
	switch value {
	case "newest", "oldest", "priority":
		return value
	default:
		if fallback == "" {
			return "newest"
		}
		return fallback
	}
}

func sanitizeText(raw string, fieldName string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", fmt.Errorf("%w: %s cannot be empty", ErrValidation, fieldName)
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

func spamRiskScore(reasonText string, attachments []string) int {
	risk := 0
	lower := strings.ToLower(reasonText)

	if len(reasonText) > 1500 {
		risk += 10
	}
	if strings.Count(lower, "http://")+strings.Count(lower, "https://") > 2 {
		risk += 20
	}
	if len(attachments) > 3 {
		risk += 15
	}
	if hasExcessiveRepetition(lower) {
		risk += 25
	}
	if strings.Contains(lower, "casino") || strings.Contains(lower, "telegram") {
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
	var previous rune
	for index, current := range content {
		if index == 0 {
			previous = current
			continue
		}

		if current == previous {
			repeated++
			if repeated >= 8 {
				return true
			}
		} else {
			repeated = 1
			previous = current
		}
	}
	return false
}

func toStringPtr(value *entity.SupportTargetType) *string {
	if value == nil {
		return nil
	}
	v := string(*value)
	return &v
}
