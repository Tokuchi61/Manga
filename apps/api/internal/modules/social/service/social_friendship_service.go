package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/social/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/social/entity"
	socialrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/social/repository"
)

func (s *SocialService) CreateFriendRequest(ctx context.Context, request dto.CreateFriendRequestRequest) (dto.CreateFriendRequestResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.CreateFriendRequestResponse{}, err
	}

	actorUserID, err := parseID(request.ActorUserID, "actor_user_id")
	if err != nil {
		return dto.CreateFriendRequestResponse{}, err
	}
	targetUserID, err := parseID(request.TargetUserID, "target_user_id")
	if err != nil {
		return dto.CreateFriendRequestResponse{}, err
	}
	if err := ensureDistinctActorTarget(actorUserID, targetUserID); err != nil {
		return dto.CreateFriendRequestResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.CreateFriendRequestResponse{}, err
	}
	if err := s.requireFriendshipWriteEnabled(cfg.FriendshipEnabled); err != nil {
		return dto.CreateFriendRequestResponse{}, err
	}

	blocked, err := s.store.IsBlockedEither(ctx, actorUserID, targetUserID)
	if err != nil {
		return dto.CreateFriendRequestResponse{}, err
	}
	if blocked {
		return dto.CreateFriendRequestResponse{}, ErrForbiddenAction
	}

	restricted, err := s.store.IsRestrictedEither(ctx, actorUserID, targetUserID)
	if err != nil {
		return dto.CreateFriendRequestResponse{}, err
	}
	if restricted {
		return dto.CreateFriendRequestResponse{}, ErrForbiddenAction
	}

	friends, err := s.store.AreFriends(ctx, actorUserID, targetUserID)
	if err != nil {
		return dto.CreateFriendRequestResponse{}, err
	}
	if friends {
		return dto.CreateFriendRequestResponse{}, ErrAlreadyFriends
	}

	now := s.now().UTC()
	stored, created, err := s.store.CreateFriendRequest(ctx, entity.FriendshipRequest{
		RequesterUserID: actorUserID,
		TargetUserID:    targetUserID,
		Status:          entity.FriendshipRequestPending,
		RequestID:       strings.TrimSpace(request.RequestID),
		CreatedAt:       now,
		UpdatedAt:       now,
	}, buildDedupKey("friendship", actorUserID, targetUserID, request.RequestID))
	if err != nil {
		if errors.Is(err, socialrepository.ErrConflict) {
			return dto.CreateFriendRequestResponse{}, ErrConflict
		}
		return dto.CreateFriendRequestResponse{}, err
	}

	return dto.CreateFriendRequestResponse{
		RequesterUserID: stored.RequesterUserID,
		TargetUserID:    stored.TargetUserID,
		Status:          string(stored.Status),
		Created:         created,
		UpdatedAt:       stored.UpdatedAt,
	}, nil
}

func (s *SocialService) RespondFriendRequest(ctx context.Context, request dto.RespondFriendRequestRequest) (dto.OperationResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.OperationResponse{}, err
	}

	actorUserID, err := parseID(request.ActorUserID, "actor_user_id")
	if err != nil {
		return dto.OperationResponse{}, err
	}
	requesterUserID, err := parseID(request.RequesterUserID, "requester_user_id")
	if err != nil {
		return dto.OperationResponse{}, err
	}
	if err := ensureDistinctActorTarget(actorUserID, requesterUserID); err != nil {
		return dto.OperationResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.OperationResponse{}, err
	}
	if err := s.requireFriendshipWriteEnabled(cfg.FriendshipEnabled); err != nil {
		return dto.OperationResponse{}, err
	}

	_, err = s.store.GetFriendRequest(ctx, requesterUserID, actorUserID)
	if err != nil {
		if errors.Is(err, socialrepository.ErrNotFound) {
			return dto.OperationResponse{}, ErrNotFound
		}
		return dto.OperationResponse{}, err
	}

	action := normalizeValue(request.Action)
	switch action {
	case "accept":
		blocked, err := s.store.IsBlockedEither(ctx, actorUserID, requesterUserID)
		if err != nil {
			return dto.OperationResponse{}, err
		}
		if blocked {
			return dto.OperationResponse{}, ErrForbiddenAction
		}

		restricted, err := s.store.IsRestrictedEither(ctx, actorUserID, requesterUserID)
		if err != nil {
			return dto.OperationResponse{}, err
		}
		if restricted {
			return dto.OperationResponse{}, ErrForbiddenAction
		}

		now := s.now().UTC()
		if err := s.store.UpsertFriendship(ctx, entity.Friendship{
			UserAID:   requesterUserID,
			UserBID:   actorUserID,
			CreatedAt: now,
			UpdatedAt: now,
		}); err != nil {
			return dto.OperationResponse{}, err
		}
		_ = s.store.DeleteFriendRequest(ctx, requesterUserID, actorUserID)
		_ = s.store.DeleteFriendRequest(ctx, actorUserID, requesterUserID)
		return dto.OperationResponse{Status: "accepted"}, nil
	case "reject":
		if err := s.store.DeleteFriendRequest(ctx, requesterUserID, actorUserID); err != nil {
			return dto.OperationResponse{}, err
		}
		return dto.OperationResponse{Status: "rejected"}, nil
	default:
		return dto.OperationResponse{}, ErrValidation
	}
}

func (s *SocialService) RemoveFriend(ctx context.Context, request dto.RemoveFriendRequest) (dto.OperationResponse, error) {
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
	if err := s.requireFriendshipWriteEnabled(cfg.FriendshipEnabled); err != nil {
		return dto.OperationResponse{}, err
	}

	if _, err := s.store.DeleteFriendship(ctx, actorUserID, targetUserID); err != nil {
		return dto.OperationResponse{}, err
	}
	_ = s.store.DeleteFriendRequest(ctx, actorUserID, targetUserID)
	_ = s.store.DeleteFriendRequest(ctx, targetUserID, actorUserID)

	return dto.OperationResponse{Status: "removed"}, nil
}

func (s *SocialService) ListFriendRequests(ctx context.Context, request dto.ListFriendRequestsRequest) (dto.ListFriendRequestsResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ListFriendRequestsResponse{}, err
	}

	actorUserID, err := parseID(request.ActorUserID, "actor_user_id")
	if err != nil {
		return dto.ListFriendRequestsResponse{}, err
	}

	items, err := s.store.ListFriendRequests(ctx, actorUserID, parseDirection(request.Direction))
	if err != nil {
		return dto.ListFriendRequestsResponse{}, err
	}

	result := make([]dto.FriendRequestItemResponse, 0, len(items))
	for _, item := range items {
		result = append(result, mapFriendRequest(item))
	}

	return dto.ListFriendRequestsResponse{Items: result, Count: len(result)}, nil
}

func (s *SocialService) ListFriends(ctx context.Context, request dto.ListFriendsRequest) (dto.ListFriendsResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ListFriendsResponse{}, err
	}

	actorUserID, err := parseID(request.ActorUserID, "actor_user_id")
	if err != nil {
		return dto.ListFriendsResponse{}, err
	}

	items, err := s.store.ListFriends(ctx, actorUserID)
	if err != nil {
		return dto.ListFriendsResponse{}, err
	}

	result := make([]dto.FriendItemResponse, 0, len(items))
	for _, item := range items {
		result = append(result, mapFriendship(item, actorUserID))
	}

	return dto.ListFriendsResponse{Items: result, Count: len(result)}, nil
}
