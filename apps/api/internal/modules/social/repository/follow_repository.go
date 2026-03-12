package repository

import (
	"context"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/social/entity"
)

func (s *MemoryStore) UpsertFollow(_ context.Context, follow entity.FollowRelation, dedupKey string) (entity.FollowRelation, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	dedupKey = normalizeValue(dedupKey)
	if dedupKey != "" {
		storedKey, exists := s.followDedup[dedupKey]
		if exists {
			existing, ok := s.followsByKey[storedKey]
			if ok {
				return cloneFollow(existing), false, nil
			}
			delete(s.followDedup, dedupKey)
		}
	}

	key := directionalKey(follow.FollowerUserID, follow.FolloweeUserID)
	if existing, exists := s.followsByKey[key]; exists {
		return cloneFollow(existing), false, nil
	}

	stored := cloneFollow(follow)
	stored.FollowerUserID = normalizeValue(follow.FollowerUserID)
	stored.FolloweeUserID = normalizeValue(follow.FolloweeUserID)
	stored.RequestID = strings.TrimSpace(follow.RequestID)
	s.followsByKey[key] = stored

	if dedupKey != "" {
		s.followDedup[dedupKey] = key
	}

	return cloneFollow(stored), true, nil
}

func (s *MemoryStore) DeleteFollow(_ context.Context, followerUserID string, followeeUserID string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := directionalKey(followerUserID, followeeUserID)
	if _, exists := s.followsByKey[key]; !exists {
		return false, nil
	}
	delete(s.followsByKey, key)
	return true, nil
}

func (s *MemoryStore) ListFollowers(_ context.Context, userID string) ([]entity.FollowRelation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	target := normalizeValue(userID)
	items := make([]entity.FollowRelation, 0)
	for _, follow := range s.followsByKey {
		if follow.FolloweeUserID != target {
			continue
		}
		items = append(items, cloneFollow(follow))
	}

	sortByTimeDesc(items, func(i int, j int) bool {
		if items[i].UpdatedAt.Equal(items[j].UpdatedAt) {
			return directionalKey(items[i].FollowerUserID, items[i].FolloweeUserID) < directionalKey(items[j].FollowerUserID, items[j].FolloweeUserID)
		}
		return items[i].UpdatedAt.After(items[j].UpdatedAt)
	})

	return items, nil
}

func (s *MemoryStore) ListFollowing(_ context.Context, userID string) ([]entity.FollowRelation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	actor := normalizeValue(userID)
	items := make([]entity.FollowRelation, 0)
	for _, follow := range s.followsByKey {
		if follow.FollowerUserID != actor {
			continue
		}
		items = append(items, cloneFollow(follow))
	}

	sortByTimeDesc(items, func(i int, j int) bool {
		if items[i].UpdatedAt.Equal(items[j].UpdatedAt) {
			return directionalKey(items[i].FollowerUserID, items[i].FolloweeUserID) < directionalKey(items[j].FollowerUserID, items[j].FolloweeUserID)
		}
		return items[i].UpdatedAt.After(items[j].UpdatedAt)
	})

	return items, nil
}
