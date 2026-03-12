package contract_test

import (
	"testing"
	"time"

	commentcontract "github.com/Tokuchi61/Manga/apps/api/internal/modules/comment/contract"
	"github.com/stretchr/testify/require"
)

func TestCommentTargetRelationContractShape(t *testing.T) {
	now := time.Date(2026, 3, 12, 21, 30, 0, 0, time.UTC)
	parentID := "d222f8ef-16ea-492a-979c-96da5f9f7f40"
	relation := commentcontract.TargetRelation{
		CommentID:        "2d2d1e5e-3578-44f2-a57d-99e402d6b7d4",
		TargetType:       "manga",
		TargetID:         "40dc7d79-9cc5-4f96-a26f-6f73dca9f4ea",
		ParentCommentID:  &parentID,
		RootCommentID:    &parentID,
		ModerationStatus: "visible",
		Deleted:          false,
		UpdatedAt:        now,
	}

	require.Equal(t, "manga", relation.TargetType)
	require.NotNil(t, relation.ParentCommentID)
	require.Equal(t, now, relation.UpdatedAt)
}

func TestCommentModerationSignalContractShape(t *testing.T) {
	now := time.Date(2026, 3, 12, 21, 31, 0, 0, time.UTC)
	signal := commentcontract.ModerationSignal{
		Event:         commentcontract.EventCommentModerated,
		CommentID:     "2d2d1e5e-3578-44f2-a57d-99e402d6b7d4",
		TargetType:    "chapter",
		TargetID:      "40dc7d79-9cc5-4f96-a26f-6f73dca9f4ea",
		OccurredAt:    now,
		RequestID:     "req-stage9",
		CorrelationID: "corr-stage9",
	}

	require.Equal(t, commentcontract.EventCommentModerated, signal.Event)
	require.Equal(t, "req-stage9", signal.RequestID)
	require.Equal(t, "corr-stage9", signal.CorrelationID)
	require.Equal(t, now, signal.OccurredAt)
}
