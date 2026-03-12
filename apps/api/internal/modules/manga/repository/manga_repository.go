package repository

import (
	"context"
	"sort"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/entity"
)

func (s *MemoryStore) CreateManga(_ context.Context, manga entity.Manga) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	slugKey := normalizeValue(manga.Slug)
	if _, exists := s.slugIndex[slugKey]; exists {
		return ErrConflict
	}

	stored := cloneManga(manga)
	stored.Slug = slugKey
	stored.AlternativeTitles = normalizeStringSlice(stored.AlternativeTitles)
	stored.Genres = normalizeStringSlice(stored.Genres)
	stored.Tags = normalizeStringSlice(stored.Tags)
	stored.Themes = normalizeStringSlice(stored.Themes)
	stored.ContentWarnings = normalizeStringSlice(stored.ContentWarnings)
	stored.CollectionKeys = normalizeStringSlice(stored.CollectionKeys)

	s.mangaByID[stored.ID] = stored
	s.slugIndex[slugKey] = stored.ID
	return nil
}

func (s *MemoryStore) GetMangaByID(_ context.Context, mangaID string) (entity.Manga, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	manga, ok := s.mangaByID[mangaID]
	if !ok {
		return entity.Manga{}, ErrNotFound
	}
	return cloneManga(manga), nil
}

func (s *MemoryStore) GetMangaBySlug(_ context.Context, slug string) (entity.Manga, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	mangaID, ok := s.slugIndex[normalizeValue(slug)]
	if !ok {
		return entity.Manga{}, ErrNotFound
	}
	manga, ok := s.mangaByID[mangaID]
	if !ok {
		return entity.Manga{}, ErrNotFound
	}
	return cloneManga(manga), nil
}

func (s *MemoryStore) UpdateManga(_ context.Context, manga entity.Manga) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, ok := s.mangaByID[manga.ID]
	if !ok {
		return ErrNotFound
	}

	next := cloneManga(manga)
	next.Slug = normalizeValue(next.Slug)
	next.AlternativeTitles = normalizeStringSlice(next.AlternativeTitles)
	next.Genres = normalizeStringSlice(next.Genres)
	next.Tags = normalizeStringSlice(next.Tags)
	next.Themes = normalizeStringSlice(next.Themes)
	next.ContentWarnings = normalizeStringSlice(next.ContentWarnings)
	next.CollectionKeys = normalizeStringSlice(next.CollectionKeys)

	if currentSlug := normalizeValue(current.Slug); currentSlug != next.Slug {
		if mappedID, exists := s.slugIndex[next.Slug]; exists && mappedID != manga.ID {
			return ErrConflict
		}
		delete(s.slugIndex, currentSlug)
		s.slugIndex[next.Slug] = manga.ID
	}

	s.mangaByID[manga.ID] = next
	return nil
}

func (s *MemoryStore) ListManga(_ context.Context, query ListQuery) ([]entity.Manga, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]entity.Manga, 0, len(s.mangaByID))
	for _, manga := range s.mangaByID {
		if !query.IncludeDeleted && manga.DeletedAt != nil {
			continue
		}
		if !query.IncludeUnpublished && manga.PublishState != entity.PublishStatePublished {
			continue
		}
		if !query.IncludeHidden && manga.Visibility != entity.VisibilityPublic {
			continue
		}
		if !matchesSearch(manga, query.Search) {
			continue
		}
		if !containsValue(manga.Genres, query.Genre) {
			continue
		}
		if !containsValue(manga.Tags, query.Tag) {
			continue
		}
		if !containsValue(manga.Themes, query.Theme) {
			continue
		}
		if !containsValue(manga.ContentWarnings, query.ContentWarning) {
			continue
		}
		result = append(result, cloneManga(manga))
	}

	switch normalizeValue(query.SortBy) {
	case "title":
		sort.Slice(result, func(i, j int) bool {
			return normalizeValue(result[i].Title) < normalizeValue(result[j].Title)
		})
	case "popular":
		sort.Slice(result, func(i, j int) bool {
			if result[i].ViewCount == result[j].ViewCount {
				return normalizeValue(result[i].Title) < normalizeValue(result[j].Title)
			}
			return result[i].ViewCount > result[j].ViewCount
		})
	case "updated":
		sort.Slice(result, func(i, j int) bool {
			if result[i].UpdatedAt.Equal(result[j].UpdatedAt) {
				return result[i].CreatedAt.After(result[j].CreatedAt)
			}
			return result[i].UpdatedAt.After(result[j].UpdatedAt)
		})
	default:
		sort.Slice(result, func(i, j int) bool {
			if result[i].CreatedAt.Equal(result[j].CreatedAt) {
				return normalizeValue(result[i].Title) < normalizeValue(result[j].Title)
			}
			return result[i].CreatedAt.After(result[j].CreatedAt)
		})
	}

	offset := query.Offset
	if offset < 0 {
		offset = 0
	}
	if offset >= len(result) {
		return []entity.Manga{}, nil
	}

	limit := query.Limit
	if limit <= 0 {
		limit = 20
	}
	end := offset + limit
	if end > len(result) {
		end = len(result)
	}

	return append([]entity.Manga(nil), result[offset:end]...), nil
}
