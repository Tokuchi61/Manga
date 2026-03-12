package service

import (
	"errors"
	"fmt"
	"time"

	commentrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/comment/repository"
	commentvalidator "github.com/Tokuchi61/Manga/apps/api/internal/modules/comment/validator"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
)

var (
	ErrValidation             = errors.New("comment_validation_failed")
	ErrCommentNotFound        = errors.New("comment_not_found")
	ErrCommentAlreadyExists   = errors.New("comment_already_exists")
	ErrCommentNotVisible      = errors.New("comment_not_visible")
	ErrEditWindowExpired      = errors.New("comment_edit_window_expired")
	ErrReplyDepthExceeded     = errors.New("comment_reply_depth_exceeded")
	ErrCommentLocked          = errors.New("comment_locked")
	ErrRateLimited            = errors.New("comment_rate_limited")
	ErrRestoreWindowExpired   = errors.New("comment_restore_window_expired")
	ErrForbiddenAction        = errors.New("comment_forbidden_action")
	ErrInvalidStateTransition = errors.New("comment_invalid_state_transition")
)

// CommentService owns stage-9 comment/thread flows.
type CommentService struct {
	store         commentrepository.Store
	validator     *validation.Validator
	now           func() time.Time
	writeCooldown time.Duration
	editWindow    time.Duration
	restoreWindow time.Duration
	maxReplyDepth int
}

func New(store commentrepository.Store, validator *validation.Validator) *CommentService {
	if store == nil {
		store = commentrepository.NewMemoryStore()
	}
	return &CommentService{
		store:         store,
		validator:     validator,
		now:           time.Now,
		writeCooldown: 15 * time.Second,
		editWindow:    30 * time.Minute,
		restoreWindow: 72 * time.Hour,
		maxReplyDepth: 3,
	}
}

func (s *CommentService) validateInput(payload any) error {
	if err := commentvalidator.ValidateStruct(s.validator, payload); err != nil {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}
	return nil
}
