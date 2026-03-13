package contract_test

import (
	"testing"

	historyevents "github.com/Tokuchi61/Manga/apps/api/internal/modules/history/events"
	"github.com/stretchr/testify/require"
)

func TestHistoryEventConstants(t *testing.T) {
	require.Equal(t, "history.checkpoint.updated", historyevents.EventHistoryCheckpointUpdated)
	require.Equal(t, "history.chapter.finished", historyevents.EventHistoryChapterFinished)
	require.Equal(t, "history.library.changed", historyevents.EventHistoryLibraryChanged)
}
