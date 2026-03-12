package user

import "errors"

var (
	ErrValidation             = errors.New("user_validation_failed")
	ErrUserNotFound           = errors.New("user_not_found")
	ErrUserAlreadyExists      = errors.New("user_already_exists")
	ErrProfileNotVisible      = errors.New("user_profile_not_visible")
	ErrInvalidStateTransition = errors.New("user_invalid_state_transition")
)
