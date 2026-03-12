package service

import (
	"context"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/social/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/social/repository"
)

func (s *SocialService) UpdateBlock(ctx context.Context, request dto.UpdateBlockRequest) (dto.RelationUpdateResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.RelationUpdateResponse{}, err
	}

	actorUserID, err := parseID(request.ActorUserID, "actor_user_id")
	if err != nil {
		return dto.RelationUpdateResponse{}, err
	}
	targetUserID, err := parseID(request.TargetUserID, "target_user_id")
	if err != nil {
		return dto.RelationUpdateResponse{}, err
	}
	if err := ensureDistinctActorTarget(actorUserID, targetUserID); err != nil {
		return dto.RelationUpdateResponse{}, err
	}

	now := s.now().UTC()
	if err := s.store.SetBlock(ctx, actorUserID, targetUserID, request.Enabled, now); err != nil {
		return dto.RelationUpdateResponse{}, err
	}

	if request.Enabled {
		_, _ = s.store.DeleteFriendship(ctx, actorUserID, targetUserID)
		_ = s.store.DeleteFriendRequest(ctx, actorUserID, targetUserID)
		_ = s.store.DeleteFriendRequest(ctx, targetUserID, actorUserID)
		_, _ = s.store.DeleteFollow(ctx, actorUserID, targetUserID)
		_, _ = s.store.DeleteFollow(ctx, targetUserID, actorUserID)
	}

	return mapRelation(repository.RelationEntry{
		ActorUserID:  actorUserID,
		TargetUserID: targetUserID,
		UpdatedAt:    now,
	}, "blocked", request.Enabled), nil
}

func (s *SocialService) UpdateMute(ctx context.Context, request dto.UpdateMuteRequest) (dto.RelationUpdateResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.RelationUpdateResponse{}, err
	}

	actorUserID, err := parseID(request.ActorUserID, "actor_user_id")
	if err != nil {
		return dto.RelationUpdateResponse{}, err
	}
	targetUserID, err := parseID(request.TargetUserID, "target_user_id")
	if err != nil {
		return dto.RelationUpdateResponse{}, err
	}
	if err := ensureDistinctActorTarget(actorUserID, targetUserID); err != nil {
		return dto.RelationUpdateResponse{}, err
	}

	now := s.now().UTC()
	if err := s.store.SetMute(ctx, actorUserID, targetUserID, request.Enabled, now); err != nil {
		return dto.RelationUpdateResponse{}, err
	}

	return mapRelation(repository.RelationEntry{
		ActorUserID:  actorUserID,
		TargetUserID: targetUserID,
		UpdatedAt:    now,
	}, "muted", request.Enabled), nil
}

func (s *SocialService) UpdateRestrict(ctx context.Context, request dto.UpdateRestrictRequest) (dto.RelationUpdateResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.RelationUpdateResponse{}, err
	}

	actorUserID, err := parseID(request.ActorUserID, "actor_user_id")
	if err != nil {
		return dto.RelationUpdateResponse{}, err
	}
	targetUserID, err := parseID(request.TargetUserID, "target_user_id")
	if err != nil {
		return dto.RelationUpdateResponse{}, err
	}
	if err := ensureDistinctActorTarget(actorUserID, targetUserID); err != nil {
		return dto.RelationUpdateResponse{}, err
	}

	now := s.now().UTC()
	if err := s.store.SetRestrict(ctx, actorUserID, targetUserID, request.Enabled, now); err != nil {
		return dto.RelationUpdateResponse{}, err
	}

	return mapRelation(repository.RelationEntry{
		ActorUserID:  actorUserID,
		TargetUserID: targetUserID,
		UpdatedAt:    now,
	}, "restricted", request.Enabled), nil
}

func (s *SocialService) ListRelations(ctx context.Context, request dto.ListRelationsRequest) (dto.ListRelationsResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ListRelationsResponse{}, err
	}

	actorUserID, err := parseID(request.ActorUserID, "actor_user_id")
	if err != nil {
		return dto.ListRelationsResponse{}, err
	}

	relationType := normalizeValue(request.RelationType)
	var items []repository.RelationEntry
	switch relationType {
	case "blocked":
		items, err = s.store.ListBlocked(ctx, actorUserID)
	case "muted":
		items, err = s.store.ListMuted(ctx, actorUserID)
	case "restricted":
		items, err = s.store.ListRestricted(ctx, actorUserID)
	default:
		return dto.ListRelationsResponse{}, ErrValidation
	}
	if err != nil {
		return dto.ListRelationsResponse{}, err
	}

	result := make([]dto.RelationItemResponse, 0, len(items))
	for _, item := range items {
		result = append(result, mapRelationItem(item))
	}

	return dto.ListRelationsResponse{Items: result, Count: len(result)}, nil
}
