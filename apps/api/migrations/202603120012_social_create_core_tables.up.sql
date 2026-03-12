CREATE TABLE IF NOT EXISTS social_friendship_requests (
    id UUID PRIMARY KEY,
    requester_user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    target_user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    request_id TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT social_friendship_requests_status_check CHECK (status IN ('pending', 'rejected')),
    CONSTRAINT social_friendship_requests_actor_target_check CHECK (requester_user_id <> target_user_id),
    CONSTRAINT social_friendship_requests_unique_pair UNIQUE (requester_user_id, target_user_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_social_friendship_requests_request_id
    ON social_friendship_requests (requester_user_id, request_id)
    WHERE request_id <> '';

CREATE INDEX IF NOT EXISTS idx_social_friendship_requests_target_updated_at
    ON social_friendship_requests (target_user_id, updated_at DESC);

CREATE TABLE IF NOT EXISTS social_friendships (
    id UUID PRIMARY KEY,
    user_a_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    user_b_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT social_friendships_actor_target_check CHECK (user_a_id <> user_b_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_social_friendships_pair
    ON social_friendships (LEAST(user_a_id, user_b_id), GREATEST(user_a_id, user_b_id));

CREATE INDEX IF NOT EXISTS idx_social_friendships_user_a_updated_at
    ON social_friendships (user_a_id, updated_at DESC);

CREATE INDEX IF NOT EXISTS idx_social_friendships_user_b_updated_at
    ON social_friendships (user_b_id, updated_at DESC);

CREATE TABLE IF NOT EXISTS social_follows (
    id UUID PRIMARY KEY,
    follower_user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    followee_user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    request_id TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT social_follows_actor_target_check CHECK (follower_user_id <> followee_user_id),
    CONSTRAINT social_follows_unique_pair UNIQUE (follower_user_id, followee_user_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_social_follows_request_id
    ON social_follows (follower_user_id, request_id)
    WHERE request_id <> '';

CREATE INDEX IF NOT EXISTS idx_social_follows_followee_updated_at
    ON social_follows (followee_user_id, updated_at DESC);

CREATE TABLE IF NOT EXISTS social_wall_posts (
    id UUID PRIMARY KEY,
    owner_user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    body TEXT NOT NULL,
    request_id TEXT NOT NULL DEFAULT '',
    correlation_id TEXT NOT NULL DEFAULT '',
    dedup_key TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_social_wall_posts_dedup_key
    ON social_wall_posts (dedup_key)
    WHERE dedup_key <> '';

CREATE INDEX IF NOT EXISTS idx_social_wall_posts_owner_created_at
    ON social_wall_posts (owner_user_id, created_at DESC);

CREATE TABLE IF NOT EXISTS social_wall_replies (
    id UUID PRIMARY KEY,
    post_id UUID NOT NULL REFERENCES social_wall_posts(id) ON DELETE CASCADE,
    owner_user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    body TEXT NOT NULL,
    request_id TEXT NOT NULL DEFAULT '',
    correlation_id TEXT NOT NULL DEFAULT '',
    dedup_key TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_social_wall_replies_dedup_key
    ON social_wall_replies (dedup_key)
    WHERE dedup_key <> '';

CREATE INDEX IF NOT EXISTS idx_social_wall_replies_post_created_at
    ON social_wall_replies (post_id, created_at ASC);

CREATE TABLE IF NOT EXISTS social_message_threads (
    id UUID PRIMARY KEY,
    user_a_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    user_b_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    last_message_id UUID NULL,
    unread_by_a INTEGER NOT NULL DEFAULT 0,
    unread_by_b INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT social_message_threads_actor_target_check CHECK (user_a_id <> user_b_id),
    CONSTRAINT social_message_threads_unread_by_a_check CHECK (unread_by_a >= 0),
    CONSTRAINT social_message_threads_unread_by_b_check CHECK (unread_by_b >= 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_social_message_threads_pair
    ON social_message_threads (LEAST(user_a_id, user_b_id), GREATEST(user_a_id, user_b_id));

CREATE INDEX IF NOT EXISTS idx_social_message_threads_user_a_updated_at
    ON social_message_threads (user_a_id, updated_at DESC);

CREATE INDEX IF NOT EXISTS idx_social_message_threads_user_b_updated_at
    ON social_message_threads (user_b_id, updated_at DESC);

CREATE TABLE IF NOT EXISTS social_messages (
    id UUID PRIMARY KEY,
    thread_id UUID NOT NULL REFERENCES social_message_threads(id) ON DELETE CASCADE,
    sender_user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    body TEXT NOT NULL,
    request_id TEXT NOT NULL DEFAULT '',
    correlation_id TEXT NOT NULL DEFAULT '',
    dedup_key TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_social_messages_dedup_key
    ON social_messages (dedup_key)
    WHERE dedup_key <> '';

CREATE INDEX IF NOT EXISTS idx_social_messages_thread_created_at
    ON social_messages (thread_id, created_at DESC);

CREATE TABLE IF NOT EXISTS social_relationship_edges (
    id UUID PRIMARY KEY,
    actor_user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    target_user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    relation_type TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT social_relationship_edges_type_check CHECK (relation_type IN ('blocked', 'muted', 'restricted')),
    CONSTRAINT social_relationship_edges_actor_target_check CHECK (actor_user_id <> target_user_id),
    CONSTRAINT social_relationship_edges_unique UNIQUE (actor_user_id, target_user_id, relation_type)
);

CREATE INDEX IF NOT EXISTS idx_social_relationship_edges_actor_type_updated_at
    ON social_relationship_edges (actor_user_id, relation_type, updated_at DESC);

CREATE TABLE IF NOT EXISTS social_runtime_controls (
    id SMALLINT PRIMARY KEY,
    friendship_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    follow_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    wall_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    messaging_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO social_runtime_controls (id)
VALUES (1)
ON CONFLICT (id) DO NOTHING;
