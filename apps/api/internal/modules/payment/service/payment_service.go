package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/repository"
	paymentvalidator "github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/validator"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
)

var (
	ErrValidation              = errors.New("payment_validation_failed")
	ErrNotFound                = errors.New("payment_not_found")
	ErrConflict                = errors.New("payment_conflict")
	ErrForbiddenAction         = errors.New("payment_forbidden_action")
	ErrManaPurchaseDisabled    = errors.New("payment_mana_purchase_disabled")
	ErrCheckoutDisabled        = errors.New("payment_checkout_disabled")
	ErrTransactionReadDisabled = errors.New("payment_transaction_read_disabled")
	ErrCallbackIntakePaused    = errors.New("payment_callback_intake_paused")
)

// PaymentService owns stage-19 payment flows.
type PaymentService struct {
	store     repository.Store
	validator *validation.Validator
	now       func() time.Time
}

func New(store repository.Store, validator *validation.Validator) *PaymentService {
	if store == nil {
		store = repository.NewMemoryStore()
	}
	return &PaymentService{
		store:     store,
		validator: validator,
		now:       time.Now,
	}
}

func (s *PaymentService) validateInput(payload any) error {
	if err := paymentvalidator.ValidateStruct(s.validator, payload); err != nil {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}
	return nil
}

func (s *PaymentService) requireManaPurchaseEnabled(enabled bool) error {
	if !enabled {
		return ErrManaPurchaseDisabled
	}
	return nil
}

func (s *PaymentService) requireCheckoutEnabled(enabled bool) error {
	if !enabled {
		return ErrCheckoutDisabled
	}
	return nil
}

func (s *PaymentService) requireTransactionReadEnabled(enabled bool) error {
	if !enabled {
		return ErrTransactionReadDisabled
	}
	return nil
}

func (s *PaymentService) requireCallbackIntakeOpen(paused bool) error {
	if paused {
		return ErrCallbackIntakePaused
	}
	return nil
}
