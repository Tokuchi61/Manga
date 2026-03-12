package chapter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestModuleNameIsCanonical(t *testing.T) {
	module := New()
	require.Equal(t, "chapter", module.Name())
}

func TestModuleMountsChapterRoutes(t *testing.T) {
	module := New()
	router := chi.NewRouter()
	module.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/manga/7b8f6f9f-4d3d-423f-96b2-2fe7180f5940/chapters", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	require.NotEqual(t, http.StatusNotFound, rec.Code)
}
