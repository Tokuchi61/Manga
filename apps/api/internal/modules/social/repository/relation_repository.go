package repository

import (
	"context"
	"strings"
	"time"
)

func (s *MemoryStore) SetBlock(_ context.Context, actorUserID string, targetUserID string, enabled bool, updatedAt time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := directionalKey(actorUserID, targetUserID)
	if enabled {
		s.blockByKey[key] = updatedAt.UTC()
		return nil
	}
	delete(s.blockByKey, key)
	return nil
}

func (s *MemoryStore) SetMute(_ context.Context, actorUserID string, targetUserID string, enabled bool, updatedAt time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := directionalKey(actorUserID, targetUserID)
	if enabled {
		s.muteByKey[key] = updatedAt.UTC()
		return nil
	}
	delete(s.muteByKey, key)
	return nil
}

func (s *MemoryStore) SetRestrict(_ context.Context, actorUserID string, targetUserID string, enabled bool, updatedAt time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := directionalKey(actorUserID, targetUserID)
	if enabled {
		s.restrictByKey[key] = updatedAt.UTC()
		return nil
	}
	delete(s.restrictByKey, key)
	return nil
}

func (s *MemoryStore) ListBlocked(_ context.Context, actorUserID string) ([]RelationEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return listRelationsForActor(s.blockByKey, actorUserID), nil
}

func (s *MemoryStore) ListMuted(_ context.Context, actorUserID string) ([]RelationEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return listRelationsForActor(s.muteByKey, actorUserID), nil
}

func (s *MemoryStore) ListRestricted(_ context.Context, actorUserID string) ([]RelationEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return listRelationsForActor(s.restrictByKey, actorUserID), nil
}

func (s *MemoryStore) IsBlockedEither(_ context.Context, userAID string, userBID string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, left := s.blockByKey[directionalKey(userAID, userBID)]
	if left {
		return true, nil
	}
	_, right := s.blockByKey[directionalKey(userBID, userAID)]
	return right, nil
}

func (s *MemoryStore) IsRestrictedEither(_ context.Context, userAID string, userBID string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, left := s.restrictByKey[directionalKey(userAID, userBID)]
	if left {
		return true, nil
	}
	_, right := s.restrictByKey[directionalKey(userBID, userAID)]
	return right, nil
}

func (s *MemoryStore) IsMutedEither(_ context.Context, userAID string, userBID string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, left := s.muteByKey[directionalKey(userAID, userBID)]
	if left {
		return true, nil
	}
	_, right := s.muteByKey[directionalKey(userBID, userAID)]
	return right, nil
}

func listRelationsForActor(source map[string]time.Time, actorUserID string) []RelationEntry {
	actor := normalizeValue(actorUserID)
	items := make([]RelationEntry, 0)
	for key, updatedAt := range source {
		parts := strings.SplitN(key, ":", 2)
		if len(parts) != 2 {
			continue
		}
		if parts[0] != actor {
			continue
		}
		items = append(items, RelationEntry{
			ActorUserID:  parts[0],
			TargetUserID: parts[1],
			UpdatedAt:    updatedAt.UTC(),
		})
	}

	sortByTimeDesc(items, func(i int, j int) bool {
		if items[i].UpdatedAt.Equal(items[j].UpdatedAt) {
			return directionalKey(items[i].ActorUserID, items[i].TargetUserID) < directionalKey(items[j].ActorUserID, items[j].TargetUserID)
		}
		return items[i].UpdatedAt.After(items[j].UpdatedAt)
	})

	return items
}
