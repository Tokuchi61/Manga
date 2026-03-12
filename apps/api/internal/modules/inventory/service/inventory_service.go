package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/inventory/repository"
	inventoryvalidator "github.com/Tokuchi61/Manga/apps/api/internal/modules/inventory/validator"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
)

var (
	ErrValidation           = errors.New("inventory_validation_failed")
	ErrNotFound             = errors.New("inventory_not_found")
	ErrForbiddenAction      = errors.New("inventory_forbidden_action")
	ErrConflict             = errors.New("inventory_conflict")
	ErrReadDisabled         = errors.New("inventory_read_disabled")
	ErrClaimDisabled        = errors.New("inventory_claim_disabled")
	ErrConsumeDisabled      = errors.New("inventory_consume_disabled")
	ErrEquipDisabled        = errors.New("inventory_equip_disabled")
	ErrInsufficientQuantity = errors.New("inventory_insufficient_quantity")
	ErrItemAlreadyOwned     = errors.New("inventory_item_already_owned")
)

// InventoryService owns stage-15 inventory flows.
type InventoryService struct {
	store     repository.Store
	validator *validation.Validator
	now       func() time.Time
}

func New(store repository.Store, validator *validation.Validator) *InventoryService {
	if store == nil {
		store = repository.NewMemoryStore()
	}
	return &InventoryService{
		store:     store,
		validator: validator,
		now:       time.Now,
	}
}

func (s *InventoryService) validateInput(payload any) error {
	if err := inventoryvalidator.ValidateStruct(s.validator, payload); err != nil {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}
	return nil
}

func (s *InventoryService) requireReadEnabled(cfgEnabled bool) error {
	if !cfgEnabled {
		return ErrReadDisabled
	}
	return nil
}

func (s *InventoryService) requireClaimEnabled(cfgEnabled bool) error {
	if !cfgEnabled {
		return ErrClaimDisabled
	}
	return nil
}

func (s *InventoryService) requireConsumeEnabled(cfgEnabled bool) error {
	if !cfgEnabled {
		return ErrConsumeDisabled
	}
	return nil
}

func (s *InventoryService) requireEquipEnabled(cfgEnabled bool) error {
	if !cfgEnabled {
		return ErrEquipDisabled
	}
	return nil
}
