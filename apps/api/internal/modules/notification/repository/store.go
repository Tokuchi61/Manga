package repository

import (
	"context"
	"errors"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/notification/entity"
)

var (
	ErrNotFound = errors.New("notification_repository_not_found")
	ErrConflict = errors.New("notification_repository_conflict")
)

// InboxQuery defines requester-owned inbox listing controls.
type InboxQuery struct {
	UserID     string
	Category   string
	State      string
	UnreadOnly bool
	SortBy     string
	Limit      int
	Offset     int
}

// Store defines notification persistence boundary.
type Store interface {
	CreateNotification(ctx context.Context, notification entity.Notification, dedupKey string) (entity.Notification, bool, error)
	GetNotificationByID(ctx context.Context, notificationID string) (entity.Notification, error)
	ListInbox(ctx context.Context, query InboxQuery) ([]entity.Notification, error)
	UpdateNotification(ctx context.Context, notification entity.Notification) error

	GetPreferenceByUserID(ctx context.Context, userID string) (entity.Preference, error)
	UpsertPreference(ctx context.Context, preference entity.Preference) (entity.Preference, error)

	GetRuntimeConfig(ctx context.Context) (entity.RuntimeConfig, error)
	UpdateRuntimeConfig(ctx context.Context, cfg entity.RuntimeConfig) error
}
