package repository

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/social/entity"
)

func (s *MemoryStore) CreateWallPost(_ context.Context, post entity.WallPost, dedupKey string) (entity.WallPost, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	dedupKey = normalizeValue(dedupKey)
	if dedupKey != "" {
		storedID, exists := s.wallPostDedup[dedupKey]
		if exists {
			existing, ok := s.wallPostsByID[storedID]
			if ok {
				return cloneWallPost(existing), false, nil
			}
			delete(s.wallPostDedup, dedupKey)
		}
	}

	if _, exists := s.wallPostsByID[post.ID]; exists {
		return entity.WallPost{}, false, ErrConflict
	}

	stored := cloneWallPost(post)
	stored.OwnerUserID = normalizeValue(post.OwnerUserID)
	s.wallPostsByID[stored.ID] = stored

	if dedupKey != "" {
		s.wallPostDedup[dedupKey] = stored.ID
	}

	return cloneWallPost(stored), true, nil
}

func (s *MemoryStore) GetWallPostByID(_ context.Context, postID string) (entity.WallPost, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	post, ok := s.wallPostsByID[postID]
	if !ok {
		return entity.WallPost{}, ErrNotFound
	}
	return cloneWallPost(post), nil
}

func (s *MemoryStore) CreateWallReply(_ context.Context, reply entity.WallReply, dedupKey string) (entity.WallReply, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.wallPostsByID[reply.PostID]; !exists {
		return entity.WallReply{}, false, ErrNotFound
	}

	dedupKey = normalizeValue(dedupKey)
	if dedupKey != "" {
		storedID, exists := s.wallReplyDedup[dedupKey]
		if exists {
			existing, ok := s.wallRepliesByID[storedID]
			if ok {
				return cloneWallReply(existing), false, nil
			}
			delete(s.wallReplyDedup, dedupKey)
		}
	}

	if _, exists := s.wallRepliesByID[reply.ID]; exists {
		return entity.WallReply{}, false, ErrConflict
	}

	stored := cloneWallReply(reply)
	stored.PostID = normalizeValue(reply.PostID)
	stored.OwnerUserID = normalizeValue(reply.OwnerUserID)
	s.wallRepliesByID[stored.ID] = stored
	s.wallRepliesByPost[stored.PostID] = append(s.wallRepliesByPost[stored.PostID], stored.ID)

	if dedupKey != "" {
		s.wallReplyDedup[dedupKey] = stored.ID
	}

	return cloneWallReply(stored), true, nil
}

func (s *MemoryStore) ListWallPosts(_ context.Context, ownerUserID string, sortBy string, limit int, offset int) ([]entity.WallPost, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	owner := normalizeValue(ownerUserID)
	items := make([]entity.WallPost, 0)
	for _, post := range s.wallPostsByID {
		if post.OwnerUserID != owner {
			continue
		}
		items = append(items, cloneWallPost(post))
	}

	if normalizeValue(sortBy) == "oldest" {
		sortByTimeDesc(items, func(i int, j int) bool {
			if items[i].CreatedAt.Equal(items[j].CreatedAt) {
				return items[i].ID < items[j].ID
			}
			return items[i].CreatedAt.Before(items[j].CreatedAt)
		})
	} else {
		sortByTimeDesc(items, func(i int, j int) bool {
			if items[i].CreatedAt.Equal(items[j].CreatedAt) {
				return items[i].ID < items[j].ID
			}
			return items[i].CreatedAt.After(items[j].CreatedAt)
		})
	}

	return applyOffsetLimit(items, offset, limit), nil
}

func (s *MemoryStore) ListWallRepliesByPostID(_ context.Context, postID string, sortBy string) ([]entity.WallReply, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ids := s.wallRepliesByPost[normalizeValue(postID)]
	items := make([]entity.WallReply, 0, len(ids))
	for _, id := range ids {
		reply, ok := s.wallRepliesByID[id]
		if !ok {
			continue
		}
		items = append(items, cloneWallReply(reply))
	}

	if normalizeValue(sortBy) == "oldest" {
		sortByTimeDesc(items, func(i int, j int) bool {
			if items[i].CreatedAt.Equal(items[j].CreatedAt) {
				return items[i].ID < items[j].ID
			}
			return items[i].CreatedAt.Before(items[j].CreatedAt)
		})
	} else {
		sortByTimeDesc(items, func(i int, j int) bool {
			if items[i].CreatedAt.Equal(items[j].CreatedAt) {
				return items[i].ID < items[j].ID
			}
			return items[i].CreatedAt.After(items[j].CreatedAt)
		})
	}

	return items, nil
}
