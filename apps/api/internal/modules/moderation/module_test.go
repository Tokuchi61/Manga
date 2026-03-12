package moderation

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestModuleNameIsCanonical(t *testing.T) {
	module := New()
	if module.Name() != "moderation" {
		t.Fatalf("expected moderation module name, got %s", module.Name())
	}
}

func TestModuleMountsModerationRoutes(t *testing.T) {
	module := New()
	router := chi.NewRouter()
	module.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/moderation/queue", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	if rec.Code == http.StatusNotFound {
		t.Fatalf("expected moderation routes to be mounted")
	}
}
