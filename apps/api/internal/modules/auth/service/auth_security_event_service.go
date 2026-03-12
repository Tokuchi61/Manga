package service

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/entity"
	"github.com/google/uuid"
)

func (s *AuthService) appendSecurityEvent(ctx context.Context, credentialID uuid.UUID, action string, result string, reason string, meta RequestMeta) {
	meta = normalizeMeta(meta)
	event := entity.SecurityEvent{
		ID:            uuid.New(),
		CredentialID:  &credentialID,
		ActorID:       &credentialID,
		TargetID:      &credentialID,
		Action:        action,
		Result:        result,
		Reason:        reason,
		RequestID:     meta.RequestID,
		CorrelationID: meta.CorrelationID,
		Device:        meta.Device,
		IP:            meta.IP,
		CreatedAt:     s.now().UTC(),
	}
	_ = s.store.AppendSecurityEvent(ctx, event)
}
