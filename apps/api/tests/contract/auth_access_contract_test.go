package contract_test

import (
	"testing"
	"time"

	authcontract "github.com/Tokuchi61/Manga/apps/api/internal/modules/auth/contract"
	"github.com/stretchr/testify/require"
)

func TestAuthAccessContractShape(t *testing.T) {
	identity := authcontract.VerifiedIdentity{
		CredentialID:    "cred-1",
		SessionID:       "sess-1",
		EmailVerified:   true,
		Suspended:       false,
		Banned:          false,
		AuthenticatedAt: time.Date(2026, 3, 12, 10, 0, 0, 0, time.UTC),
	}
	require.True(t, identity.EmailVerified)
	require.False(t, identity.Suspended)
	require.False(t, identity.Banned)

	signal := authcontract.SecuritySignal{
		CredentialID:  "cred-1",
		Signal:        "auth.security.suspicious_login",
		OccurredAt:    time.Date(2026, 3, 12, 10, 5, 0, 0, time.UTC),
		RequestID:     "req-1",
		CorrelationID: "corr-1",
	}
	require.Equal(t, "auth.security.suspicious_login", signal.Signal)
	require.NotEmpty(t, signal.RequestID)
	require.NotEmpty(t, signal.CorrelationID)
}
