package social

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestModuleNameIsCanonical(t *testing.T) {
	module := New()
	require.Equal(t, "social", module.Name())
}

func TestModuleMountsSocialRoutes(t *testing.T) {
	module := New()
	router := chi.NewRouter()
	module.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/social/friendships", nil)
	token, err := identity.IssueAccessToken(identity.TokenClaims{UserID: uuid.NewString(), ExpiresAt: time.Now().UTC().Add(time.Minute)})
	require.NoError(t, err)
	req.Header.Set(identity.HeaderAuthorization, "Bearer "+token)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	require.NotEqual(t, http.StatusNotFound, rec.Code)
}
