package service

import (
	"fmt"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/access/entity"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/settings"
	"github.com/google/uuid"
)

func parseUUID(raw string, fieldName string) (uuid.UUID, error) {
	parsed, err := uuid.Parse(strings.TrimSpace(raw))
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: invalid %s", ErrValidation, fieldName)
	}
	return parsed, nil
}

func normalizeName(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

func normalizeOptional(raw string, fallback string) string {
	value := normalizeName(raw)
	if value == "" {
		return fallback
	}
	return value
}

func toPolicyEffect(raw string) (entity.PolicyEffect, error) {
	effect := entity.PolicyEffect(normalizeName(raw))
	if !entity.IsValidPolicyEffect(effect) {
		return "", fmt.Errorf("%w: invalid policy effect", ErrValidation)
	}
	return effect, nil
}

func toAudienceKind(raw string) (string, error) {
	value := normalizeName(raw)
	for _, audience := range settings.AllAudiences {
		if string(audience) == value {
			return value, nil
		}
	}
	return "", fmt.Errorf("%w: invalid audience kind", ErrValidation)
}

func toScopeKind(raw string) (string, error) {
	value := normalizeName(raw)
	for _, scopeKind := range settings.AllScopeKinds {
		if string(scopeKind) == value {
			return value, nil
		}
	}
	return "", fmt.Errorf("%w: invalid scope kind", ErrValidation)
}
