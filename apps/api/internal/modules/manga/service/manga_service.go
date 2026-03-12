package service

import (
	"errors"
	"fmt"
	"time"

	mangarepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/repository"
	mangavalidator "github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/validator"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
)

var (
	ErrValidation             = errors.New("manga_validation_failed")
	ErrMangaNotFound          = errors.New("manga_not_found")
	ErrMangaAlreadyExists     = errors.New("manga_already_exists")
	ErrMangaNotVisible        = errors.New("manga_not_visible")
	ErrInvalidStateTransition = errors.New("manga_invalid_state_transition")
	ErrInvalidCounterValue    = errors.New("manga_invalid_counter_value")
)

// MangaService owns stage-7 manga metadata/discovery/public surfaces.
type MangaService struct {
	store     mangarepository.Store
	validator *validation.Validator
	now       func() time.Time
}

func New(store mangarepository.Store, validator *validation.Validator) *MangaService {
	if store == nil {
		store = mangarepository.NewMemoryStore()
	}
	return &MangaService{
		store:     store,
		validator: validator,
		now:       time.Now,
	}
}

func (s *MangaService) validateInput(payload any) error {
	if err := mangavalidator.ValidateStruct(s.validator, payload); err != nil {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}
	return nil
}
