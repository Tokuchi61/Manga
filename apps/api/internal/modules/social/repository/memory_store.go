package repository

import (
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/social/entity"
)

// MemoryStore is stage-14 bootstrap persistence for social flows.
type MemoryStore struct {
	mu sync.RWMutex

	friendRequestsByKey map[string]entity.FriendshipRequest
	friendshipsByPair   map[string]entity.Friendship
	friendRequestDedup  map[string]string

	followsByKey map[string]entity.FollowRelation
	followDedup  map[string]string

	wallPostsByID     map[string]entity.WallPost
	wallRepliesByID   map[string]entity.WallReply
	wallRepliesByPost map[string][]string
	wallPostDedup     map[string]string
	wallReplyDedup    map[string]string

	threadsByID        map[string]entity.MessageThread
	threadByPair       map[string]string
	messagesByID       map[string]entity.Message
	messageIDsByThread map[string][]string
	messageDedup       map[string]string

	blockByKey    map[string]time.Time
	muteByKey     map[string]time.Time
	restrictByKey map[string]time.Time

	runtimeConfig entity.RuntimeConfig
}

func NewMemoryStore() *MemoryStore {
	now := time.Now().UTC()
	return &MemoryStore{
		friendRequestsByKey: make(map[string]entity.FriendshipRequest),
		friendshipsByPair:   make(map[string]entity.Friendship),
		friendRequestDedup:  make(map[string]string),
		followsByKey:        make(map[string]entity.FollowRelation),
		followDedup:         make(map[string]string),
		wallPostsByID:       make(map[string]entity.WallPost),
		wallRepliesByID:     make(map[string]entity.WallReply),
		wallRepliesByPost:   make(map[string][]string),
		wallPostDedup:       make(map[string]string),
		wallReplyDedup:      make(map[string]string),
		threadsByID:         make(map[string]entity.MessageThread),
		threadByPair:        make(map[string]string),
		messagesByID:        make(map[string]entity.Message),
		messageIDsByThread:  make(map[string][]string),
		messageDedup:        make(map[string]string),
		blockByKey:          make(map[string]time.Time),
		muteByKey:           make(map[string]time.Time),
		restrictByKey:       make(map[string]time.Time),
		runtimeConfig: entity.RuntimeConfig{
			FriendshipEnabled: true,
			FollowEnabled:     true,
			WallEnabled:       true,
			MessagingEnabled:  true,
			UpdatedAt:         now,
		},
	}
}

func normalizeValue(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

func directionalKey(left string, right string) string {
	return normalizeValue(left) + ":" + normalizeValue(right)
}

func pairKey(left string, right string) string {
	a := normalizeValue(left)
	b := normalizeValue(right)
	if a <= b {
		return a + ":" + b
	}
	return b + ":" + a
}

func cloneFriendRequest(in entity.FriendshipRequest) entity.FriendshipRequest {
	out := in
	out.CreatedAt = in.CreatedAt.UTC()
	out.UpdatedAt = in.UpdatedAt.UTC()
	return out
}

func cloneFriendship(in entity.Friendship) entity.Friendship {
	out := in
	out.CreatedAt = in.CreatedAt.UTC()
	out.UpdatedAt = in.UpdatedAt.UTC()
	return out
}

func cloneFollow(in entity.FollowRelation) entity.FollowRelation {
	out := in
	out.CreatedAt = in.CreatedAt.UTC()
	out.UpdatedAt = in.UpdatedAt.UTC()
	return out
}

func cloneWallPost(in entity.WallPost) entity.WallPost {
	out := in
	out.CreatedAt = in.CreatedAt.UTC()
	out.UpdatedAt = in.UpdatedAt.UTC()
	return out
}

func cloneWallReply(in entity.WallReply) entity.WallReply {
	out := in
	out.CreatedAt = in.CreatedAt.UTC()
	out.UpdatedAt = in.UpdatedAt.UTC()
	return out
}

func cloneThread(in entity.MessageThread) entity.MessageThread {
	out := in
	out.CreatedAt = in.CreatedAt.UTC()
	out.UpdatedAt = in.UpdatedAt.UTC()
	return out
}

func cloneMessage(in entity.Message) entity.Message {
	out := in
	out.CreatedAt = in.CreatedAt.UTC()
	return out
}

func cloneRuntimeConfig(in entity.RuntimeConfig) entity.RuntimeConfig {
	out := in
	out.UpdatedAt = in.UpdatedAt.UTC()
	return out
}

func sortByTimeDesc[T any](items []T, less func(i int, j int) bool) {
	sort.Slice(items, less)
}

func applyOffsetLimit[T any](items []T, offset int, limit int) []T {
	if offset < 0 {
		offset = 0
	}
	if offset >= len(items) {
		return []T{}
	}
	if limit <= 0 {
		limit = 50
	}
	end := offset + limit
	if end > len(items) {
		end = len(items)
	}
	return append([]T(nil), items[offset:end]...)
}
