package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/social/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/social/entity"
	socialrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/social/repository"
)

func (s *SocialService) FollowUser(ctx context.Context, request dto.FollowUserRequest) (dto.FollowUserResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.FollowUserResponse{}, err
	}

	actorUserID, err := parseID(request.ActorUserID, "actor_user_id")
	if err != nil {
		return dto.FollowUserResponse{}, err
	}
	targetUserID, err := parseID(request.TargetUserID, "target_user_id")
	if err != nil {
		return dto.FollowUserResponse{}, err
	}
	if err := ensureDistinctActorTarget(actorUserID, targetUserID); err != nil {
		return dto.FollowUserResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.FollowUserResponse{}, err
	}
	if err := s.requireFollowWriteEnabled(cfg.FollowEnabled); err != nil {
		return dto.FollowUserResponse{}, err
	}

	blocked, err := s.store.IsBlockedEither(ctx, actorUserID, targetUserID)
	if err != nil {
		return dto.FollowUserResponse{}, err
	}
	if blocked {
		return dto.FollowUserResponse{}, ErrForbiddenAction
	}

	restricted, err := s.store.IsRestrictedEither(ctx, actorUserID, targetUserID)
	if err != nil {
		return dto.FollowUserResponse{}, err
	}
	if restricted {
		return dto.FollowUserResponse{}, ErrForbiddenAction
	}

	now := s.now().UTC()
	stored, created, err := s.store.UpsertFollow(ctx, entity.FollowRelation{
		FollowerUserID: actorUserID,
		FolloweeUserID: targetUserID,
		RequestID:      strings.TrimSpace(request.RequestID),
		CreatedAt:      now,
		UpdatedAt:      now,
	}, buildDedupKey("follow", actorUserID, targetUserID, request.RequestID))
	if err != nil {
		if errors.Is(err, socialrepository.ErrConflict) {
			return dto.FollowUserResponse{}, ErrConflict
		}
		return dto.FollowUserResponse{}, err
	}

	return dto.FollowUserResponse{
		FollowerUserID: stored.FollowerUserID,
		FolloweeUserID: stored.FolloweeUserID,
		Following:      true,
		Created:        created,
		UpdatedAt:      stored.UpdatedAt,
	}, nil
}

func (s *SocialService) UnfollowUser(ctx context.Context, request dto.UnfollowUserRequest) (dto.OperationResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.OperationResponse{}, err
	}

	actorUserID, err := parseID(request.ActorUserID, "actor_user_id")
	if err != nil {
		return dto.OperationResponse{}, err
	}
	targetUserID, err := parseID(request.TargetUserID, "target_user_id")
	if err != nil {
		return dto.OperationResponse{}, err
	}
	if err := ensureDistinctActorTarget(actorUserID, targetUserID); err != nil {
		return dto.OperationResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.OperationResponse{}, err
	}
	if err := s.requireFollowWriteEnabled(cfg.FollowEnabled); err != nil {
		return dto.OperationResponse{}, err
	}

	if _, err := s.store.DeleteFollow(ctx, actorUserID, targetUserID); err != nil {
		return dto.OperationResponse{}, err
	}
	return dto.OperationResponse{Status: "removed"}, nil
}

func (s *SocialService) ListFollowers(ctx context.Context, request dto.ListFollowersRequest) (dto.ListFollowsResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ListFollowsResponse{}, err
	}

	actorUserID, err := parseID(request.ActorUserID, "actor_user_id")
	if err != nil {
		return dto.ListFollowsResponse{}, err
	}

	items, err := s.store.ListFollowers(ctx, actorUserID)
	if err != nil {
		return dto.ListFollowsResponse{}, err
	}

	result := make([]dto.FollowItemResponse, 0, len(items))
	for _, item := range items {
		result = append(result, mapFollowFollower(item))
	}

	return dto.ListFollowsResponse{Items: result, Count: len(result)}, nil
}

func (s *SocialService) ListFollowing(ctx context.Context, request dto.ListFollowingRequest) (dto.ListFollowsResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ListFollowsResponse{}, err
	}

	actorUserID, err := parseID(request.ActorUserID, "actor_user_id")
	if err != nil {
		return dto.ListFollowsResponse{}, err
	}

	items, err := s.store.ListFollowing(ctx, actorUserID)
	if err != nil {
		return dto.ListFollowsResponse{}, err
	}

	result := make([]dto.FollowItemResponse, 0, len(items))
	for _, item := range items {
		result = append(result, mapFollowFollowing(item))
	}

	return dto.ListFollowsResponse{Items: result, Count: len(result)}, nil
}
