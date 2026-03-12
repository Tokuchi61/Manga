package contract_test

import (
	"testing"
	"time"

	mangacontract "github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/contract"
	"github.com/stretchr/testify/require"
)

func TestMangaChapterDefaultsContractShape(t *testing.T) {
	now := time.Date(2026, 3, 12, 14, 0, 0, 0, time.UTC)
	defaults := mangacontract.ChapterAccessDefaults{
		MangaID:            "e6d2b929-9f69-4783-a8bf-c99751d6f93f",
		ReadAccessLevel:    "authenticated",
		EarlyAccessEnabled: true,
		EarlyAccessLevel:   "vip",
		UpdatedAt:          now,
	}

	require.Equal(t, "authenticated", defaults.ReadAccessLevel)
	require.True(t, defaults.EarlyAccessEnabled)
	require.Equal(t, "vip", defaults.EarlyAccessLevel)
	require.Equal(t, now, defaults.UpdatedAt)
}

func TestMangaTargetReferenceContractShape(t *testing.T) {
	reference := mangacontract.TargetReference{
		TargetType: mangacontract.TargetTypeManga,
		TargetID:   "e6d2b929-9f69-4783-a8bf-c99751d6f93f",
	}
	require.Equal(t, "manga", reference.TargetType)
	require.NotEmpty(t, reference.TargetID)
}
