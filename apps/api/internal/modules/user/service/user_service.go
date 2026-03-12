package service

import (
	"errors"
	"fmt"
	"time"

	userrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/user/repository"
	uservalidator "github.com/Tokuchi61/Manga/apps/api/internal/modules/user/validator"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
)

var (
	ErrValidation             = errors.New("user_validation_failed")
	ErrUserNotFound           = errors.New("user_not_found")
	ErrUserAlreadyExists      = errors.New("user_already_exists")
	ErrProfileNotVisible      = errors.New("user_profile_not_visible")
	ErrInvalidStateTransition = errors.New("user_invalid_state_transition")
	ErrAccountDeactivated     = errors.New("user_account_deactivated")
	ErrAccountBanned          = errors.New("user_account_banned")
)

// UserService owns stage-5 user account/profile/membership flows.
type UserService struct {
	store     userrepository.Store
	validator *validation.Validator
	now       func() time.Time
}

func New(store userrepository.Store, validator *validation.Validator) *UserService {
	if store == nil {
		store = userrepository.NewMemoryStore()
	}
	return &UserService{
		store:     store,
		validator: validator,
		now:       time.Now,
	}
}

func (s *UserService) validateInput(payload any) error {
	if err := uservalidator.ValidateStruct(s.validator, payload); err != nil {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}
	return nil
}
