package password

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashAndVerify(t *testing.T) {
	hash, err := Hash("StrongPass123!")
	require.NoError(t, err)
	require.NotEmpty(t, hash)

	ok, err := Verify(hash, "StrongPass123!")
	require.NoError(t, err)
	require.True(t, ok)

	ok, err = Verify(hash, "WrongPass")
	require.NoError(t, err)
	require.False(t, ok)
}
