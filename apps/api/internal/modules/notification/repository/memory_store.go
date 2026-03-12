package repository

import (
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/notification/entity"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/catalog"
)

// MemoryStore is stage-12 bootstrap persistence for notification flows.
type MemoryStore struct {
	mu sync.RWMutex

	notificationsByID map[string]entity.Notification
	dedupIndex        map[string]string
	preferencesByUser map[string]entity.Preference
	runtimeConfig     entity.RuntimeConfig
}

func NewMemoryStore() *MemoryStore {
	now := time.Now().UTC()
	categoryEnabled := make(map[catalog.NotificationCategory]bool, len(catalog.AllNotificationCategories))
	for _, category := range catalog.AllNotificationCategories {
		categoryEnabled[category] = true
	}

	channelEnabled := map[entity.DeliveryChannel]bool{
		entity.DeliveryChannelInApp: true,
		entity.DeliveryChannelEmail: true,
		entity.DeliveryChannelPush:  true,
	}

	return &MemoryStore{
		notificationsByID: make(map[string]entity.Notification),
		dedupIndex:        make(map[string]string),
		preferencesByUser: make(map[string]entity.Preference),
		runtimeConfig: entity.RuntimeConfig{
			CategoryEnabled: categoryEnabled,
			ChannelEnabled:  channelEnabled,
			DigestEnabled:   true,
			DeliveryPaused:  false,
			UpdatedAt:       now,
		},
	}
}

func normalizeValue(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

func cloneTimePtr(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	v := value.UTC()
	return &v
}

func cloneMutedCategories(values []catalog.NotificationCategory) []catalog.NotificationCategory {
	if len(values) == 0 {
		return nil
	}
	out := make([]catalog.NotificationCategory, len(values))
	copy(out, values)
	return out
}

func cloneNotification(in entity.Notification) entity.Notification {
	out := in
	out.ReadAt = cloneTimePtr(in.ReadAt)
	return out
}

func clonePreference(in entity.Preference) entity.Preference {
	out := in
	out.MutedCategories = cloneMutedCategories(in.MutedCategories)
	return out
}

func cloneCategoryEnabled(in map[catalog.NotificationCategory]bool) map[catalog.NotificationCategory]bool {
	if len(in) == 0 {
		return map[catalog.NotificationCategory]bool{}
	}
	out := make(map[catalog.NotificationCategory]bool, len(in))
	for key, value := range in {
		out[key] = value
	}
	return out
}

func cloneChannelEnabled(in map[entity.DeliveryChannel]bool) map[entity.DeliveryChannel]bool {
	if len(in) == 0 {
		return map[entity.DeliveryChannel]bool{}
	}
	out := make(map[entity.DeliveryChannel]bool, len(in))
	for key, value := range in {
		out[key] = value
	}
	return out
}

func cloneRuntimeConfig(in entity.RuntimeConfig) entity.RuntimeConfig {
	out := in
	out.CategoryEnabled = cloneCategoryEnabled(in.CategoryEnabled)
	out.ChannelEnabled = cloneChannelEnabled(in.ChannelEnabled)
	return out
}

func defaultPreference(userID string, now time.Time) entity.Preference {
	return entity.Preference{
		UserID:            userID,
		MutedCategories:   nil,
		QuietHoursEnabled: false,
		QuietHoursStart:   22,
		QuietHoursEnd:     7,
		InAppEnabled:      true,
		EmailEnabled:      true,
		PushEnabled:       true,
		DigestEnabled:     true,
		UpdatedAt:         now,
	}
}

func sortNotifications(items []entity.Notification, sortBy string) {
	switch normalizeValue(sortBy) {
	case "oldest":
		sort.Slice(items, func(i, j int) bool {
			if items[i].CreatedAt.Equal(items[j].CreatedAt) {
				return items[i].ID < items[j].ID
			}
			return items[i].CreatedAt.Before(items[j].CreatedAt)
		})
	default:
		sort.Slice(items, func(i, j int) bool {
			if items[i].CreatedAt.Equal(items[j].CreatedAt) {
				return items[i].ID < items[j].ID
			}
			return items[i].CreatedAt.After(items[j].CreatedAt)
		})
	}
}

func applyOffsetLimit(items []entity.Notification, offset int, limit int) []entity.Notification {
	if offset < 0 {
		offset = 0
	}
	if offset >= len(items) {
		return []entity.Notification{}
	}
	if limit <= 0 {
		limit = 50
	}
	end := offset + limit
	if end > len(items) {
		end = len(items)
	}
	return append([]entity.Notification(nil), items[offset:end]...)
}
