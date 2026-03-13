package ads

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestModuleNameIsCanonical(t *testing.T) {
	module := New()
	require.Equal(t, "ads", module.Name())
}

func TestModuleMountsAdsRoutes(t *testing.T) {
	module := New()
	router := chi.NewRouter()
	module.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodPost, "/ads/impressions", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}
