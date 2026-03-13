package repository

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/history/entity"
)

func (s *MemoryStore) GetLibraryEntry(_ context.Context, userID string, mangaID string) (entity.LibraryEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.libraryByKey[libraryKey(userID, mangaID)]
	if !ok {
		return entity.LibraryEntry{}, ErrNotFound
	}
	return cloneLibraryEntry(entry), nil
}

func (s *MemoryStore) UpsertLibraryEntry(_ context.Context, entry entity.LibraryEntry) (entity.LibraryEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := libraryKey(entry.UserID, entry.MangaID)
	stored := cloneLibraryEntry(entry)
	stored.UserID = normalizeValue(entry.UserID)
	stored.MangaID = normalizeValue(entry.MangaID)
	if stored.ID == "" {
		stored.ID = key
	}
	s.libraryByKey[key] = stored
	return cloneLibraryEntry(stored), nil
}

func (s *MemoryStore) ListContinueReading(_ context.Context, query ContinueReadingQuery) ([]entity.LibraryEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	userID := normalizeValue(query.UserID)
	items := make([]entity.LibraryEntry, 0)
	for _, entry := range s.libraryByKey {
		if normalizeValue(entry.UserID) != userID {
			continue
		}
		if entry.Status != entity.ReadingStatusInProgress {
			continue
		}
		items = append(items, cloneLibraryEntry(entry))
	}

	sortLibrary(items, query.SortBy, "last_read")
	return applyOffsetLimitEntries(items, query.Offset, query.Limit), nil
}

func (s *MemoryStore) ListLibrary(_ context.Context, query LibraryQuery) ([]entity.LibraryEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	userID := normalizeValue(query.UserID)
	status := normalizeValue(query.Status)
	items := make([]entity.LibraryEntry, 0)
	for _, entry := range s.libraryByKey {
		if normalizeValue(entry.UserID) != userID {
			continue
		}
		if status != "" && normalizeValue(string(entry.Status)) != status {
			continue
		}
		if query.Bookmarked != nil && entry.Bookmarked != *query.Bookmarked {
			continue
		}
		if query.Favorited != nil && entry.Favorited != *query.Favorited {
			continue
		}
		if query.SharedOnly && !entry.SharePublic {
			continue
		}
		items = append(items, cloneLibraryEntry(entry))
	}

	sortLibrary(items, query.SortBy, "updated")
	return applyOffsetLimitEntries(items, query.Offset, query.Limit), nil
}
