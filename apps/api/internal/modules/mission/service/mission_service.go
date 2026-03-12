package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/mission/repository"
	missionvalidator "github.com/Tokuchi61/Manga/apps/api/internal/modules/mission/validator"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
)

var (
	ErrValidation             = errors.New("mission_validation_failed")
	ErrNotFound               = errors.New("mission_not_found")
	ErrConflict               = errors.New("mission_conflict")
	ErrForbiddenAction        = errors.New("mission_forbidden_action")
	ErrReadDisabled           = errors.New("mission_read_disabled")
	ErrClaimDisabled          = errors.New("mission_claim_disabled")
	ErrProgressIngestDisabled = errors.New("mission_progress_ingest_disabled")
	ErrMissionInactive        = errors.New("mission_inactive")
	ErrMissionNotEligible     = errors.New("mission_not_eligible")
	ErrAlreadyClaimed         = errors.New("mission_already_claimed")
)

// MissionService owns stage-16 mission flows.
type MissionService struct {
	store     repository.Store
	validator *validation.Validator
	now       func() time.Time
}

func New(store repository.Store, validator *validation.Validator) *MissionService {
	if store == nil {
		store = repository.NewMemoryStore()
	}
	return &MissionService{
		store:     store,
		validator: validator,
		now:       time.Now,
	}
}

func (s *MissionService) validateInput(payload any) error {
	if err := missionvalidator.ValidateStruct(s.validator, payload); err != nil {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}
	return nil
}

func (s *MissionService) requireReadEnabled(cfgEnabled bool) error {
	if !cfgEnabled {
		return ErrReadDisabled
	}
	return nil
}

func (s *MissionService) requireClaimEnabled(cfgEnabled bool) error {
	if !cfgEnabled {
		return ErrClaimDisabled
	}
	return nil
}

func (s *MissionService) requireProgressIngestEnabled(cfgEnabled bool) error {
	if !cfgEnabled {
		return ErrProgressIngestDisabled
	}
	return nil
}
