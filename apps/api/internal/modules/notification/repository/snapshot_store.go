package repository

import (
	"bytes"
	"encoding/gob"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/notification/entity"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/catalog"
)

type notificationSnapshotState struct {
	NotificationsByID map[string]entity.Notification
	DedupIndex        map[string]string
	PreferencesByUser map[string]entity.Preference
	RuntimeConfig     entity.RuntimeConfig
}

func (s *MemoryStore) Snapshot() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state := notificationSnapshotState{
		NotificationsByID: s.notificationsByID,
		DedupIndex:        s.dedupIndex,
		PreferencesByUser: s.preferencesByUser,
		RuntimeConfig:     s.runtimeConfig,
	}

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(state); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *MemoryStore) RestoreSnapshot(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	var state notificationSnapshotState
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&state); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if state.NotificationsByID == nil {
		state.NotificationsByID = make(map[string]entity.Notification)
	}
	if state.DedupIndex == nil {
		state.DedupIndex = make(map[string]string)
	}
	if state.PreferencesByUser == nil {
		state.PreferencesByUser = make(map[string]entity.Preference)
	}
	if state.RuntimeConfig.CategoryEnabled == nil {
		state.RuntimeConfig.CategoryEnabled = make(map[catalog.NotificationCategory]bool)
		for _, category := range catalog.AllNotificationCategories {
			state.RuntimeConfig.CategoryEnabled[category] = true
		}
	}
	if state.RuntimeConfig.ChannelEnabled == nil {
		state.RuntimeConfig.ChannelEnabled = map[entity.DeliveryChannel]bool{
			entity.DeliveryChannelInApp: true,
			entity.DeliveryChannelEmail: true,
			entity.DeliveryChannelPush:  true,
		}
	}

	s.notificationsByID = state.NotificationsByID
	s.dedupIndex = state.DedupIndex
	s.preferencesByUser = state.PreferencesByUser
	s.runtimeConfig = state.RuntimeConfig

	return nil
}
