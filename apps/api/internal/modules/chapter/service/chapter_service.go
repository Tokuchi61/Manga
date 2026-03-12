package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	chapterrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/repository"
	chaptervalidator "github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/validator"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
)

var (
	ErrValidation             = errors.New("chapter_validation_failed")
	ErrChapterNotFound        = errors.New("chapter_not_found")
	ErrChapterAlreadyExists   = errors.New("chapter_already_exists")
	ErrMangaNotFound          = errors.New("chapter_manga_not_found")
	ErrChapterNotVisible      = errors.New("chapter_not_visible")
	ErrInvalidStateTransition = errors.New("chapter_invalid_state_transition")
	ErrInvalidMediaState      = errors.New("chapter_invalid_media_state")
)

// MangaLookup exposes manga target existence checks.
type MangaLookup interface {
	TargetExists(ctx context.Context, mangaID string) (bool, error)
}

// ChapterService owns stage-8 chapter read/release flows.
type ChapterService struct {
	store       chapterrepository.Store
	validator   *validation.Validator
	mangaLookup MangaLookup
	now         func() time.Time
}

func New(store chapterrepository.Store, validator *validation.Validator) *ChapterService {
	if store == nil {
		store = chapterrepository.NewMemoryStore()
	}
	return &ChapterService{
		store:     store,
		validator: validator,
		now:       time.Now,
	}
}

func (s *ChapterService) SetMangaLookup(lookup MangaLookup) {
	s.mangaLookup = lookup
}

func (s *ChapterService) validateInput(payload any) error {
	if err := chaptervalidator.ValidateStruct(s.validator, payload); err != nil {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}
	return nil
}
