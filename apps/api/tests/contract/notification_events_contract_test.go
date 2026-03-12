package contract_test

import (
	"testing"

	notificationevents "github.com/Tokuchi61/Manga/apps/api/internal/modules/notification/events"
	"github.com/stretchr/testify/require"
)

func TestNotificationEventConstants(t *testing.T) {
	require.Equal(t, "notification.created", notificationevents.EventNotificationCreated)
	require.Equal(t, "notification.delivered", notificationevents.EventNotificationDelivered)
	require.Equal(t, "notification.failed", notificationevents.EventNotificationFailed)
	require.Equal(t, "notification.read", notificationevents.EventNotificationRead)
}
