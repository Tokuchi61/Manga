package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/ads/repository"
	adsvalidator "github.com/Tokuchi61/Manga/apps/api/internal/modules/ads/validator"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
)

var (
	ErrValidation          = errors.New("ads_validation_failed")
	ErrNotFound            = errors.New("ads_not_found")
	ErrConflict            = errors.New("ads_conflict")
	ErrSurfaceDisabled     = errors.New("ads_surface_disabled")
	ErrPlacementDisabled   = errors.New("ads_placement_disabled")
	ErrCampaignDisabled    = errors.New("ads_campaign_disabled")
	ErrClickIntakeDisabled = errors.New("ads_click_intake_disabled")
)

// AdsService owns stage-20 ads flows.
type AdsService struct {
	store     repository.Store
	validator *validation.Validator
	now       func() time.Time
}

func New(store repository.Store, validator *validation.Validator) *AdsService {
	if store == nil {
		store = repository.NewMemoryStore()
	}
	return &AdsService{
		store:     store,
		validator: validator,
		now:       time.Now,
	}
}

func (s *AdsService) validateInput(payload any) error {
	if err := adsvalidator.ValidateStruct(s.validator, payload); err != nil {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}
	return nil
}

func (s *AdsService) requireSurfaceEnabled(enabled bool) error {
	if !enabled {
		return ErrSurfaceDisabled
	}
	return nil
}

func (s *AdsService) requirePlacementEnabled(enabled bool) error {
	if !enabled {
		return ErrPlacementDisabled
	}
	return nil
}

func (s *AdsService) requireCampaignEnabled(enabled bool) error {
	if !enabled {
		return ErrCampaignDisabled
	}
	return nil
}

func (s *AdsService) requireClickIntakeEnabled(enabled bool) error {
	if !enabled {
		return ErrClickIntakeDisabled
	}
	return nil
}
