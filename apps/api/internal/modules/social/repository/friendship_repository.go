package repository

import (
	"context"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/social/entity"
)

func (s *MemoryStore) CreateFriendRequest(_ context.Context, request entity.FriendshipRequest, dedupKey string) (entity.FriendshipRequest, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	dedupKey = normalizeValue(dedupKey)
	if dedupKey != "" {
		storedKey, exists := s.friendRequestDedup[dedupKey]
		if exists {
			existing, ok := s.friendRequestsByKey[storedKey]
			if ok {
				return cloneFriendRequest(existing), false, nil
			}
			delete(s.friendRequestDedup, dedupKey)
		}
	}

	key := directionalKey(request.RequesterUserID, request.TargetUserID)
	if existing, exists := s.friendRequestsByKey[key]; exists {
		return cloneFriendRequest(existing), false, nil
	}

	stored := cloneFriendRequest(request)
	stored.RequesterUserID = normalizeValue(request.RequesterUserID)
	stored.TargetUserID = normalizeValue(request.TargetUserID)
	stored.RequestID = strings.TrimSpace(request.RequestID)
	s.friendRequestsByKey[key] = stored
	if dedupKey != "" {
		s.friendRequestDedup[dedupKey] = key
	}

	return cloneFriendRequest(stored), true, nil
}

func (s *MemoryStore) GetFriendRequest(_ context.Context, requesterUserID string, targetUserID string) (entity.FriendshipRequest, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	request, ok := s.friendRequestsByKey[directionalKey(requesterUserID, targetUserID)]
	if !ok {
		return entity.FriendshipRequest{}, ErrNotFound
	}
	return cloneFriendRequest(request), nil
}

func (s *MemoryStore) DeleteFriendRequest(_ context.Context, requesterUserID string, targetUserID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.friendRequestsByKey, directionalKey(requesterUserID, targetUserID))
	return nil
}

func (s *MemoryStore) ListFriendRequests(_ context.Context, userID string, direction string) ([]entity.FriendshipRequest, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	actor := normalizeValue(userID)
	direction = normalizeValue(direction)
	items := make([]entity.FriendshipRequest, 0)

	for _, request := range s.friendRequestsByKey {
		switch direction {
		case "incoming":
			if request.TargetUserID != actor {
				continue
			}
		case "outgoing":
			if request.RequesterUserID != actor {
				continue
			}
		default:
			if request.TargetUserID != actor && request.RequesterUserID != actor {
				continue
			}
		}
		items = append(items, cloneFriendRequest(request))
	}

	sortByTimeDesc(items, func(i int, j int) bool {
		if items[i].UpdatedAt.Equal(items[j].UpdatedAt) {
			return directionalKey(items[i].RequesterUserID, items[i].TargetUserID) < directionalKey(items[j].RequesterUserID, items[j].TargetUserID)
		}
		return items[i].UpdatedAt.After(items[j].UpdatedAt)
	})

	return items, nil
}

func (s *MemoryStore) UpsertFriendship(_ context.Context, friendship entity.Friendship) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	stored := cloneFriendship(friendship)
	stored.UserAID = normalizeValue(friendship.UserAID)
	stored.UserBID = normalizeValue(friendship.UserBID)
	s.friendshipsByPair[pairKey(stored.UserAID, stored.UserBID)] = stored
	return nil
}

func (s *MemoryStore) DeleteFriendship(_ context.Context, userAID string, userBID string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := pairKey(userAID, userBID)
	if _, exists := s.friendshipsByPair[key]; !exists {
		return false, nil
	}
	delete(s.friendshipsByPair, key)
	return true, nil
}

func (s *MemoryStore) ListFriends(_ context.Context, userID string) ([]entity.Friendship, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	actor := normalizeValue(userID)
	items := make([]entity.Friendship, 0)
	for _, friendship := range s.friendshipsByPair {
		if friendship.UserAID != actor && friendship.UserBID != actor {
			continue
		}
		items = append(items, cloneFriendship(friendship))
	}

	sortByTimeDesc(items, func(i int, j int) bool {
		if items[i].UpdatedAt.Equal(items[j].UpdatedAt) {
			return pairKey(items[i].UserAID, items[i].UserBID) < pairKey(items[j].UserAID, items[j].UserBID)
		}
		return items[i].UpdatedAt.After(items[j].UpdatedAt)
	})

	return items, nil
}

func (s *MemoryStore) AreFriends(_ context.Context, userAID string, userBID string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, exists := s.friendshipsByPair[pairKey(userAID, userBID)]
	return exists, nil
}
