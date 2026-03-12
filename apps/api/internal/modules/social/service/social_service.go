package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/social/repository"
	socialvalidator "github.com/Tokuchi61/Manga/apps/api/internal/modules/social/validator"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
)

var (
	ErrValidation         = errors.New("social_validation_failed")
	ErrNotFound           = errors.New("social_not_found")
	ErrForbiddenAction    = errors.New("social_forbidden_action")
	ErrConflict           = errors.New("social_conflict")
	ErrFriendshipDisabled = errors.New("social_friendship_disabled")
	ErrFollowDisabled     = errors.New("social_follow_disabled")
	ErrWallDisabled       = errors.New("social_wall_disabled")
	ErrMessagingDisabled  = errors.New("social_messaging_disabled")
	ErrAlreadyFriends     = errors.New("social_already_friends")
	ErrThreadAccessDenied = errors.New("social_thread_access_denied")
)

// SocialService owns stage-14 social flows.
type SocialService struct {
	store     repository.Store
	validator *validation.Validator
	now       func() time.Time
}

func New(store repository.Store, validator *validation.Validator) *SocialService {
	if store == nil {
		store = repository.NewMemoryStore()
	}
	return &SocialService{
		store:     store,
		validator: validator,
		now:       time.Now,
	}
}

func (s *SocialService) validateInput(payload any) error {
	if err := socialvalidator.ValidateStruct(s.validator, payload); err != nil {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}
	return nil
}

func (s *SocialService) requireFriendshipWriteEnabled(cfgEnabled bool) error {
	if !cfgEnabled {
		return ErrFriendshipDisabled
	}
	return nil
}

func (s *SocialService) requireFollowWriteEnabled(cfgEnabled bool) error {
	if !cfgEnabled {
		return ErrFollowDisabled
	}
	return nil
}

func (s *SocialService) requireWallEnabled(cfgEnabled bool) error {
	if !cfgEnabled {
		return ErrWallDisabled
	}
	return nil
}

func (s *SocialService) requireMessagingWriteEnabled(cfgEnabled bool) error {
	if !cfgEnabled {
		return ErrMessagingDisabled
	}
	return nil
}
