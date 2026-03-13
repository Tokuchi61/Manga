package repository

import (
	"context"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/history/entity"
	"github.com/google/uuid"
)

func (s *MemoryStore) UpsertCheckpoint(_ context.Context, checkpoint entity.Checkpoint, dedupKey string) (entity.LibraryEntry, string, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	dedupKey = normalizeValue(dedupKey)
	if dedupKey != "" {
		timelineEventID, exists := s.checkpointDedup[dedupKey]
		if exists {
			key := libraryKey(checkpoint.UserID, checkpoint.MangaID)
			if existing, ok := s.libraryByKey[key]; ok {
				return cloneLibraryEntry(existing), timelineEventID, false, nil
			}
			delete(s.checkpointDedup, dedupKey)
		}
	}

	userID := normalizeValue(checkpoint.UserID)
	mangaID := normalizeValue(checkpoint.MangaID)
	chapterID := normalizeValue(checkpoint.ChapterID)
	eventName := normalizeValue(checkpoint.Event)
	if eventName == "" {
		eventName = "chapter.read.checkpoint"
	}

	now := time.Now().UTC()
	occurredAt := checkpoint.OccurredAt.UTC()
	if occurredAt.IsZero() {
		occurredAt = now
	}

	key := libraryKey(userID, mangaID)
	entry, exists := s.libraryByKey[key]
	if !exists {
		entry = entity.LibraryEntry{
			ID:         uuid.NewString(),
			UserID:     userID,
			MangaID:    mangaID,
			Status:     entity.ReadingStatusInProgress,
			CreatedAt:  now,
			UpdatedAt:  now,
			LastReadAt: occurredAt,
		}
	}

	entry.LastChapterID = chapterID
	if checkpoint.PageNumber < 0 {
		checkpoint.PageNumber = 0
	}
	if checkpoint.PageCount < 0 {
		checkpoint.PageCount = 0
	}
	if eventName == "chapter.read.finished" && checkpoint.PageCount > 0 {
		entry.LastPageNumber = checkpoint.PageCount
	} else {
		entry.LastPageNumber = checkpoint.PageNumber
	}
	entry.PageCount = checkpoint.PageCount
	entry.LastReadAt = occurredAt
	entry.UpdatedAt = now

	switch eventName {
	case "chapter.read.finished":
		entry.Status = entity.ReadingStatusCompleted
	case "chapter.read.started", "chapter.read.checkpoint":
		entry.Status = entity.ReadingStatusInProgress
	default:
		entry.Status = entity.ReadingStatusInProgress
	}

	s.libraryByKey[key] = cloneLibraryEntry(entry)

	timelineEventID := uuid.NewString()
	timelineEvent := entity.TimelineEvent{
		ID:            timelineEventID,
		UserID:        userID,
		MangaID:       mangaID,
		ChapterID:     chapterID,
		Event:         eventName,
		PageNumber:    entry.LastPageNumber,
		PageCount:     checkpoint.PageCount,
		RequestID:     normalizeValue(checkpoint.RequestID),
		CorrelationID: normalizeValue(checkpoint.CorrelationID),
		OccurredAt:    occurredAt,
		CreatedAt:     now,
	}
	userTimeline := s.timelineByUser[userID]
	userTimeline = append(userTimeline, cloneTimelineEvent(timelineEvent))
	s.timelineByUser[userID] = userTimeline

	if dedupKey != "" {
		s.checkpointDedup[dedupKey] = timelineEventID
	}

	return cloneLibraryEntry(entry), timelineEventID, true, nil
}
