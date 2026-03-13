package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/history/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/history/entity"
)

func (s *HistoryService) IngestChapterSignal(ctx context.Context, request dto.IngestChapterSignalRequest) (dto.IngestChapterSignalResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.IngestChapterSignalResponse{}, err
	}
	if s.chapterSignalProvider == nil {
		return dto.IngestChapterSignalResponse{}, ErrChapterSignalUnavailable
	}

	userID, err := parseID(request.UserID, "user_id")
	if err != nil {
		return dto.IngestChapterSignalResponse{}, err
	}
	chapterID, err := parseID(request.ChapterID, "chapter_id")
	if err != nil {
		return dto.IngestChapterSignalResponse{}, err
	}

	anchor, err := s.chapterSignalProvider.GetResumeAnchor(ctx, chapterID, request.PageNumber)
	if err != nil {
		return dto.IngestChapterSignalResponse{}, fmt.Errorf("%w: %v", ErrChapterSignalInvalid, err)
	}
	mangaID, err := parseID(anchor.MangaID, "manga_id")
	if err != nil {
		return dto.IngestChapterSignalResponse{}, fmt.Errorf("%w: %v", ErrChapterSignalInvalid, err)
	}

	eventName := parseChapterEvent(request.Event)
	signal := s.chapterSignalProvider.BuildReadSignal(
		chapterID,
		mangaID,
		anchor.PageNumber,
		anchor.PageCount,
		eventName,
		strings.TrimSpace(request.RequestID),
		strings.TrimSpace(request.CorrelationID),
	)

	dedupKey := buildCheckpointDedupKey(userID, signal)
	entry, timelineEventID, created, err := s.store.UpsertCheckpoint(ctx, entity.Checkpoint{
		UserID:        userID,
		MangaID:       mangaID,
		ChapterID:     chapterID,
		Event:         signal.Event,
		PageNumber:    signal.PageNumber,
		PageCount:     signal.PageCount,
		RequestID:     signal.RequestID,
		CorrelationID: signal.CorrelationID,
		OccurredAt:    signal.OccurredAt,
	}, dedupKey)
	if err != nil {
		return dto.IngestChapterSignalResponse{}, err
	}

	return dto.IngestChapterSignalResponse{
		LibraryEntryID:  entry.ID,
		TimelineEventID: timelineEventID,
		MangaID:         entry.MangaID,
		ChapterID:       entry.LastChapterID,
		Event:           signal.Event,
		Status:          string(entry.Status),
		Created:         created,
	}, nil
}
