package repository

import (
	"context"
	"sort"

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

func (s *MemoryStore) ResolveNavigation(_ context.Context, mangaID string, chapterID string) (NavigationResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	mangaKey := normalizeValue(mangaID)
	sequenceMap, ok := s.sequenceIndex[mangaKey]
	if !ok || len(sequenceMap) == 0 {
		return NavigationResult{}, nil
	}

	sequences := make([]int, 0, len(sequenceMap))
	for sequenceNo := range sequenceMap {
		sequences = append(sequences, sequenceNo)
	}
	sort.Ints(sequences)

	orderedIDs := make([]string, 0, len(sequences))
	for _, sequenceNo := range sequences {
		id := sequenceMap[sequenceNo]
		chapter, exists := s.chaptersByID[id]
		if !exists {
			continue
		}
		if chapter.DeletedAt != nil || chapter.PublishState != entity.PublishStatePublished {
			continue
		}
		orderedIDs = append(orderedIDs, chapter.ID)
	}

	if len(orderedIDs) == 0 {
		return NavigationResult{}, nil
	}

	result := NavigationResult{
		FirstID: stringPtr(orderedIDs[0]),
		LastID:  stringPtr(orderedIDs[len(orderedIDs)-1]),
	}

	currentIndex := -1
	for i := range orderedIDs {
		if orderedIDs[i] == chapterID {
			currentIndex = i
			break
		}
	}
	if currentIndex == -1 {
		return result, nil
	}

	result.FoundCurrent = true
	if currentIndex > 0 {
		result.PreviousID = stringPtr(orderedIDs[currentIndex-1])
	}
	if currentIndex+1 < len(orderedIDs) {
		result.NextID = stringPtr(orderedIDs[currentIndex+1])
	}

	return result, nil
}

func stringPtr(value string) *string {
	resolved := value
	return &resolved
}
