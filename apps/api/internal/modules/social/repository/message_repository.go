package repository

import (
	"context"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/social/entity"
)

func (s *MemoryStore) OpenThread(_ context.Context, thread entity.MessageThread) (entity.MessageThread, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := pairKey(thread.UserAID, thread.UserBID)
	if threadID, exists := s.threadByPair[key]; exists {
		existing, ok := s.threadsByID[threadID]
		if ok {
			return cloneThread(existing), false, nil
		}
		delete(s.threadByPair, key)
	}

	if _, exists := s.threadsByID[thread.ID]; exists {
		return entity.MessageThread{}, false, ErrConflict
	}

	stored := cloneThread(thread)
	stored.UserAID = normalizeValue(thread.UserAID)
	stored.UserBID = normalizeValue(thread.UserBID)
	stored.LastMessageID = strings.TrimSpace(thread.LastMessageID)
	s.threadsByID[stored.ID] = stored
	s.threadByPair[key] = stored.ID

	return cloneThread(stored), true, nil
}

func (s *MemoryStore) GetThreadByID(_ context.Context, threadID string) (entity.MessageThread, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	thread, ok := s.threadsByID[strings.TrimSpace(threadID)]
	if !ok {
		return entity.MessageThread{}, ErrNotFound
	}
	return cloneThread(thread), nil
}

func (s *MemoryStore) UpdateThread(_ context.Context, thread entity.MessageThread) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	threadID := strings.TrimSpace(thread.ID)
	if _, exists := s.threadsByID[threadID]; !exists {
		return ErrNotFound
	}

	stored := cloneThread(thread)
	stored.UserAID = normalizeValue(thread.UserAID)
	stored.UserBID = normalizeValue(thread.UserBID)
	stored.LastMessageID = strings.TrimSpace(thread.LastMessageID)
	s.threadsByID[threadID] = stored
	s.threadByPair[pairKey(stored.UserAID, stored.UserBID)] = threadID
	return nil
}

func (s *MemoryStore) ListThreadsByUser(_ context.Context, userID string) ([]entity.MessageThread, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	actor := normalizeValue(userID)
	items := make([]entity.MessageThread, 0)
	for _, thread := range s.threadsByID {
		if thread.UserAID != actor && thread.UserBID != actor {
			continue
		}
		items = append(items, cloneThread(thread))
	}

	sortByTimeDesc(items, func(i int, j int) bool {
		if items[i].UpdatedAt.Equal(items[j].UpdatedAt) {
			return items[i].ID < items[j].ID
		}
		return items[i].UpdatedAt.After(items[j].UpdatedAt)
	})

	return items, nil
}

func (s *MemoryStore) CreateMessage(_ context.Context, message entity.Message, dedupKey string) (entity.Message, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	threadID := strings.TrimSpace(message.ThreadID)
	if _, exists := s.threadsByID[threadID]; !exists {
		return entity.Message{}, false, ErrNotFound
	}

	dedupKey = normalizeValue(dedupKey)
	if dedupKey != "" {
		storedID, exists := s.messageDedup[dedupKey]
		if exists {
			existing, ok := s.messagesByID[storedID]
			if ok {
				return cloneMessage(existing), false, nil
			}
			delete(s.messageDedup, dedupKey)
		}
	}

	if _, exists := s.messagesByID[message.ID]; exists {
		return entity.Message{}, false, ErrConflict
	}

	stored := cloneMessage(message)
	stored.ThreadID = threadID
	stored.SenderUserID = normalizeValue(message.SenderUserID)
	stored.Body = strings.TrimSpace(message.Body)
	stored.RequestID = strings.TrimSpace(message.RequestID)
	stored.CorrelationID = strings.TrimSpace(message.CorrelationID)
	s.messagesByID[stored.ID] = stored
	s.messageIDsByThread[threadID] = append(s.messageIDsByThread[threadID], stored.ID)

	if dedupKey != "" {
		s.messageDedup[dedupKey] = stored.ID
	}

	return cloneMessage(stored), true, nil
}

func (s *MemoryStore) ListMessagesByThreadID(_ context.Context, threadID string, sortBy string, limit int, offset int) ([]entity.Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	threadID = strings.TrimSpace(threadID)
	if _, exists := s.threadsByID[threadID]; !exists {
		return []entity.Message{}, nil
	}

	ids := s.messageIDsByThread[threadID]
	items := make([]entity.Message, 0, len(ids))
	for _, id := range ids {
		message, ok := s.messagesByID[id]
		if !ok {
			continue
		}
		items = append(items, cloneMessage(message))
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
