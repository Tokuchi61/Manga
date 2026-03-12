package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	supportcontract "github.com/Tokuchi61/Manga/apps/api/internal/modules/support/contract"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/notification/repository"
	notificationvalidator "github.com/Tokuchi61/Manga/apps/api/internal/modules/notification/validator"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
)

var (
	ErrValidation               = errors.New("notification_validation_failed")
	ErrNotFound                 = errors.New("notification_not_found")
	ErrForbiddenAction          = errors.New("notification_forbidden_action")
	ErrCategoryDisabled         = errors.New("notification_category_disabled")
	ErrCategoryMuted            = errors.New("notification_category_muted")
	ErrChannelDisabled          = errors.New("notification_channel_disabled")
	ErrDeliveryPaused           = errors.New("notification_delivery_paused")
	ErrDigestDisabled           = errors.New("notification_digest_disabled")
	ErrSupportSignalUnavailable = errors.New("notification_support_signal_unavailable")
	ErrSupportSignalInvalid     = errors.New("notification_support_signal_invalid")
)

// SupportSignalProvider exposes support->notification signal build boundary.
type SupportSignalProvider interface {
	BuildNotificationSignal(ctx context.Context, supportID string, event string, requestID string, correlationID string) (supportcontract.NotificationSignal, error)
}

// NotificationService owns stage-12 notification flows.
type NotificationService struct {
	store                 repository.Store
	validator             *validation.Validator
	now                   func() time.Time
	supportSignalProvider SupportSignalProvider
}

func New(store repository.Store, validator *validation.Validator) *NotificationService {
	if store == nil {
		store = repository.NewMemoryStore()
	}
	return &NotificationService{
		store:     store,
		validator: validator,
		now:       time.Now,
	}
}

func (s *NotificationService) SetSupportSignalProvider(provider SupportSignalProvider) {
	s.supportSignalProvider = provider
}

func (s *NotificationService) validateInput(payload any) error {
	if err := notificationvalidator.ValidateStruct(s.validator, payload); err != nil {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}
	return nil
}
