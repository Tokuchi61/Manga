package main

import (
	"context"
	"encoding/json"
	"strings"

	paymentservice "github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/service"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/outbox"
	"github.com/google/uuid"
)

type paymentCallbackOutboxPublisher struct {
	store outbox.Store
}

func (p paymentCallbackOutboxPublisher) PublishPaymentCallback(ctx context.Context, event paymentservice.PaymentCallbackEvent) error {
	if p.store == nil {
		return nil
	}
	payload, err := json.Marshal(map[string]any{
		"provider_event_id":  event.ProviderEventID,
		"session_id":         event.SessionID,
		"transaction_id":     event.TransactionID,
		"transaction_status": event.TransactionStatus,
		"processed_at":       event.ProcessedAt.UTC(),
	})
	if err != nil {
		return err
	}

	requestID := strings.TrimSpace(event.RequestID)
	if requestID == "" {
		requestID = event.ProviderEventID
	}
	return p.store.Enqueue(ctx, outbox.Event{
		EventID:       uuid.NewString(),
		Topic:         "payment.callback.processed",
		Payload:       payload,
		RequestID:     requestID,
		CorrelationID: strings.TrimSpace(event.CorrelationID),
		CausationID:   strings.TrimSpace(event.ProviderEventID),
	})
}
