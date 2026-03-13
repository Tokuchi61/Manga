package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/shop/repository"
	shopvalidator "github.com/Tokuchi61/Manga/apps/api/internal/modules/shop/validator"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
)

var (
	ErrValidation        = errors.New("shop_validation_failed")
	ErrNotFound          = errors.New("shop_not_found")
	ErrConflict          = errors.New("shop_conflict")
	ErrForbiddenAction   = errors.New("shop_forbidden_action")
	ErrCatalogDisabled   = errors.New("shop_catalog_disabled")
	ErrPurchaseDisabled  = errors.New("shop_purchase_disabled")
	ErrCampaignDisabled  = errors.New("shop_campaign_disabled")
	ErrAlreadyOwned      = errors.New("shop_already_owned")
	ErrEligibilityFailed = errors.New("shop_eligibility_failed")
)

// ShopService owns stage-18 shop flows.
type ShopService struct {
	store     repository.Store
	validator *validation.Validator
	now       func() time.Time
}

func New(store repository.Store, validator *validation.Validator) *ShopService {
	if store == nil {
		store = repository.NewMemoryStore()
	}
	return &ShopService{
		store:     store,
		validator: validator,
		now:       time.Now,
	}
}

func (s *ShopService) validateInput(payload any) error {
	if err := shopvalidator.ValidateStruct(s.validator, payload); err != nil {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}
	return nil
}

func (s *ShopService) requireCatalogEnabled(enabled bool) error {
	if !enabled {
		return ErrCatalogDisabled
	}
	return nil
}

func (s *ShopService) requirePurchaseEnabled(enabled bool) error {
	if !enabled {
		return ErrPurchaseDisabled
	}
	return nil
}

func (s *ShopService) requireCampaignEnabled(enabled bool) error {
	if !enabled {
		return ErrCampaignDisabled
	}
	return nil
}
