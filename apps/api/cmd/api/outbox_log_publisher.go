package main

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/platform/outbox"
	"go.uber.org/zap"
)

type outboxLogPublisher struct {
	logger *zap.Logger
}

func (p outboxLogPublisher) Publish(ctx context.Context, event outbox.Event) error {
	if p.logger != nil {
		p.logger.Info("outbox_event_published",
			zap.String("event_id", event.EventID),
			zap.String("topic", event.Topic),
			zap.String("request_id", event.RequestID),
			zap.String("correlation_id", event.CorrelationID),
		)
	}
	return nil
}
