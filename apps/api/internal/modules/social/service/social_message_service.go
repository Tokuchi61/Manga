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

func (s *SocialService) OpenThread(ctx context.Context, request dto.OpenThreadRequest) (dto.OpenThreadResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.OpenThreadResponse{}, err
	}

	actorUserID, err := parseID(request.ActorUserID, "actor_user_id")
	if err != nil {
		return dto.OpenThreadResponse{}, err
	}
	targetUserID, err := parseID(request.TargetUserID, "target_user_id")
	if err != nil {
		return dto.OpenThreadResponse{}, err
	}
	if err := ensureDistinctActorTarget(actorUserID, targetUserID); err != nil {
		return dto.OpenThreadResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.OpenThreadResponse{}, err
	}
	if err := s.requireMessagingWriteEnabled(cfg.MessagingEnabled); err != nil {
		return dto.OpenThreadResponse{}, err
	}

	blocked, err := s.store.IsBlockedEither(ctx, actorUserID, targetUserID)
	if err != nil {
		return dto.OpenThreadResponse{}, err
	}
	if blocked {
		return dto.OpenThreadResponse{}, ErrForbiddenAction
	}

	restricted, err := s.store.IsRestrictedEither(ctx, actorUserID, targetUserID)
	if err != nil {
		return dto.OpenThreadResponse{}, err
	}
	if restricted {
		return dto.OpenThreadResponse{}, ErrForbiddenAction
	}

	now := s.now().UTC()
	thread, created, err := s.store.OpenThread(ctx, entity.MessageThread{
		ID:        uuid.NewString(),
		UserAID:   actorUserID,
		UserBID:   targetUserID,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		if errors.Is(err, socialrepository.ErrConflict) {
			return dto.OpenThreadResponse{}, ErrConflict
		}
		return dto.OpenThreadResponse{}, err
	}

	peerUserID := resolveThreadPeer(thread, actorUserID)
	return dto.OpenThreadResponse{ThreadID: thread.ID, PeerUserID: peerUserID, Created: created}, nil
}

func (s *SocialService) SendMessage(ctx context.Context, request dto.SendMessageRequest) (dto.SendMessageResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.SendMessageResponse{}, err
	}

	actorUserID, err := parseID(request.ActorUserID, "actor_user_id")
	if err != nil {
		return dto.SendMessageResponse{}, err
	}
	threadID, err := parseID(request.ThreadID, "thread_id")
	if err != nil {
		return dto.SendMessageResponse{}, err
	}

	cfg, err := s.store.GetRuntimeConfig(ctx)
	if err != nil {
		return dto.SendMessageResponse{}, err
	}
	if err := s.requireMessagingWriteEnabled(cfg.MessagingEnabled); err != nil {
		return dto.SendMessageResponse{}, err
	}

	thread, err := s.store.GetThreadByID(ctx, threadID)
	if err != nil {
		if errors.Is(err, socialrepository.ErrNotFound) {
			return dto.SendMessageResponse{}, ErrNotFound
		}
		return dto.SendMessageResponse{}, err
	}
	if !isThreadParticipant(thread, actorUserID) {
		return dto.SendMessageResponse{}, ErrThreadAccessDenied
	}

	peerUserID := resolveThreadPeer(thread, actorUserID)
	blocked, err := s.store.IsBlockedEither(ctx, actorUserID, peerUserID)
	if err != nil {
		return dto.SendMessageResponse{}, err
	}
	if blocked {
		return dto.SendMessageResponse{}, ErrForbiddenAction
	}

	restricted, err := s.store.IsRestrictedEither(ctx, actorUserID, peerUserID)
	if err != nil {
		return dto.SendMessageResponse{}, err
	}
	if restricted {
		return dto.SendMessageResponse{}, ErrForbiddenAction
	}

	body, err := sanitizeBody(request.Body, "body", 4000)
	if err != nil {
		return dto.SendMessageResponse{}, err
	}

	now := s.now().UTC()
	storedMessage, created, err := s.store.CreateMessage(ctx, entity.Message{
		ID:            uuid.NewString(),
		ThreadID:      threadID,
		SenderUserID:  actorUserID,
		Body:          body,
		RequestID:     strings.TrimSpace(request.RequestID),
		CorrelationID: strings.TrimSpace(request.CorrelationID),
		CreatedAt:     now,
	}, buildDedupKey("message", threadID, actorUserID, request.RequestID, request.CorrelationID))
	if err != nil {
		if errors.Is(err, socialrepository.ErrNotFound) {
			return dto.SendMessageResponse{}, ErrNotFound
		}
		if errors.Is(err, socialrepository.ErrConflict) {
			return dto.SendMessageResponse{}, ErrConflict
		}
		return dto.SendMessageResponse{}, err
	}

	thread.LastMessageID = storedMessage.ID
	thread.UpdatedAt = now
	if normalizeValue(thread.UserAID) == normalizeValue(actorUserID) {
		thread.UnreadByB++
	} else {
		thread.UnreadByA++
	}
	if err := s.store.UpdateThread(ctx, thread); err != nil {
		if errors.Is(err, socialrepository.ErrNotFound) {
			return dto.SendMessageResponse{}, ErrNotFound
		}
		return dto.SendMessageResponse{}, err
	}

	return dto.SendMessageResponse{MessageID: storedMessage.ID, ThreadID: storedMessage.ThreadID, Created: created}, nil
}

func (s *SocialService) MarkThreadRead(ctx context.Context, request dto.MarkThreadReadRequest) (dto.OperationResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.OperationResponse{}, err
	}

	actorUserID, err := parseID(request.ActorUserID, "actor_user_id")
	if err != nil {
		return dto.OperationResponse{}, err
	}
	threadID, err := parseID(request.ThreadID, "thread_id")
	if err != nil {
		return dto.OperationResponse{}, err
	}

	thread, err := s.store.GetThreadByID(ctx, threadID)
	if err != nil {
		if errors.Is(err, socialrepository.ErrNotFound) {
			return dto.OperationResponse{}, ErrNotFound
		}
		return dto.OperationResponse{}, err
	}
	if !isThreadParticipant(thread, actorUserID) {
		return dto.OperationResponse{}, ErrThreadAccessDenied
	}

	if normalizeValue(thread.UserAID) == normalizeValue(actorUserID) {
		thread.UnreadByA = 0
	} else {
		thread.UnreadByB = 0
	}
	thread.UpdatedAt = s.now().UTC()

	if err := s.store.UpdateThread(ctx, thread); err != nil {
		if errors.Is(err, socialrepository.ErrNotFound) {
			return dto.OperationResponse{}, ErrNotFound
		}
		return dto.OperationResponse{}, err
	}

	return dto.OperationResponse{Status: "ok"}, nil
}

func (s *SocialService) ListThreads(ctx context.Context, request dto.ListThreadsRequest) (dto.ListThreadsResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ListThreadsResponse{}, err
	}

	actorUserID, err := parseID(request.ActorUserID, "actor_user_id")
	if err != nil {
		return dto.ListThreadsResponse{}, err
	}

	threads, err := s.store.ListThreadsByUser(ctx, actorUserID)
	if err != nil {
		return dto.ListThreadsResponse{}, err
	}

	items := make([]dto.ThreadItemResponse, 0, len(threads))
	for _, thread := range threads {
		items = append(items, mapThread(thread, actorUserID))
	}

	return dto.ListThreadsResponse{Items: items, Count: len(items)}, nil
}

func (s *SocialService) ListThreadMessages(ctx context.Context, request dto.ListThreadMessagesRequest) (dto.ListThreadMessagesResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.ListThreadMessagesResponse{}, err
	}

	actorUserID, err := parseID(request.ActorUserID, "actor_user_id")
	if err != nil {
		return dto.ListThreadMessagesResponse{}, err
	}
	threadID, err := parseID(request.ThreadID, "thread_id")
	if err != nil {
		return dto.ListThreadMessagesResponse{}, err
	}

	thread, err := s.store.GetThreadByID(ctx, threadID)
	if err != nil {
		if errors.Is(err, socialrepository.ErrNotFound) {
			return dto.ListThreadMessagesResponse{}, ErrNotFound
		}
		return dto.ListThreadMessagesResponse{}, err
	}
	if !isThreadParticipant(thread, actorUserID) {
		return dto.ListThreadMessagesResponse{}, ErrThreadAccessDenied
	}

	limit := request.Limit
	if limit <= 0 {
		limit = 50
	}

	messages, err := s.store.ListMessagesByThreadID(ctx, threadID, parseSortBy(request.SortBy, "newest"), limit, request.Offset)
	if err != nil {
		return dto.ListThreadMessagesResponse{}, err
	}

	items := make([]dto.MessageItemResponse, 0, len(messages))
	for _, message := range messages {
		items = append(items, mapMessage(message))
	}

	return dto.ListThreadMessagesResponse{Items: items, Count: len(items)}, nil
}
