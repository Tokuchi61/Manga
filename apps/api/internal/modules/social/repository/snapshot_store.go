package repository

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/social/entity"
)

type socialSnapshotState struct {
	FriendRequestsByKey map[string]entity.FriendshipRequest
	FriendshipsByPair   map[string]entity.Friendship
	FriendRequestDedup  map[string]string

	FollowsByKey map[string]entity.FollowRelation
	FollowDedup  map[string]string

	WallPostsByID     map[string]entity.WallPost
	WallRepliesByID   map[string]entity.WallReply
	WallRepliesByPost map[string][]string
	WallPostDedup     map[string]string
	WallReplyDedup    map[string]string

	ThreadsByID        map[string]entity.MessageThread
	ThreadByPair       map[string]string
	MessagesByID       map[string]entity.Message
	MessageIDsByThread map[string][]string
	MessageDedup       map[string]string

	BlockByKey    map[string]time.Time
	MuteByKey     map[string]time.Time
	RestrictByKey map[string]time.Time

	RuntimeConfig entity.RuntimeConfig
}

func (s *MemoryStore) Snapshot() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state := socialSnapshotState{
		FriendRequestsByKey: s.friendRequestsByKey,
		FriendshipsByPair:   s.friendshipsByPair,
		FriendRequestDedup:  s.friendRequestDedup,
		FollowsByKey:        s.followsByKey,
		FollowDedup:         s.followDedup,
		WallPostsByID:       s.wallPostsByID,
		WallRepliesByID:     s.wallRepliesByID,
		WallRepliesByPost:   s.wallRepliesByPost,
		WallPostDedup:       s.wallPostDedup,
		WallReplyDedup:      s.wallReplyDedup,
		ThreadsByID:         s.threadsByID,
		ThreadByPair:        s.threadByPair,
		MessagesByID:        s.messagesByID,
		MessageIDsByThread:  s.messageIDsByThread,
		MessageDedup:        s.messageDedup,
		BlockByKey:          s.blockByKey,
		MuteByKey:           s.muteByKey,
		RestrictByKey:       s.restrictByKey,
		RuntimeConfig:       s.runtimeConfig,
	}

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(state); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *MemoryStore) RestoreSnapshot(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	var state socialSnapshotState
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&state); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if state.FriendRequestsByKey == nil {
		state.FriendRequestsByKey = make(map[string]entity.FriendshipRequest)
	}
	if state.FriendshipsByPair == nil {
		state.FriendshipsByPair = make(map[string]entity.Friendship)
	}
	if state.FriendRequestDedup == nil {
		state.FriendRequestDedup = make(map[string]string)
	}
	if state.FollowsByKey == nil {
		state.FollowsByKey = make(map[string]entity.FollowRelation)
	}
	if state.FollowDedup == nil {
		state.FollowDedup = make(map[string]string)
	}
	if state.WallPostsByID == nil {
		state.WallPostsByID = make(map[string]entity.WallPost)
	}
	if state.WallRepliesByID == nil {
		state.WallRepliesByID = make(map[string]entity.WallReply)
	}
	if state.WallRepliesByPost == nil {
		state.WallRepliesByPost = make(map[string][]string)
	}
	if state.WallPostDedup == nil {
		state.WallPostDedup = make(map[string]string)
	}
	if state.WallReplyDedup == nil {
		state.WallReplyDedup = make(map[string]string)
	}
	if state.ThreadsByID == nil {
		state.ThreadsByID = make(map[string]entity.MessageThread)
	}
	if state.ThreadByPair == nil {
		state.ThreadByPair = make(map[string]string)
	}
	if state.MessagesByID == nil {
		state.MessagesByID = make(map[string]entity.Message)
	}
	if state.MessageIDsByThread == nil {
		state.MessageIDsByThread = make(map[string][]string)
	}
	if state.MessageDedup == nil {
		state.MessageDedup = make(map[string]string)
	}
	if state.BlockByKey == nil {
		state.BlockByKey = make(map[string]time.Time)
	}
	if state.MuteByKey == nil {
		state.MuteByKey = make(map[string]time.Time)
	}
	if state.RestrictByKey == nil {
		state.RestrictByKey = make(map[string]time.Time)
	}
	if state.RuntimeConfig.UpdatedAt.IsZero() {
		state.RuntimeConfig = entity.RuntimeConfig{
			FriendshipEnabled: true,
			FollowEnabled:     true,
			WallEnabled:       true,
			MessagingEnabled:  true,
			UpdatedAt:         time.Now().UTC(),
		}
	}

	s.friendRequestsByKey = state.FriendRequestsByKey
	s.friendshipsByPair = state.FriendshipsByPair
	s.friendRequestDedup = state.FriendRequestDedup
	s.followsByKey = state.FollowsByKey
	s.followDedup = state.FollowDedup
	s.wallPostsByID = state.WallPostsByID
	s.wallRepliesByID = state.WallRepliesByID
	s.wallRepliesByPost = state.WallRepliesByPost
	s.wallPostDedup = state.WallPostDedup
	s.wallReplyDedup = state.WallReplyDedup
	s.threadsByID = state.ThreadsByID
	s.threadByPair = state.ThreadByPair
	s.messagesByID = state.MessagesByID
	s.messageIDsByThread = state.MessageIDsByThread
	s.messageDedup = state.MessageDedup
	s.blockByKey = state.BlockByKey
	s.muteByKey = state.MuteByKey
	s.restrictByKey = state.RestrictByKey
	s.runtimeConfig = state.RuntimeConfig

	return nil
}
