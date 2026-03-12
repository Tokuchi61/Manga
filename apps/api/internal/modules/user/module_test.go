package user

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestModuleNameIsCanonical(t *testing.T) {
	module := New()
	require.Equal(t, "user", module.Name())
}

func TestModuleMountsUserRoutes(t *testing.T) {
	module := New()
	router := chi.NewRouter()
	module.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"invalid":true}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.NotEqual(t, http.StatusNotFound, rec.Code)
}
