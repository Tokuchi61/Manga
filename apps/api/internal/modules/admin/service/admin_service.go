package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/admin/repository"
	adminvalidator "github.com/Tokuchi61/Manga/apps/api/internal/modules/admin/validator"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
)

var (
	ErrValidation                 = errors.New("admin_validation_failed")
	ErrNotFound                   = errors.New("admin_not_found")
	ErrConflict                   = errors.New("admin_conflict")
	ErrForbiddenAction            = errors.New("admin_forbidden_action")
	ErrDoubleConfirmationRequired = errors.New("admin_double_confirmation_required")
)

// AdminService owns stage-21 admin flows.
type AdminService struct {
	store     repository.Store
	validator *validation.Validator
	now       func() time.Time
}

func New(store repository.Store, validator *validation.Validator) *AdminService {
	if store == nil {
		store = repository.NewMemoryStore()
	}
	return &AdminService{
		store:     store,
		validator: validator,
		now:       time.Now,
	}
}

func (s *AdminService) validateInput(payload any) error {
	if err := adminvalidator.ValidateStruct(s.validator, payload); err != nil {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}
	return nil
}
