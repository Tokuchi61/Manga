package service

import (
	"errors"
	"fmt"
	"time"

	supportrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/support/repository"
	supportvalidator "github.com/Tokuchi61/Manga/apps/api/internal/modules/support/validator"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
)

var (
	ErrValidation                  = errors.New("support_validation_failed")
	ErrSupportNotFound             = errors.New("support_not_found")
	ErrSupportAlreadyExists        = errors.New("support_already_exists")
	ErrSupportNotVisible           = errors.New("support_not_visible")
	ErrForbiddenAction             = errors.New("support_forbidden_action")
	ErrDuplicateDetected           = errors.New("support_duplicate_detected")
	ErrInvalidStateTransition      = errors.New("support_invalid_state_transition")
	ErrTargetRequired              = errors.New("support_target_required")
	ErrInvalidTarget               = errors.New("support_invalid_target")
	ErrModerationHandoffNotAllowed = errors.New("support_moderation_handoff_not_allowed")
	ErrAlreadyHandedOff            = errors.New("support_already_handed_off")
)

// SupportService owns stage-10 support intake/review flows.
type SupportService struct {
	store           supportrepository.Store
	validator       *validation.Validator
	now             func() time.Time
	duplicateWindow time.Duration
}

func New(store supportrepository.Store, validator *validation.Validator) *SupportService {
	if store == nil {
		store = supportrepository.NewMemoryStore()
	}
	return &SupportService{
		store:           store,
		validator:       validator,
		now:             time.Now,
		duplicateWindow: 10 * time.Minute,
	}
}

func (s *SupportService) validateInput(payload any) error {
	if err := supportvalidator.ValidateStruct(s.validator, payload); err != nil {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}
	return nil
}
