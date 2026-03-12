package repository

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/notification/entity"
)

func (s *MemoryStore) CreateNotification(_ context.Context, notification entity.Notification, dedupKey string) (entity.Notification, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	dedupKey = normalizeValue(dedupKey)
	if dedupKey != "" {
		notificationID, exists := s.dedupIndex[dedupKey]
		if exists {
			existing, ok := s.notificationsByID[notificationID]
			if ok {
				return cloneNotification(existing), false, nil
			}
			delete(s.dedupIndex, dedupKey)
		}
	}

	if _, exists := s.notificationsByID[notification.ID]; exists {
		return entity.Notification{}, false, ErrConflict
	}

	stored := cloneNotification(notification)
	s.notificationsByID[stored.ID] = stored
	if dedupKey != "" {
		s.dedupIndex[dedupKey] = stored.ID
	}

	return cloneNotification(stored), true, nil
}

func (s *MemoryStore) GetNotificationByID(_ context.Context, notificationID string) (entity.Notification, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	notification, ok := s.notificationsByID[notificationID]
	if !ok {
		return entity.Notification{}, ErrNotFound
	}
	return cloneNotification(notification), nil
}

func (s *MemoryStore) ListInbox(_ context.Context, query InboxQuery) ([]entity.Notification, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	userIDKey := normalizeValue(query.UserID)
	categoryKey := normalizeValue(query.Category)
	stateKey := normalizeValue(query.State)

	items := make([]entity.Notification, 0)
	for _, notification := range s.notificationsByID {
		if normalizeValue(notification.UserID) != userIDKey {
			continue
		}
		if categoryKey != "" && normalizeValue(string(notification.Category)) != categoryKey {
			continue
		}
		if stateKey != "" && normalizeValue(string(notification.State)) != stateKey {
			continue
		}
		if query.UnreadOnly && notification.State == entity.NotificationStateRead {
			continue
		}
		items = append(items, cloneNotification(notification))
	}

	sortNotifications(items, query.SortBy)
	return applyOffsetLimit(items, query.Offset, query.Limit), nil
}

func (s *MemoryStore) UpdateNotification(_ context.Context, notification entity.Notification) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.notificationsByID[notification.ID]; !exists {
		return ErrNotFound
	}
	s.notificationsByID[notification.ID] = cloneNotification(notification)
	return nil
}
