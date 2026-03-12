package service

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/access/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/access/entity"
	accessrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/access/repository"
	"github.com/google/uuid"
)

func (s *AccessService) CreatePolicyRule(ctx context.Context, request dto.CreatePolicyRuleRequest) (dto.CreatePolicyRuleResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.CreatePolicyRuleResponse{}, err
	}

	effect, err := toPolicyEffect(request.Effect)
	if err != nil {
		return dto.CreatePolicyRuleResponse{}, err
	}
	audienceKind, err := toAudienceKind(request.AudienceKind)
	if err != nil {
		return dto.CreatePolicyRuleResponse{}, err
	}
	scopeKind, err := toScopeKind(request.ScopeKind)
	if err != nil {
		return dto.CreatePolicyRuleResponse{}, err
	}

	now := s.now().UTC()
	key := normalizeName(request.Key)
	audienceSelector := normalizeOptional(request.AudienceSelector, "-")
	scopeSelector := normalizeOptional(request.ScopeSelector, "-")
	active := true
	if request.Active != nil {
		active = *request.Active
	}

	existing, err := s.store.ListPolicyRulesByKey(ctx, key)
	if err != nil {
		return dto.CreatePolicyRuleResponse{}, err
	}
	version := 1
	for _, item := range existing {
		if item.Version >= version {
			version = item.Version + 1
		}
	}

	rule := entity.PolicyRule{
		ID:               uuid.New(),
		Key:              key,
		Effect:           effect,
		AudienceKind:     audienceKind,
		AudienceSelector: audienceSelector,
		ScopeKind:        scopeKind,
		ScopeSelector:    scopeSelector,
		Active:           active,
		Version:          version,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := s.store.CreatePolicyRule(ctx, rule); err != nil {
		if errors.Is(err, accessrepository.ErrConflict) {
			return dto.CreatePolicyRuleResponse{}, ErrPolicyConflict
		}
		return dto.CreatePolicyRuleResponse{}, err
	}

	return dto.CreatePolicyRuleResponse{
		PolicyRuleID:     rule.ID.String(),
		Key:              rule.Key,
		Effect:           string(rule.Effect),
		AudienceKind:     rule.AudienceKind,
		AudienceSelector: rule.AudienceSelector,
		ScopeKind:        rule.ScopeKind,
		ScopeSelector:    rule.ScopeSelector,
		Active:           rule.Active,
		Version:          rule.Version,
		CreatedAt:        rule.CreatedAt,
	}, nil
}
