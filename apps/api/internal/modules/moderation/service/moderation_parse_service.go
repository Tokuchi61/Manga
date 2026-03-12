package service

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/entity"
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

func parseOptionalCaseStatus(raw string) (entity.CaseStatus, error) {
	value := normalizeValue(raw)
	if value == "" {
		return "", nil
	}
	switch entity.CaseStatus(value) {
	case entity.CaseStatusNew,
		entity.CaseStatusQueued,
		entity.CaseStatusAssigned,
		entity.CaseStatusInReview,
		entity.CaseStatusEscalated,
		entity.CaseStatusResolved,
		entity.CaseStatusRejected,
		entity.CaseStatusClosed:
		return entity.CaseStatus(value), nil
	default:
		return "", fmt.Errorf("%w: invalid status", ErrValidation)
	}
}

func parseTargetType(raw string) (entity.TargetType, error) {
	value := normalizeValue(raw)
	switch entity.TargetType(value) {
	case entity.TargetTypeManga,
		entity.TargetTypeChapter,
		entity.TargetTypeComment:
		return entity.TargetType(value), nil
	default:
		return "", fmt.Errorf("%w: invalid target_type", ErrValidation)
	}
}

func parseOptionalTargetType(raw string) (entity.TargetType, error) {
	value := normalizeValue(raw)
	if value == "" {
		return "", nil
	}
	return parseTargetType(value)
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

func parseActionType(raw string) (entity.ActionType, error) {
	value := normalizeValue(raw)
	switch entity.ActionType(value) {
	case entity.ActionTypeHide,
		entity.ActionTypeUnhide,
		entity.ActionTypeLock,
		entity.ActionTypeUnlock,
		entity.ActionTypeWarning,
		entity.ActionTypeReviewComplete,
		entity.ActionTypeEscalate:
		return entity.ActionType(value), nil
	default:
		return "", fmt.Errorf("%w: invalid action_type", ErrValidation)
	}
}

func mapActionResult(actionType entity.ActionType) entity.ActionResult {
	switch actionType {
	case entity.ActionTypeHide:
		return entity.ActionResultContentHidden
	case entity.ActionTypeUnhide:
		return entity.ActionResultContentRestored
	case entity.ActionTypeWarning:
		return entity.ActionResultWarningSent
	default:
		return entity.ActionResultNoAction
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
