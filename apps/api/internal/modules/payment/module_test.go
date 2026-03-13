package payment

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
	module := New(RuntimeConfig{})
	require.Equal(t, "payment", module.Name())
}

func TestModuleMountsPaymentRoutes(t *testing.T) {
	module := New(RuntimeConfig{})
	router := chi.NewRouter()
	module.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodPost, "/payment/checkout/sessions", nil)
	token, err := identity.IssueAccessToken(identity.TokenClaims{UserID: uuid.NewString(), ExpiresAt: time.Now().UTC().Add(time.Minute)})
	require.NoError(t, err)
	req.Header.Set(identity.HeaderAuthorization, "Bearer "+token)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}
