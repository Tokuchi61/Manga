package service

import (
	"fmt"
	"strings"

	chaptercontract "github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/contract"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/history/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/history/entity"
	"github.com/google/uuid"
)

func normalizeValue(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

func parseID(raw string, fieldName string) (string, error) {
	parsed, err := uuid.Parse(strings.TrimSpace(raw))
	if err != nil {
		return "", fmt.Errorf("%w: invalid %s", ErrValidation, fieldName)
	}
	return parsed.String(), nil
}

func parseSortBy(raw string, fallback string) string {
	value := normalizeValue(raw)
	switch value {
	case "newest", "oldest":
		return value
	default:
		if fallback == "" {
			return "newest"
		}
		return fallback
	}
}

func parseOptionalStatus(raw string) (entity.ReadingStatus, error) {
	value := normalizeValue(raw)
	if value == "" {
		return "", nil
	}
	switch entity.ReadingStatus(value) {
	case entity.ReadingStatusInProgress,
		entity.ReadingStatusCompleted,
		entity.ReadingStatusDropped:
		return entity.ReadingStatus(value), nil
	default:
		return "", fmt.Errorf("%w: invalid status", ErrValidation)
	}
}

func parseChapterEvent(raw string) string {
	value := normalizeValue(raw)
	switch value {
	case chaptercontract.EventReadStarted,
		chaptercontract.EventReadCheckpoint,
		chaptercontract.EventReadFinished:
		return value
	default:
		return chaptercontract.EventReadCheckpoint
	}
}

func parseOptionalTimelineEvent(raw string) string {
	value := normalizeValue(raw)
	switch value {
	case "":
		return ""
	case chaptercontract.EventReadStarted,
		chaptercontract.EventReadCheckpoint,
		chaptercontract.EventReadFinished:
		return value
	default:
		return ""
	}
}

func pointerBoolFromQuery(raw string) *bool {
	value := normalizeValue(raw)
	switch value {
	case "true":
		v := true
		return &v
	case "false":
		v := false
		return &v
	default:
		return nil
	}
}

func buildCheckpointDedupKey(userID string, signal chaptercontract.ReadSignal) string {
	eventName := normalizeValue(signal.Event)
	requestID := normalizeValue(signal.RequestID)
	correlationID := normalizeValue(signal.CorrelationID)
	chapterID := normalizeValue(signal.ChapterID)
	return fmt.Sprintf("%s:%s:%s:%d:%d:%s:%s", normalizeValue(userID), eventName, chapterID, signal.PageNumber, signal.PageCount, requestID, correlationID)
}

func mapEntry(entry entity.LibraryEntry, includeUser bool) dto.HistoryEntryResponse {
	response := dto.HistoryEntryResponse{
		LibraryEntryID: entry.ID,
		MangaID:        entry.MangaID,
		LastChapterID:  entry.LastChapterID,
		LastPageNumber: entry.LastPageNumber,
		PageCount:      entry.PageCount,
		Status:         string(entry.Status),
		Bookmarked:     entry.Bookmarked,
		Favorited:      entry.Favorited,
		SharePublic:    entry.SharePublic,
		LastReadAt:     entry.LastReadAt,
		CreatedAt:      entry.CreatedAt,
		UpdatedAt:      entry.UpdatedAt,
	}
	if includeUser {
		response.UserID = entry.UserID
	}
	return response
}

func mapTimeline(event entity.TimelineEvent) dto.TimelineItemResponse {
	return dto.TimelineItemResponse{
		TimelineEventID: event.ID,
		MangaID:         event.MangaID,
		ChapterID:       event.ChapterID,
		Event:           event.Event,
		PageNumber:      event.PageNumber,
		PageCount:       event.PageCount,
		RequestID:       event.RequestID,
		CorrelationID:   event.CorrelationID,
		OccurredAt:      event.OccurredAt,
	}
}

func mapRuntimeConfig(cfg entity.RuntimeConfig) dto.RuntimeConfigResponse {
	return dto.RuntimeConfigResponse{
		ContinueReadingEnabled: cfg.ContinueReadingEnabled,
		LibraryEnabled:         cfg.LibraryEnabled,
		TimelineEnabled:        cfg.TimelineEnabled,
		BookmarkWriteEnabled:   cfg.BookmarkWriteEnabled,
		UpdatedAt:              cfg.UpdatedAt,
	}
}
