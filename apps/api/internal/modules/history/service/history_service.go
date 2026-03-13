package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	chaptercontract "github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/contract"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/history/repository"
	historyvalidator "github.com/Tokuchi61/Manga/apps/api/internal/modules/history/validator"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
)

var (
	ErrValidation               = errors.New("history_validation_failed")
	ErrNotFound                 = errors.New("history_not_found")
	ErrForbiddenAction          = errors.New("history_forbidden_action")
	ErrContinueReadingDisabled  = errors.New("history_continue_reading_disabled")
	ErrLibraryDisabled          = errors.New("history_library_disabled")
	ErrTimelineDisabled         = errors.New("history_timeline_disabled")
	ErrBookmarkWriteDisabled    = errors.New("history_bookmark_write_disabled")
	ErrChapterSignalUnavailable = errors.New("history_chapter_signal_unavailable")
	ErrChapterSignalInvalid     = errors.New("history_chapter_signal_invalid")
)

// ChapterSignalProvider exposes chapter->history signal build boundary.
type ChapterSignalProvider interface {
	GetResumeAnchor(ctx context.Context, chapterID string, pageNumber int) (chaptercontract.ResumeAnchor, error)
	BuildReadSignal(chapterID string, mangaID string, pageNumber int, pageCount int, event string, requestID string, correlationID string) chaptercontract.ReadSignal
}

// HistoryService owns stage-13 history flows.
type HistoryService struct {
	store                 repository.Store
	validator             *validation.Validator
	now                   func() time.Time
	chapterSignalProvider ChapterSignalProvider
}

func New(store repository.Store, validator *validation.Validator) *HistoryService {
	if store == nil {
		store = repository.NewMemoryStore()
	}
	return &HistoryService{
		store:     store,
		validator: validator,
		now:       time.Now,
	}
}

func (s *HistoryService) SetChapterSignalProvider(provider ChapterSignalProvider) {
	s.chapterSignalProvider = provider
}

func (s *HistoryService) validateInput(payload any) error {
	if err := historyvalidator.ValidateStruct(s.validator, payload); err != nil {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}
	return nil
}
