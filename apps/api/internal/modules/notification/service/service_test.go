package service

import (
	"context"
	"testing"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/notification/dto"
	notificationrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/notification/repository"
	supportcontract "github.com/Tokuchi61/Manga/apps/api/internal/modules/support/contract"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type supportSignalProviderStub struct {
	signal supportcontract.NotificationSignal
	err    error
}

func (s supportSignalProviderStub) BuildNotificationSignal(_ context.Context, _ string, _ string, _ string, _ string) (supportcontract.NotificationSignal, error) {
	if s.err != nil {
		return supportcontract.NotificationSignal{}, s.err
	}
	return s.signal, nil
}

func TestCreateFromSupportSignalInboxAndReadFlow(t *testing.T) {
	store := notificationrepository.NewMemoryStore()
	svc := New(store, validation.New())
	now := time.Date(2026, 3, 12, 23, 0, 0, 0, time.UTC)
	svc.now = func() time.Time { return now }

	requesterID := uuid.NewString()
	supportID := uuid.NewString()
	svc.SetSupportSignalProvider(supportSignalProviderStub{signal: supportcontract.NotificationSignal{
		Event:           supportcontract.EventSupportResolved,
		SupportID:       supportID,
		RequesterUserID: requesterID,
		Status:          "resolved",
		OccurredAt:      now,
		RequestID:       "req-stage12-support",
		CorrelationID:   "corr-stage12-support",
	}})

	created, err := svc.CreateFromSupportSignal(context.Background(), dto.CreateFromSupportSignalRequest{SupportID: supportID})
	require.NoError(t, err)
	require.True(t, created.Created)
	firstID := created.NotificationID

	again, err := svc.CreateFromSupportSignal(context.Background(), dto.CreateFromSupportSignalRequest{SupportID: supportID})
	require.NoError(t, err)
	require.False(t, again.Created)
	require.Equal(t, firstID, again.NotificationID)

	inbox, err := svc.ListInbox(context.Background(), dto.ListInboxRequest{UserID: requesterID})
	require.NoError(t, err)
	require.Equal(t, 1, inbox.Count)
	require.Equal(t, 1, inbox.UnreadCount)

	_, err = svc.MarkRead(context.Background(), dto.MarkReadRequest{UserID: requesterID, NotificationID: firstID})
	require.NoError(t, err)

	inboxAfterRead, err := svc.ListInbox(context.Background(), dto.ListInboxRequest{UserID: requesterID})
	require.NoError(t, err)
	require.Equal(t, 0, inboxAfterRead.UnreadCount)
}

func TestUpdatePreferenceAndMuteBlocksSupportSignal(t *testing.T) {
	store := notificationrepository.NewMemoryStore()
	svc := New(store, validation.New())

	requesterID := uuid.NewString()
	supportID := uuid.NewString()
	svc.SetSupportSignalProvider(supportSignalProviderStub{signal: supportcontract.NotificationSignal{
		Event:           supportcontract.EventSupportCreated,
		SupportID:       supportID,
		RequesterUserID: requesterID,
		Status:          "open",
		RequestID:       "req-stage12-pref",
		CorrelationID:   "corr-stage12-pref",
	}})

	_, err := svc.UpdatePreference(context.Background(), dto.UpdatePreferenceRequest{
		UserID:            requesterID,
		MutedCategories:   []string{"support"},
		QuietHoursEnabled: true,
		QuietHoursStart:   23,
		QuietHoursEnd:     7,
		InAppEnabled:      true,
		EmailEnabled:      true,
		PushEnabled:       false,
		DigestEnabled:     true,
	})
	require.NoError(t, err)

	_, err = svc.CreateFromSupportSignal(context.Background(), dto.CreateFromSupportSignalRequest{SupportID: supportID})
	require.ErrorIs(t, err, ErrCategoryMuted)
}

func TestRuntimePauseBlocksSupportSignal(t *testing.T) {
	store := notificationrepository.NewMemoryStore()
	svc := New(store, validation.New())

	requesterID := uuid.NewString()
	supportID := uuid.NewString()
	svc.SetSupportSignalProvider(supportSignalProviderStub{signal: supportcontract.NotificationSignal{
		Event:           supportcontract.EventSupportCreated,
		SupportID:       supportID,
		RequesterUserID: requesterID,
		Status:          "open",
		RequestID:       "req-stage12-runtime",
		CorrelationID:   "corr-stage12-runtime",
	}})

	_, err := svc.UpdateDeliveryPause(context.Background(), dto.UpdateDeliveryPauseRequest{Paused: true})
	require.NoError(t, err)

	_, err = svc.CreateFromSupportSignal(context.Background(), dto.CreateFromSupportSignalRequest{SupportID: supportID})
	require.ErrorIs(t, err, ErrDeliveryPaused)
}
