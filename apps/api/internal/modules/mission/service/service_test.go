package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/mission/dto"
	missionrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/mission/repository"
	"github.com/Tokuchi61/Manga/apps/api/internal/platform/validation"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestMissionServiceProgressClaimAndListFlow(t *testing.T) {
	store := missionrepository.NewMemoryStore()
	svc := New(store, validation.New())
	now := time.Date(2026, 3, 16, 2, 0, 0, 0, time.UTC)
	svc.now = func() time.Time { return now }

	actorID := uuid.NewString()

	_, err := svc.UpsertMissionDefinition(context.Background(), dto.UpsertMissionDefinitionRequest{
		MissionID:      "daily_read_3",
		Category:       "daily",
		Title:          "Read 3 chapters",
		ObjectiveType:  "chapter_read",
		TargetCount:    3,
		RewardItemID:   "mana_potion",
		RewardQuantity: 1,
		Active:         true,
	})
	require.NoError(t, err)

	progressRes, err := svc.IngestMissionProgress(context.Background(), dto.IngestMissionProgressRequest{
		ActorUserID: actorID,
		MissionID:   "daily_read_3",
		Delta:       2,
		SourceType:  "history",
		RequestID:   "req-mission-progress-1",
	})
	require.NoError(t, err)
	require.False(t, progressRes.Completed)
	require.Equal(t, 2, progressRes.Mission.ProgressCount)

	progressRes, err = svc.IngestMissionProgress(context.Background(), dto.IngestMissionProgressRequest{
		ActorUserID: actorID,
		MissionID:   "daily_read_3",
		Delta:       1,
		SourceType:  "history",
		RequestID:   "req-mission-progress-2",
	})
	require.NoError(t, err)
	require.True(t, progressRes.Completed)
	require.Equal(t, 3, progressRes.Mission.ProgressCount)

	claimRes, err := svc.ClaimMission(context.Background(), dto.ClaimMissionRequest{
		ActorUserID: actorID,
		MissionID:   "daily_read_3",
		RequestID:   "req-mission-claim-1",
	})
	require.NoError(t, err)
	require.Equal(t, "claim_requested", claimRes.Status)
	require.Equal(t, 1, claimRes.RewardQuantity)
	require.Equal(t, "claimed", claimRes.Mission.Status)

	_, err = svc.ClaimMission(context.Background(), dto.ClaimMissionRequest{
		ActorUserID: actorID,
		MissionID:   "daily_read_3",
		RequestID:   "req-mission-claim-2",
	})
	require.True(t, errors.Is(err, ErrAlreadyClaimed))

	listRes, err := svc.ListActorMissions(context.Background(), dto.ListActorMissionsRequest{ActorUserID: actorID})
	require.NoError(t, err)
	require.Equal(t, 1, listRes.Count)
	require.Equal(t, "claimed", listRes.Items[0].Status)
}

func TestMissionServiceRuntimeTogglesAffectOperations(t *testing.T) {
	store := missionrepository.NewMemoryStore()
	svc := New(store, validation.New())
	actorID := uuid.NewString()

	_, err := svc.UpsertMissionDefinition(context.Background(), dto.UpsertMissionDefinitionRequest{
		MissionID:      "daily_comment_2",
		Category:       "daily",
		Title:          "Write 2 comments",
		ObjectiveType:  "comment_create",
		TargetCount:    2,
		RewardQuantity: 1,
		Active:         true,
	})
	require.NoError(t, err)

	_, err = svc.UpdateProgressIngestState(context.Background(), dto.UpdateProgressIngestStateRequest{Enabled: false})
	require.NoError(t, err)
	_, err = svc.IngestMissionProgress(context.Background(), dto.IngestMissionProgressRequest{ActorUserID: actorID, MissionID: "daily_comment_2", Delta: 1, SourceType: "comment", RequestID: "req-disabled-ingest"})
	require.True(t, errors.Is(err, ErrProgressIngestDisabled))

	_, err = svc.UpdateProgressIngestState(context.Background(), dto.UpdateProgressIngestStateRequest{Enabled: true})
	require.NoError(t, err)
	_, err = svc.IngestMissionProgress(context.Background(), dto.IngestMissionProgressRequest{ActorUserID: actorID, MissionID: "daily_comment_2", Delta: 2, SourceType: "comment", RequestID: "req-complete"})
	require.NoError(t, err)

	_, err = svc.UpdateClaimState(context.Background(), dto.UpdateClaimStateRequest{Enabled: false})
	require.NoError(t, err)
	_, err = svc.ClaimMission(context.Background(), dto.ClaimMissionRequest{ActorUserID: actorID, MissionID: "daily_comment_2", RequestID: "req-disabled-claim"})
	require.True(t, errors.Is(err, ErrClaimDisabled))

	_, err = svc.UpdateReadState(context.Background(), dto.UpdateReadStateRequest{Enabled: false})
	require.NoError(t, err)
	_, err = svc.ListActorMissions(context.Background(), dto.ListActorMissionsRequest{ActorUserID: actorID})
	require.True(t, errors.Is(err, ErrReadDisabled))
}

func TestMissionServiceResetMissionProgress(t *testing.T) {
	store := missionrepository.NewMemoryStore()
	svc := New(store, validation.New())
	actorID := uuid.NewString()

	_, err := svc.UpsertMissionDefinition(context.Background(), dto.UpsertMissionDefinitionRequest{
		MissionID:      "weekly_social_4",
		Category:       "weekly",
		Title:          "Send 4 social messages",
		ObjectiveType:  "social_message",
		TargetCount:    4,
		RewardQuantity: 1,
		Active:         true,
	})
	require.NoError(t, err)

	_, err = svc.IngestMissionProgress(context.Background(), dto.IngestMissionProgressRequest{
		ActorUserID: actorID,
		MissionID:   "weekly_social_4",
		Delta:       3,
		SourceType:  "social",
		RequestID:   "req-weekly-social-progress",
	})
	require.NoError(t, err)

	resetRes, err := svc.ResetMissionProgress(context.Background(), dto.ResetMissionProgressRequest{
		TargetUserID: actorID,
		MissionID:    "weekly_social_4",
	})
	require.NoError(t, err)
	require.Equal(t, "reset", resetRes.Status)
	require.Equal(t, 1, resetRes.DeletedCount)

	detailRes, err := svc.GetActorMissionDetail(context.Background(), dto.GetActorMissionDetailRequest{
		ActorUserID: actorID,
		MissionID:   "weekly_social_4",
	})
	require.NoError(t, err)
	require.Equal(t, 0, detailRes.ProgressCount)
}
