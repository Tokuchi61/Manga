package user

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/user/handler"
	userrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/user/repository"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/user/service"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
	"github.com/go-chi/chi/v5"
)

// Module wires user handlers to the central module registry.
type Module struct {
	httpHandler *handler.HTTPHandler
	svc         *service.UserService
}

func New() Module {
	store := userrepository.NewMemoryStore()
	validator := validation.New()
	svc := service.New(store, validator)

	return Module{httpHandler: handler.New(svc), svc: svc}
}

func (m *Module) SetCredentialLookup(lookup service.CredentialLookup) {
	if m == nil || m.svc == nil {
		return
	}
	m.svc.SetCredentialLookup(lookup)
}

func (m Module) Name() string {
	return "user"
}

func (m Module) RegisterRoutes(router chi.Router) {
	registerRoutes(router, m.httpHandler)
}
