package entity

import (
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/shared/catalog"
)

// DeliveryChannel defines notification delivery lane.
type DeliveryChannel string

const (
	DeliveryChannelInApp DeliveryChannel = "in_app"
	DeliveryChannelEmail DeliveryChannel = "email"
	DeliveryChannelPush  DeliveryChannel = "push"
)

// NotificationState defines notification lifecycle state.
type NotificationState string

const (
	NotificationStateCreated   NotificationState = "created"
	NotificationStateDelivered NotificationState = "delivered"
	NotificationStateFailed    NotificationState = "failed"
	NotificationStateRead      NotificationState = "read"
)

// Notification is notification owner aggregate for stage-12 flows.
type Notification struct {
	ID                   string
	UserID               string
	Category             catalog.NotificationCategory
	Channel              DeliveryChannel
	TemplateKey          string
	Title                string
	Body                 string
	State                NotificationState
	SourceEvent          string
	SourceRefID          string
	RequestID            string
	CorrelationID        string
	DeliveryAttemptCount int
	LastFailureReason    string
	ReadAt               *time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// Preference is user-owned notification preference state.
type Preference struct {
	UserID            string
	MutedCategories   []catalog.NotificationCategory
	QuietHoursEnabled bool
	QuietHoursStart   int
	QuietHoursEnd     int
	InAppEnabled      bool
	EmailEnabled      bool
	PushEnabled       bool
	DigestEnabled     bool
	UpdatedAt         time.Time
}

// RuntimeConfig stores module-level operational controls.
type RuntimeConfig struct {
	CategoryEnabled map[catalog.NotificationCategory]bool
	ChannelEnabled  map[DeliveryChannel]bool
	DigestEnabled   bool
	DeliveryPaused  bool
	UpdatedAt       time.Time
}
