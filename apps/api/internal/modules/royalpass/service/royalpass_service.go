package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/royalpass/repository"
	rpvalidator "github.com/Tokuchi61/Manga/apps/api/internal/modules/royalpass/validator"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
)

var (
	ErrValidation         = errors.New("royalpass_validation_failed")
	ErrNotFound           = errors.New("royalpass_not_found")
	ErrConflict           = errors.New("royalpass_conflict")
	ErrForbiddenAction    = errors.New("royalpass_forbidden_action")
	ErrSeasonDisabled     = errors.New("royalpass_season_disabled")
	ErrClaimDisabled      = errors.New("royalpass_claim_disabled")
	ErrPremiumDisabled    = errors.New("royalpass_premium_disabled")
	ErrSeasonUnavailable  = errors.New("royalpass_season_unavailable")
	ErrTierNotEligible    = errors.New("royalpass_tier_not_eligible")
	ErrTierAlreadyClaimed = errors.New("royalpass_tier_already_claimed")
)

// RoyalPassService owns stage-17 royalpass flows.
type RoyalPassService struct {
	store     repository.Store
	validator *validation.Validator
	now       func() time.Time
}

func New(store repository.Store, validator *validation.Validator) *RoyalPassService {
	if store == nil {
		store = repository.NewMemoryStore()
	}
	return &RoyalPassService{
		store:     store,
		validator: validator,
		now:       time.Now,
	}
}

func (s *RoyalPassService) validateInput(payload any) error {
	if err := rpvalidator.ValidateStruct(s.validator, payload); err != nil {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}
	return nil
}

func (s *RoyalPassService) requireSeasonEnabled(enabled bool) error {
	if !enabled {
		return ErrSeasonDisabled
	}
	return nil
}

func (s *RoyalPassService) requireClaimEnabled(enabled bool) error {
	if !enabled {
		return ErrClaimDisabled
	}
	return nil
}

func (s *RoyalPassService) requirePremiumEnabled(enabled bool) error {
	if !enabled {
		return ErrPremiumDisabled
	}
	return nil
}
