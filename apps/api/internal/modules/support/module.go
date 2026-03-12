package support

import (
	"context"
	"errors"

	supportcontract "github.com/Tokuchi61/Manga/apps/api/internal/modules/support/contract"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/support/handler"
	supportrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/support/repository"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/support/service"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/snapshot"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
	"github.com/go-chi/chi/v5"
)

// Module wires support handlers to the central module registry.
type Module struct {
	httpHandler *handler.HTTPHandler
	svc         *service.SupportService
	snapshotter snapshot.Snapshotter
}

func New() Module {
	store := supportrepository.NewMemoryStore()
	validator := validation.New()
	svc := service.New(store, validator)

	return Module{httpHandler: handler.New(svc), svc: svc, snapshotter: store}
}

func (m *Module) SetTargetLookups(mangaLookup service.MangaTargetLookup, chapterLookup service.ChapterTargetLookup, commentLookup service.CommentTargetLookup) {
	if m == nil || m.svc == nil {
		return
	}
	m.svc.SetTargetLookups(mangaLookup, chapterLookup, commentLookup)
}

func (m Module) GetModerationHandoffReference(ctx context.Context, supportID string, requestID string, correlationID string) (supportcontract.ModerationHandoffReference, error) {
	if m.svc == nil {
		return supportcontract.ModerationHandoffReference{}, errors.New("support_service_unavailable")
	}
	return m.svc.GetModerationHandoffReference(ctx, supportID, requestID, correlationID)
}

func (m Module) LinkModerationCase(ctx context.Context, supportID string, moderationCaseID string) error {
	if m.svc == nil {
		return errors.New("support_service_unavailable")
	}
	return m.svc.LinkModerationCase(ctx, supportID, moderationCaseID)
}

func (m Module) BuildNotificationSignal(ctx context.Context, supportID string, event string, requestID string, correlationID string) (supportcontract.NotificationSignal, error) {
	if m.svc == nil {
		return supportcontract.NotificationSignal{}, errors.New("support_service_unavailable")
	}
	return m.svc.BuildNotificationSignal(ctx, supportID, event, requestID, correlationID)
}

func (m Module) Name() string {
	return "support"
}

func (m Module) RegisterRoutes(router chi.Router) {
	registerRoutes(router, m.httpHandler)
}

func (m *Module) Snapshot() ([]byte, error) {
	if m == nil || m.snapshotter == nil {
		return nil, nil
	}
	return m.snapshotter.Snapshot()
}

func (m *Module) RestoreSnapshot(data []byte) error {
	if m == nil || m.snapshotter == nil {
		return nil
	}
	return m.snapshotter.RestoreSnapshot(data)
}
