package repository

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/entity"
)

func (s *MemoryStore) CreateChapter(_ context.Context, chapter entity.Chapter) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	mangaID := normalizeValue(chapter.MangaID)
	slugKey := normalizeValue(chapter.Slug)

	if _, ok := s.slugIndex[mangaID]; !ok {
		s.slugIndex[mangaID] = make(map[string]string)
	}
	if _, exists := s.slugIndex[mangaID][slugKey]; exists {
		return ErrConflict
	}

	if _, ok := s.sequenceIndex[mangaID]; !ok {
		s.sequenceIndex[mangaID] = make(map[int]string)
	}
	if _, exists := s.sequenceIndex[mangaID][chapter.SequenceNo]; exists {
		return ErrConflict
	}

	stored := cloneChapter(chapter)
	stored.MangaID = mangaID
	stored.Slug = slugKey

	s.chaptersByID[stored.ID] = stored
	s.slugIndex[mangaID][slugKey] = stored.ID
	s.sequenceIndex[mangaID][stored.SequenceNo] = stored.ID
	return nil
}

func (s *MemoryStore) GetChapterByID(_ context.Context, chapterID string) (entity.Chapter, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	chapter, ok := s.chaptersByID[chapterID]
	if !ok {
		return entity.Chapter{}, ErrNotFound
	}
	return cloneChapter(chapter), nil
}

func (s *MemoryStore) GetChapterBySlug(_ context.Context, mangaID string, slug string) (entity.Chapter, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	mangaKey := normalizeValue(mangaID)
	slugKey := normalizeValue(slug)
	chapterID, ok := s.slugIndex[mangaKey][slugKey]
	if !ok {
		return entity.Chapter{}, ErrNotFound
	}
	chapter, ok := s.chaptersByID[chapterID]
	if !ok {
		return entity.Chapter{}, ErrNotFound
	}
	return cloneChapter(chapter), nil
}

func (s *MemoryStore) UpdateChapter(_ context.Context, chapter entity.Chapter) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, ok := s.chaptersByID[chapter.ID]
	if !ok {
		return ErrNotFound
	}

	mangaID := normalizeValue(chapter.MangaID)
	slugKey := normalizeValue(chapter.Slug)

	if _, ok := s.slugIndex[mangaID]; !ok {
		s.slugIndex[mangaID] = make(map[string]string)
	}
	if ownerID, exists := s.slugIndex[mangaID][slugKey]; exists && ownerID != chapter.ID {
		return ErrConflict
	}

	if _, ok := s.sequenceIndex[mangaID]; !ok {
		s.sequenceIndex[mangaID] = make(map[int]string)
	}
	if ownerID, exists := s.sequenceIndex[mangaID][chapter.SequenceNo]; exists && ownerID != chapter.ID {
		return ErrConflict
	}

	delete(s.slugIndex[normalizeValue(current.MangaID)], normalizeValue(current.Slug))
	delete(s.sequenceIndex[normalizeValue(current.MangaID)], current.SequenceNo)

	next := cloneChapter(chapter)
	next.MangaID = mangaID
	next.Slug = slugKey

	s.chaptersByID[chapter.ID] = next
	s.slugIndex[mangaID][slugKey] = chapter.ID
	s.sequenceIndex[mangaID][chapter.SequenceNo] = chapter.ID
	return nil
}

func (s *MemoryStore) ListChaptersByManga(_ context.Context, query ListQuery) ([]entity.Chapter, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	mangaID := normalizeValue(query.MangaID)
	result := make([]entity.Chapter, 0)
	for _, chapter := range s.chaptersByID {
		if normalizeValue(chapter.MangaID) != mangaID {
			continue
		}
		if !query.IncludeDeleted && chapter.DeletedAt != nil {
			continue
		}
		if !query.IncludeUnpublished && chapter.PublishState != entity.PublishStatePublished {
			continue
		}
		result = append(result, cloneChapter(chapter))
	}

	sortChapters(result, query.SortBy)

	offset := query.Offset
	if offset < 0 {
		offset = 0
	}
	if offset >= len(result) {
		return []entity.Chapter{}, nil
	}

	limit := query.Limit
	if limit <= 0 {
		limit = 20
	}
	end := offset + limit
	if end > len(result) {
		end = len(result)
	}

	return append([]entity.Chapter(nil), result[offset:end]...), nil
}
