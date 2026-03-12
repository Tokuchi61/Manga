package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	supportcontract "github.com/Tokuchi61/Manga/apps/api/internal/modules/support/contract"
	moderationrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/repository"
	moderationvalidator "github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/validator"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
)

var (
	ErrValidation              = errors.New("moderation_validation_failed")
	ErrCaseNotFound            = errors.New("moderation_case_not_found")
	ErrCaseAlreadyExists       = errors.New("moderation_case_already_exists")
	ErrForbiddenAction         = errors.New("moderation_forbidden_action")
	ErrInvalidStateTransition  = errors.New("moderation_invalid_state_transition")
	ErrSupportReferenceRequired = errors.New("moderation_support_reference_required")
	ErrSupportReferenceInvalid = errors.New("moderation_support_reference_invalid")
	ErrAlreadyEscalated        = errors.New("moderation_already_escalated")
)

// SupportHandoffLookup exposes support->moderation reference lookup boundary.
type SupportHandoffLookup interface {
	GetModerationHandoffReference(ctx context.Context, supportID string, requestID string, correlationID string) (supportcontract.ModerationHandoffReference, error)
}

// SupportCaseLinker exposes support record linking boundary.
type SupportCaseLinker interface {
	LinkModerationCase(ctx context.Context, supportID string, moderationCaseID string) error
}

// ModerationService owns stage-11 moderation queue and case flows.
type ModerationService struct {
	store         moderationrepository.Store
	validator     *validation.Validator
	now           func() time.Time
	supportLookup SupportHandoffLookup
	supportLinker SupportCaseLinker
}

func New(store moderationrepository.Store, validator *validation.Validator) *ModerationService {
	if store == nil {
		store = moderationrepository.NewMemoryStore()
	}
	return &ModerationService{
		store:     store,
		validator: validator,
		now:       time.Now,
	}
}

func (s *ModerationService) SetSupportContracts(lookup SupportHandoffLookup, linker SupportCaseLinker) {
	s.supportLookup = lookup
	s.supportLinker = linker
}

func (s *ModerationService) validateInput(payload any) error {
	if err := moderationvalidator.ValidateStruct(s.validator, payload); err != nil {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}
	return nil
}
