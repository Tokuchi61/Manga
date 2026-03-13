package service

import (
	"context"
	"errors"
	"testing"
	"time"

	chaptercontract "github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/contract"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/history/dto"
	historyrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/history/repository"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type chapterSignalProviderStub struct {
	anchor chaptercontract.ResumeAnchor
	err    error
}

func (s chapterSignalProviderStub) GetResumeAnchor(_ context.Context, _ string, pageNumber int) (chaptercontract.ResumeAnchor, error) {
	if s.err != nil {
		return chaptercontract.ResumeAnchor{}, s.err
	}
	anchor := s.anchor
	if pageNumber > 0 {
		anchor.PageNumber = pageNumber
	}
	return anchor, nil
}

func (s chapterSignalProviderStub) BuildReadSignal(chapterID string, mangaID string, pageNumber int, pageCount int, event string, requestID string, correlationID string) chaptercontract.ReadSignal {
	if event == "" {
		event = chaptercontract.EventReadCheckpoint
	}
	return chaptercontract.ReadSignal{
		Event:         event,
		ChapterID:     chapterID,
		MangaID:       mangaID,
		PageNumber:    pageNumber,
		PageCount:     pageCount,
		OccurredAt:    time.Now().UTC(),
		RequestID:     requestID,
		CorrelationID: correlationID,
	}
}

func TestIngestChapterSignalAndLibraryFlow(t *testing.T) {
	store := historyrepository.NewMemoryStore()
	svc := New(store, validation.New())

	userID := uuid.NewString()
	mangaID := uuid.NewString()
	chapterID := uuid.NewString()
	now := time.Date(2026, 3, 13, 0, 10, 0, 0, time.UTC)
	svc.now = func() time.Time { return now }

	svc.SetChapterSignalProvider(chapterSignalProviderStub{anchor: chaptercontract.ResumeAnchor{
		ChapterID:  chapterID,
		MangaID:    mangaID,
		PageNumber: 7,
		PageCount:  20,
		UpdatedAt:  now,
	}})

	created, err := svc.IngestChapterSignal(context.Background(), dto.IngestChapterSignalRequest{
		UserID:        userID,
		ChapterID:     chapterID,
		Event:         chaptercontract.EventReadCheckpoint,
		PageNumber:    7,
		RequestID:     "req-stage13-intake-1",
		CorrelationID: "corr-stage13-intake-1",
	})
	require.NoError(t, err)
	require.True(t, created.Created)
	require.Equal(t, "in_progress", created.Status)

	again, err := svc.IngestChapterSignal(context.Background(), dto.IngestChapterSignalRequest{
		UserID:        userID,
		ChapterID:     chapterID,
		Event:         chaptercontract.EventReadCheckpoint,
		PageNumber:    7,
		RequestID:     "req-stage13-intake-1",
		CorrelationID: "corr-stage13-intake-1",
	})
	require.NoError(t, err)
	require.False(t, again.Created)
	require.Equal(t, created.LibraryEntryID, again.LibraryEntryID)

	continueReading, err := svc.ListContinueReading(context.Background(), dto.ListContinueReadingRequest{UserID: userID})
	require.NoError(t, err)
	require.Equal(t, 1, continueReading.Count)

	timeline, err := svc.ListTimeline(context.Background(), dto.ListTimelineRequest{UserID: userID})
	require.NoError(t, err)
	require.Equal(t, 1, timeline.Count)

	bookmark, err := svc.UpdateBookmark(context.Background(), dto.UpdateBookmarkRequest{
		UserID:     userID,
		MangaID:    mangaID,
		Bookmarked: true,
		Favorited:  true,
	})
	require.NoError(t, err)
	require.True(t, bookmark.Bookmarked)
	require.True(t, bookmark.Favorited)

	shared, err := svc.UpdateShare(context.Background(), dto.UpdateShareRequest{UserID: userID, MangaID: mangaID, SharePublic: true})
	require.NoError(t, err)
	require.True(t, shared.SharePublic)

	publicLibrary, err := svc.ListPublicLibrary(context.Background(), dto.ListPublicLibraryRequest{OwnerUserID: userID})
	require.NoError(t, err)
	require.Equal(t, 1, publicLibrary.Count)
}

func TestRuntimeControlsAffectSurfaces(t *testing.T) {
	store := historyrepository.NewMemoryStore()
	svc := New(store, validation.New())

	userID := uuid.NewString()
	mangaID := uuid.NewString()

	_, err := svc.UpdateContinueReadingState(context.Background(), dto.UpdateContinueReadingStateRequest{Enabled: false})
	require.NoError(t, err)

	_, err = svc.ListContinueReading(context.Background(), dto.ListContinueReadingRequest{UserID: userID})
	require.True(t, errors.Is(err, ErrContinueReadingDisabled))

	_, err = svc.UpdateBookmarkWriteState(context.Background(), dto.UpdateBookmarkWriteStateRequest{Enabled: false})
	require.NoError(t, err)

	_, err = svc.UpdateBookmark(context.Background(), dto.UpdateBookmarkRequest{UserID: userID, MangaID: mangaID, Bookmarked: true})
	require.True(t, errors.Is(err, ErrBookmarkWriteDisabled))
}
