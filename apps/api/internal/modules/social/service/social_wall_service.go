package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/social/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/social/entity"
	socialrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/social/repository"
	"github.com/google/uuid"
)

func (s *SocialService) CreateWallPost(ctx context.Context, request dto.CreateWallPostRequest) (dto.WallWriteResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.WallWriteResponse{}, err
	}

	actorUserID, err := parseID(request.ActorUserID, "actor_user_id")
	if err != nil {
		return dto.WallWriteResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.WallWriteResponse{}, err
	}
	if err := s.requireWallEnabled(cfg.WallEnabled); err != nil {
		return dto.WallWriteResponse{}, err
	}

	body, err := sanitizeBody(request.Body, "body", 2000)
	if err != nil {
		return dto.WallWriteResponse{}, err
	}

	now := s.now().UTC()
	stored, created, err := s.store.CreateWallPost(ctx, entity.WallPost{
		ID:          uuid.NewString(),
		OwnerUserID: actorUserID,
		Body:        body,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, buildDedupKey("wall_post", actorUserID, request.RequestID, request.CorrelationID, body))
	if err != nil {
		if errors.Is(err, socialrepository.ErrConflict) {
			return dto.WallWriteResponse{}, ErrConflict
		}
		return dto.WallWriteResponse{}, err
	}

	return dto.WallWriteResponse{
		PostID:      stored.ID,
		OwnerUserID: stored.OwnerUserID,
		Created:     created,
		UpdatedAt:   stored.UpdatedAt,
	}, nil
}

func (s *SocialService) CreateWallReply(ctx context.Context, request dto.CreateWallReplyRequest) (dto.WallWriteResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.WallWriteResponse{}, err
	}

	actorUserID, err := parseID(request.ActorUserID, "actor_user_id")
	if err != nil {
		return dto.WallWriteResponse{}, err
	}
	postID, err := parseID(request.PostID, "post_id")
	if err != nil {
		return dto.WallWriteResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.WallWriteResponse{}, err
	}
	if err := s.requireWallEnabled(cfg.WallEnabled); err != nil {
		return dto.WallWriteResponse{}, err
	}

	post, err := s.store.GetWallPostByID(ctx, postID)
	if err != nil {
		if errors.Is(err, socialrepository.ErrNotFound) {
			return dto.WallWriteResponse{}, ErrNotFound
		}
		return dto.WallWriteResponse{}, err
	}

	blocked, err := s.store.IsBlockedEither(ctx, actorUserID, post.OwnerUserID)
	if err != nil {
		return dto.WallWriteResponse{}, err
	}
	if blocked {
		return dto.WallWriteResponse{}, ErrForbiddenAction
	}

	restricted, err := s.store.IsRestrictedEither(ctx, actorUserID, post.OwnerUserID)
	if err != nil {
		return dto.WallWriteResponse{}, err
	}
	if restricted {
		return dto.WallWriteResponse{}, ErrForbiddenAction
	}

	body, err := sanitizeBody(request.Body, "body", 2000)
	if err != nil {
		return dto.WallWriteResponse{}, err
	}

	now := s.now().UTC()
	stored, created, err := s.store.CreateWallReply(ctx, entity.WallReply{
		ID:          uuid.NewString(),
		PostID:      postID,
		OwnerUserID: actorUserID,
		Body:        body,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, buildDedupKey("wall_reply", actorUserID, postID, request.RequestID, request.CorrelationID, body))
	if err != nil {
		if errors.Is(err, socialrepository.ErrNotFound) {
			return dto.WallWriteResponse{}, ErrNotFound
		}
		if errors.Is(err, socialrepository.ErrConflict) {
			return dto.WallWriteResponse{}, ErrConflict
		}
		return dto.WallWriteResponse{}, err
	}

	return dto.WallWriteResponse{
		ReplyID:     stored.ID,
		OwnerUserID: stored.OwnerUserID,
		Created:     created,
		UpdatedAt:   stored.UpdatedAt,
	}, nil
}

func (s *SocialService) ListWallPosts(ctx context.Context, request dto.ListWallPostsRequest) (dto.ListWallPostsResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ListWallPostsResponse{}, err
	}

	actorUserID, err := parseID(request.ActorUserID, "actor_user_id")
	if err != nil {
		return dto.ListWallPostsResponse{}, err
	}

	ownerUserID := strings.TrimSpace(request.OwnerUserID)
	if ownerUserID == "" {
		ownerUserID = actorUserID
	} else {
		ownerUserID, err = parseID(ownerUserID, "owner_user_id")
		if err != nil {
			return dto.ListWallPostsResponse{}, err
		}
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.ListWallPostsResponse{}, err
	}
	if err := s.requireWallEnabled(cfg.WallEnabled); err != nil {
		return dto.ListWallPostsResponse{}, err
	}

	if normalizeValue(actorUserID) != normalizeValue(ownerUserID) {
		blocked, err := s.store.IsBlockedEither(ctx, actorUserID, ownerUserID)
		if err != nil {
			return dto.ListWallPostsResponse{}, err
		}
		if blocked {
			return dto.ListWallPostsResponse{}, ErrForbiddenAction
		}

		restricted, err := s.store.IsRestrictedEither(ctx, actorUserID, ownerUserID)
		if err != nil {
			return dto.ListWallPostsResponse{}, err
		}
		if restricted {
			return dto.ListWallPostsResponse{}, ErrForbiddenAction
		}
	}

	limit := request.Limit
	if limit <= 0 {
		limit = 50
	}

	posts, err := s.store.ListWallPosts(ctx, ownerUserID, parseSortBy(request.SortBy, "newest"), limit, request.Offset)
	if err != nil {
		return dto.ListWallPostsResponse{}, err
	}

	result := make([]dto.WallPostItemResponse, 0, len(posts))
	for _, post := range posts {
		replies := []entity.WallReply{}
		if request.IncludeReplies {
			replies, err = s.store.ListWallRepliesByPostID(ctx, post.ID, "oldest")
			if err != nil {
				return dto.ListWallPostsResponse{}, err
			}
		}
		result = append(result, mapWallPost(post, replies))
	}

	return dto.ListWallPostsResponse{Items: result, Count: len(result)}, nil
}
