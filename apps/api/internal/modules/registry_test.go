package modules

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

type testModule struct {
	name       string
	wasMounted bool
}

func (m *testModule) Name() string {
	return m.name
}

func (m *testModule) RegisterRoutes(router chi.Router) {
	m.wasMounted = true
	router.Get("/"+m.name, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
}

func TestNewRegistryRejectsNilModule(t *testing.T) {
	_, err := NewRegistry(nil)
	require.Error(t, err)
}

func TestNewRegistryRejectsInvalidName(t *testing.T) {
	_, err := NewRegistry(&testModule{name: "Auth"})
	require.Error(t, err)
}

func TestNewRegistryRejectsDuplicateNames(t *testing.T) {
	_, err := NewRegistry(
		&testModule{name: "auth"},
		&testModule{name: "auth"},
	)
	require.Error(t, err)
}

func TestRegistryNamesSorted(t *testing.T) {
	registry, err := NewRegistry(
		&testModule{name: "user"},
		&testModule{name: "auth"},
	)
	require.NoError(t, err)

	require.Equal(t, []string{"auth", "user"}, registry.Names())
}

func TestRegistryMountRoutes(t *testing.T) {
	authModule := &testModule{name: "auth"}
	registry, err := NewRegistry(authModule)
	require.NoError(t, err)

	router := chi.NewRouter()
	registry.MountRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/auth", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.True(t, authModule.wasMounted)
	require.Equal(t, http.StatusNoContent, rec.Code)
}
