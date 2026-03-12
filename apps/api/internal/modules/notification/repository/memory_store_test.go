package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/notification/entity"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/catalog"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestMemoryStoreCreateNotificationDedupAndInbox(t *testing.T) {
	store := NewMemoryStore()
	now := time.Date(2026, 3, 12, 22, 0, 0, 0, time.UTC)
	userID := uuid.NewString()

	notification := entity.Notification{
		ID:        uuid.NewString(),
		UserID:    userID,
		Category:  catalog.NotificationSupport,
		Channel:   entity.DeliveryChannelInApp,
		Title:     "Support",
		Body:      "Case updated",
		State:     entity.NotificationStateDelivered,
		CreatedAt: now,
		UpdatedAt: now,
	}

	created, inserted, err := store.CreateNotification(context.Background(), notification, "support:event:1")
	require.NoError(t, err)
	require.True(t, inserted)
	require.Equal(t, notification.ID, created.ID)

	again, inserted, err := store.CreateNotification(context.Background(), notification, "support:event:1")
	require.NoError(t, err)
	require.False(t, inserted)
	require.Equal(t, notification.ID, again.ID)

	items, err := store.ListInbox(context.Background(), InboxQuery{UserID: userID, Limit: 10})
	require.NoError(t, err)
	require.Len(t, items, 1)
	require.Equal(t, notification.ID, items[0].ID)
}

func TestMemoryStorePreferenceAndRuntimeConfig(t *testing.T) {
	store := NewMemoryStore()
	userID := uuid.NewString()

	preference, err := store.GetPreferenceByUserID(context.Background(), userID)
	require.NoError(t, err)
	require.True(t, preference.InAppEnabled)

	preference.UserID = userID
	preference.MutedCategories = []catalog.NotificationCategory{catalog.NotificationSupport}
	preference.QuietHoursEnabled = true
	preference.QuietHoursStart = 23
	preference.QuietHoursEnd = 7
	preference.InAppEnabled = false
	preference.EmailEnabled = true
	preference.PushEnabled = false
	preference.DigestEnabled = true
	preference.UpdatedAt = time.Now().UTC()

	updatedPreference, err := store.UpsertPreference(context.Background(), preference)
	require.NoError(t, err)
	require.False(t, updatedPreference.InAppEnabled)
	require.Len(t, updatedPreference.MutedCategories, 1)

	runtimeConfig, err := store.GetRuntimeConfig(context.Background())
	require.NoError(t, err)
	runtimeConfig.DeliveryPaused = true
	runtimeConfig.CategoryEnabled[catalog.NotificationSupport] = false
	runtimeConfig.ChannelEnabled[entity.DeliveryChannelEmail] = false
	runtimeConfig.UpdatedAt = time.Now().UTC()
	require.NoError(t, store.UpdateRuntimeConfig(context.Background(), runtimeConfig))

	persistedRuntimeConfig, err := store.GetRuntimeConfig(context.Background())
	require.NoError(t, err)
	require.True(t, persistedRuntimeConfig.DeliveryPaused)
	require.False(t, persistedRuntimeConfig.CategoryEnabled[catalog.NotificationSupport])
	require.False(t, persistedRuntimeConfig.ChannelEnabled[entity.DeliveryChannelEmail])
}
