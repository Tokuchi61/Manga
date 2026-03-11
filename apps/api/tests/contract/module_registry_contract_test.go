package contract_test

import (
	"testing"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

type contractModule struct {
	name string
}

func (m contractModule) Name() string {
	return m.name
}

func (m contractModule) RegisterRoutes(router chi.Router) {
	// Contract test only verifies module registry constraints.
}

func TestModuleRegistryContractAcceptsCanonicalNames(t *testing.T) {
	registry, err := modules.NewRegistry(
		contractModule{name: "auth"},
		contractModule{name: "user"},
	)
	require.NoError(t, err)
	require.Equal(t, []string{"auth", "user"}, registry.Names())
}

func TestModuleRegistryContractRejectsNonCanonicalNames(t *testing.T) {
	_, err := modules.NewRegistry(contractModule{name: "auth-service"})
	require.Error(t, err)
}
