package repository

import (
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/entity"
)

// MemoryStore is a stage-7 bootstrap persistence for manga flows.
type MemoryStore struct {
	mu sync.RWMutex

	mangaByID map[string]entity.Manga
	slugIndex map[string]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		mangaByID: make(map[string]entity.Manga),
		slugIndex: make(map[string]string),
	}
}

func normalizeValue(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

func normalizeStringSlice(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	result := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		normalized := normalizeValue(value)
		if normalized == "" {
			continue
		}
		if _, exists := seen[normalized]; exists {
			continue
		}
		seen[normalized] = struct{}{}
		result = append(result, normalized)
	}
	sort.Strings(result)
	if len(result) == 0 {
		return nil
	}
	return result
}

func cloneTimePtr(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	v := value.UTC()
	return &v
}

func cloneManga(in entity.Manga) entity.Manga {
	out := in
	out.AlternativeTitles = append([]string(nil), in.AlternativeTitles...)
	out.Genres = append([]string(nil), in.Genres...)
	out.Tags = append([]string(nil), in.Tags...)
	out.Themes = append([]string(nil), in.Themes...)
	out.ContentWarnings = append([]string(nil), in.ContentWarnings...)
	out.CollectionKeys = append([]string(nil), in.CollectionKeys...)
	out.ScheduledAt = cloneTimePtr(in.ScheduledAt)
	out.PublishedAt = cloneTimePtr(in.PublishedAt)
	out.ArchivedAt = cloneTimePtr(in.ArchivedAt)
	out.DeletedAt = cloneTimePtr(in.DeletedAt)
	return out
}

func containsValue(values []string, target string) bool {
	if target == "" {
		return true
	}
	target = normalizeValue(target)
	for _, value := range values {
		if normalizeValue(value) == target {
			return true
		}
	}
	return false
}

func matchesSearch(manga entity.Manga, search string) bool {
	search = normalizeValue(search)
	if search == "" {
		return true
	}
	if strings.Contains(normalizeValue(manga.Title), search) {
		return true
	}
	if strings.Contains(normalizeValue(manga.Summary), search) {
		return true
	}
	if strings.Contains(normalizeValue(manga.ShortSummary), search) {
		return true
	}
	if strings.Contains(normalizeValue(manga.Slug), search) {
		return true
	}
	for _, title := range manga.AlternativeTitles {
		if strings.Contains(normalizeValue(title), search) {
			return true
		}
	}
	return false
}
