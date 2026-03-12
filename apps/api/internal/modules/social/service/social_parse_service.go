package service

import (
	"fmt"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/social/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/social/entity"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/social/repository"
	"github.com/google/uuid"
)

func normalizeValue(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

func parseID(raw string, fieldName string) (string, error) {
	parsed, err := uuid.Parse(strings.TrimSpace(raw))
	if err != nil {
		return "", fmt.Errorf("%w: invalid %s", ErrValidation, fieldName)
	}
	return parsed.String(), nil
}

func parseSortBy(raw string, fallback string) string {
	value := normalizeValue(raw)
	switch value {
	case "newest", "oldest":
		return value
	default:
		if fallback == "" {
			return "newest"
		}
		return fallback
	}
}

func parseDirection(raw string) string {
	value := normalizeValue(raw)
	switch value {
	case "incoming", "outgoing":
		return value
	default:
		return ""
	}
}

func sanitizeBody(raw string, fieldName string, maxLen int) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", fmt.Errorf("%w: %s cannot be empty", ErrValidation, fieldName)
	}
	if maxLen > 0 && len(trimmed) > maxLen {
		return "", fmt.Errorf("%w: %s too long", ErrValidation, fieldName)
	}
	trimmed = strings.ReplaceAll(trimmed, "<", "&lt;")
	trimmed = strings.ReplaceAll(trimmed, ">", "&gt;")
	return trimmed, nil
}

func ensureDistinctActorTarget(actorUserID string, targetUserID string) error {
	if normalizeValue(actorUserID) == normalizeValue(targetUserID) {
		return fmt.Errorf("%w: actor and target cannot be same", ErrValidation)
	}
	return nil
}

func buildDedupKey(parts ...string) string {
	items := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := normalizeValue(part)
		if trimmed == "" {
			continue
		}
		items = append(items, trimmed)
	}
	if len(items) == 0 {
		return ""
	}
	return strings.Join(items, ":")
}

func mapRuntimeConfig(cfg entity.RuntimeConfig) dto.RuntimeConfigResponse {
	return dto.RuntimeConfigResponse{
		FriendshipEnabled: cfg.FriendshipEnabled,
		FollowEnabled:     cfg.FollowEnabled,
		WallEnabled:       cfg.WallEnabled,
		MessagingEnabled:  cfg.MessagingEnabled,
		UpdatedAt:         cfg.UpdatedAt,
	}
}

func mapFriendRequest(item entity.FriendshipRequest) dto.FriendRequestItemResponse {
	return dto.FriendRequestItemResponse{
		RequesterUserID: item.RequesterUserID,
		TargetUserID:    item.TargetUserID,
		Status:          string(item.Status),
		UpdatedAt:       item.UpdatedAt,
	}
}

func mapFriendship(item entity.Friendship, actorUserID string) dto.FriendItemResponse {
	actor := normalizeValue(actorUserID)
	peer := item.UserAID
	if normalizeValue(item.UserAID) == actor {
		peer = item.UserBID
	}
	return dto.FriendItemResponse{
		UserID:    peer,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}

func mapFollowFollower(item entity.FollowRelation) dto.FollowItemResponse {
	return dto.FollowItemResponse{
		UserID:    item.FollowerUserID,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}

func mapFollowFollowing(item entity.FollowRelation) dto.FollowItemResponse {
	return dto.FollowItemResponse{
		UserID:    item.FolloweeUserID,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}

func mapWallReply(item entity.WallReply) dto.WallReplyItemResponse {
	return dto.WallReplyItemResponse{
		ReplyID:     item.ID,
		PostID:      item.PostID,
		OwnerUserID: item.OwnerUserID,
		Body:        item.Body,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}

func mapWallPost(item entity.WallPost, replies []entity.WallReply) dto.WallPostItemResponse {
	response := dto.WallPostItemResponse{
		PostID:      item.ID,
		OwnerUserID: item.OwnerUserID,
		Body:        item.Body,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
	if len(replies) == 0 {
		return response
	}
	response.Replies = make([]dto.WallReplyItemResponse, 0, len(replies))
	for _, reply := range replies {
		response.Replies = append(response.Replies, mapWallReply(reply))
	}
	return response
}

func mapThread(thread entity.MessageThread, actorUserID string) dto.ThreadItemResponse {
	actor := normalizeValue(actorUserID)
	peerUserID := thread.UserAID
	unread := thread.UnreadByB
	if normalizeValue(thread.UserAID) == actor {
		peerUserID = thread.UserBID
		unread = thread.UnreadByA
	}

	return dto.ThreadItemResponse{
		ThreadID:      thread.ID,
		PeerUserID:    peerUserID,
		LastMessageID: thread.LastMessageID,
		UnreadCount:   unread,
		CreatedAt:     thread.CreatedAt,
		UpdatedAt:     thread.UpdatedAt,
	}
}

func mapMessage(item entity.Message) dto.MessageItemResponse {
	return dto.MessageItemResponse{
		MessageID:     item.ID,
		ThreadID:      item.ThreadID,
		SenderUserID:  item.SenderUserID,
		Body:          item.Body,
		RequestID:     item.RequestID,
		CorrelationID: item.CorrelationID,
		CreatedAt:     item.CreatedAt,
	}
}

func mapRelation(item repository.RelationEntry, relationType string, enabled bool) dto.RelationUpdateResponse {
	return dto.RelationUpdateResponse{
		ActorUserID:  item.ActorUserID,
		TargetUserID: item.TargetUserID,
		RelationType: relationType,
		Enabled:      enabled,
		UpdatedAt:    item.UpdatedAt,
	}
}

func mapRelationItem(item repository.RelationEntry) dto.RelationItemResponse {
	return dto.RelationItemResponse{
		TargetUserID: item.TargetUserID,
		UpdatedAt:    item.UpdatedAt,
	}
}

func isThreadParticipant(thread entity.MessageThread, actorUserID string) bool {
	actor := normalizeValue(actorUserID)
	return normalizeValue(thread.UserAID) == actor || normalizeValue(thread.UserBID) == actor
}

func resolveThreadPeer(thread entity.MessageThread, actorUserID string) string {
	actor := normalizeValue(actorUserID)
	if normalizeValue(thread.UserAID) == actor {
		return thread.UserBID
	}
	return thread.UserAID
}
