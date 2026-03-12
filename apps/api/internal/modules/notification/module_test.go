package notification

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestModuleNameIsCanonical(t *testing.T) {
	module := New()
	require.Equal(t, "notification", module.Name())
}

func TestModuleMountsNotificationRoutes(t *testing.T) {
	module := New()
	router := chi.NewRouter()
	module.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/notifications/inbox", nil)
	req.Header.Set("X-Actor-User-ID", uuid.NewString())
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	require.NotEqual(t, http.StatusNotFound, rec.Code)
}
