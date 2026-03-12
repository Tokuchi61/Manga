package repository

import (
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/entity"
)

// MemoryStore is a stage-8 bootstrap persistence for chapter flows.
type MemoryStore struct {
	mu sync.RWMutex

	chaptersByID map[string]entity.Chapter
	slugIndex    map[string]map[string]string
	sequenceIndex map[string]map[int]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		chaptersByID: make(map[string]entity.Chapter),
		slugIndex: make(map[string]map[string]string),
		sequenceIndex: make(map[string]map[int]string),
	}
}

func normalizeValue(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

func cloneTimePtr(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	v := value.UTC()
	return &v
}

func clonePages(pages []entity.ChapterPage) []entity.ChapterPage {
	if len(pages) == 0 {
		return nil
	}
	result := make([]entity.ChapterPage, len(pages))
	copy(result, pages)
	return result
}

func cloneChapter(in entity.Chapter) entity.Chapter {
	out := in
	out.ScheduledAt = cloneTimePtr(in.ScheduledAt)
	out.PublishedAt = cloneTimePtr(in.PublishedAt)
	out.ArchivedAt = cloneTimePtr(in.ArchivedAt)
	out.DeletedAt = cloneTimePtr(in.DeletedAt)
	out.Pages = clonePages(in.Pages)
	return out
}

func sortChapters(items []entity.Chapter, sortBy string) {
	sortKey := normalizeValue(sortBy)
	switch sortKey {
	case "sequence_oldest":
		sort.Slice(items, func(i, j int) bool {
			if items[i].SequenceNo == items[j].SequenceNo {
				return items[i].CreatedAt.Before(items[j].CreatedAt)
			}
			return items[i].SequenceNo < items[j].SequenceNo
		})
	case "published_oldest":
		sort.Slice(items, func(i, j int) bool {
			left := items[i].PublishedAt
			right := items[j].PublishedAt
			if left == nil && right == nil {
				return items[i].SequenceNo < items[j].SequenceNo
			}
			if left == nil {
				return false
			}
			if right == nil {
				return true
			}
			if left.Equal(*right) {
				return items[i].SequenceNo < items[j].SequenceNo
			}
			return left.Before(*right)
		})
	case "published_newest":
		sort.Slice(items, func(i, j int) bool {
			left := items[i].PublishedAt
			right := items[j].PublishedAt
			if left == nil && right == nil {
				return items[i].SequenceNo > items[j].SequenceNo
			}
			if left == nil {
				return false
			}
			if right == nil {
				return true
			}
			if left.Equal(*right) {
				return items[i].SequenceNo > items[j].SequenceNo
			}
			return left.After(*right)
		})
	default:
		sort.Slice(items, func(i, j int) bool {
			if items[i].SequenceNo == items[j].SequenceNo {
				return items[i].CreatedAt.After(items[j].CreatedAt)
			}
			return items[i].SequenceNo > items[j].SequenceNo
		})
	}
}
