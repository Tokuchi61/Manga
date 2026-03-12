package service

import (
	"fmt"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/user/entity"
	"github.com/google/uuid"
)

func parseID(raw string, fieldName string) (string, error) {
	parsed, err := uuid.Parse(strings.TrimSpace(raw))
	if err != nil {
		return "", fmt.Errorf("%w: invalid %s", ErrValidation, fieldName)
	}
	return parsed.String(), nil
}

func toProfileVisibility(raw string) (entity.ProfileVisibility, error) {
	value := strings.TrimSpace(strings.ToLower(raw))
	switch entity.ProfileVisibility(value) {
	case entity.ProfileVisibilityPublic, entity.ProfileVisibilityPrivate:
		return entity.ProfileVisibility(value), nil
	default:
		return "", fmt.Errorf("%w: invalid profile_visibility", ErrValidation)
	}
}

func toHistoryVisibility(raw string) (entity.HistoryVisibilityPreference, error) {
	value := strings.TrimSpace(strings.ToLower(raw))
	switch entity.HistoryVisibilityPreference(value) {
	case entity.HistoryVisibilityPublic, entity.HistoryVisibilityPrivate:
		return entity.HistoryVisibilityPreference(value), nil
	default:
		return "", fmt.Errorf("%w: invalid history_visibility_preference", ErrValidation)
	}
}

func toAccountState(raw string) (entity.AccountState, error) {
	value := strings.TrimSpace(strings.ToLower(raw))
	switch entity.AccountState(value) {
	case entity.AccountStateActive, entity.AccountStateDeactivated, entity.AccountStateBanned:
		return entity.AccountState(value), nil
	default:
		return "", fmt.Errorf("%w: invalid account_state", ErrValidation)
	}
}

func toVIPAction(raw string) (entity.VIPAction, error) {
	value := strings.TrimSpace(strings.ToLower(raw))
	switch entity.VIPAction(value) {
	case entity.VIPActionActivate, entity.VIPActionFreeze, entity.VIPActionResume, entity.VIPActionDeactivate:
		return entity.VIPAction(value), nil
	default:
		return "", fmt.Errorf("%w: invalid vip action", ErrValidation)
	}
}
