package repository

import (
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/history/entity"
)

// MemoryStore is stage-13 bootstrap persistence for history flows.
type MemoryStore struct {
	mu sync.RWMutex

	libraryByKey    map[string]entity.LibraryEntry
	timelineByUser  map[string][]entity.TimelineEvent
	checkpointDedup map[string]string
	runtimeConfig   entity.RuntimeConfig
}

func NewMemoryStore() *MemoryStore {
	now := time.Now().UTC()
	return &MemoryStore{
		libraryByKey:    make(map[string]entity.LibraryEntry),
		timelineByUser:  make(map[string][]entity.TimelineEvent),
		checkpointDedup: make(map[string]string),
		runtimeConfig: entity.RuntimeConfig{
			ContinueReadingEnabled: true,
			LibraryEnabled:         true,
			TimelineEnabled:        true,
			BookmarkWriteEnabled:   true,
			UpdatedAt:              now,
		},
	}
}

func normalizeValue(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

func libraryKey(userID string, mangaID string) string {
	return normalizeValue(userID) + ":" + normalizeValue(mangaID)
}

func cloneLibraryEntry(in entity.LibraryEntry) entity.LibraryEntry {
	out := in
	out.LastReadAt = in.LastReadAt.UTC()
	out.CreatedAt = in.CreatedAt.UTC()
	out.UpdatedAt = in.UpdatedAt.UTC()
	return out
}

func cloneTimelineEvent(in entity.TimelineEvent) entity.TimelineEvent {
	out := in
	out.OccurredAt = in.OccurredAt.UTC()
	out.CreatedAt = in.CreatedAt.UTC()
	return out
}

func cloneRuntimeConfig(in entity.RuntimeConfig) entity.RuntimeConfig {
	out := in
	out.UpdatedAt = in.UpdatedAt.UTC()
	return out
}

func sortLibrary(items []entity.LibraryEntry, sortBy string, field string) {
	switch normalizeValue(sortBy) {
	case "oldest":
		sort.Slice(items, func(i, j int) bool {
			left := items[i]
			right := items[j]
			if field == "updated" {
				if left.UpdatedAt.Equal(right.UpdatedAt) {
					return left.ID < right.ID
				}
				return left.UpdatedAt.Before(right.UpdatedAt)
			}
			if left.LastReadAt.Equal(right.LastReadAt) {
				return left.ID < right.ID
			}
			return left.LastReadAt.Before(right.LastReadAt)
		})
	default:
		sort.Slice(items, func(i, j int) bool {
			left := items[i]
			right := items[j]
			if field == "updated" {
				if left.UpdatedAt.Equal(right.UpdatedAt) {
					return left.ID < right.ID
				}
				return left.UpdatedAt.After(right.UpdatedAt)
			}
			if left.LastReadAt.Equal(right.LastReadAt) {
				return left.ID < right.ID
			}
			return left.LastReadAt.After(right.LastReadAt)
		})
	}
}

func sortTimeline(items []entity.TimelineEvent, sortBy string) {
	switch normalizeValue(sortBy) {
	case "oldest":
		sort.Slice(items, func(i, j int) bool {
			if items[i].OccurredAt.Equal(items[j].OccurredAt) {
				return items[i].ID < items[j].ID
			}
			return items[i].OccurredAt.Before(items[j].OccurredAt)
		})
	default:
		sort.Slice(items, func(i, j int) bool {
			if items[i].OccurredAt.Equal(items[j].OccurredAt) {
				return items[i].ID < items[j].ID
			}
			return items[i].OccurredAt.After(items[j].OccurredAt)
		})
	}
}

func applyOffsetLimitEntries(items []entity.LibraryEntry, offset int, limit int) []entity.LibraryEntry {
	if offset < 0 {
		offset = 0
	}
	if offset >= len(items) {
		return []entity.LibraryEntry{}
	}
	if limit <= 0 {
		limit = 50
	}
	end := offset + limit
	if end > len(items) {
		end = len(items)
	}
	return append([]entity.LibraryEntry(nil), items[offset:end]...)
}

func applyOffsetLimitTimeline(items []entity.TimelineEvent, offset int, limit int) []entity.TimelineEvent {
	if offset < 0 {
		offset = 0
	}
	if offset >= len(items) {
		return []entity.TimelineEvent{}
	}
	if limit <= 0 {
		limit = 50
	}
	end := offset + limit
	if end > len(items) {
		end = len(items)
	}
	return append([]entity.TimelineEvent(nil), items[offset:end]...)
}
