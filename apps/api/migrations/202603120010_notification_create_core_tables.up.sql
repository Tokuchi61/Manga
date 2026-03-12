CREATE TABLE IF NOT EXISTS notification_entries (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    category TEXT NOT NULL,
    channel TEXT NOT NULL DEFAULT 'in_app',
    template_key TEXT NOT NULL DEFAULT '',
    title TEXT NOT NULL,
    body TEXT NOT NULL,
    state TEXT NOT NULL DEFAULT 'created',
    source_event TEXT NOT NULL DEFAULT '',
    source_ref_id UUID NULL,
    request_id TEXT NOT NULL DEFAULT '',
    correlation_id TEXT NOT NULL DEFAULT '',
    dedup_key TEXT NOT NULL DEFAULT '',
    delivery_attempt_count INT NOT NULL DEFAULT 0,
    last_failure_reason TEXT NOT NULL DEFAULT '',
    read_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT notification_entries_category_check CHECK (category IN ('account_security', 'social', 'comment', 'support', 'moderation', 'mission', 'royalpass', 'shop', 'payment', 'system_ops')),
    CONSTRAINT notification_entries_channel_check CHECK (channel IN ('in_app', 'email', 'push')),
    CONSTRAINT notification_entries_state_check CHECK (state IN ('created', 'delivered', 'failed', 'read'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_notification_entries_dedup_key
    ON notification_entries (dedup_key)
    WHERE dedup_key <> '';

CREATE INDEX IF NOT EXISTS idx_notification_entries_user_created_at
    ON notification_entries (user_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_notification_entries_user_state
    ON notification_entries (user_id, state);

CREATE TABLE IF NOT EXISTS notification_delivery_attempts (
    id UUID PRIMARY KEY,
    notification_id UUID NOT NULL REFERENCES notification_entries(id) ON DELETE CASCADE,
    channel TEXT NOT NULL,
    attempt_no INT NOT NULL,
    result TEXT NOT NULL DEFAULT 'pending',
    provider_ref TEXT NOT NULL DEFAULT '',
    failure_reason TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT notification_delivery_attempts_channel_check CHECK (channel IN ('in_app', 'email', 'push')),
    CONSTRAINT notification_delivery_attempts_result_check CHECK (result IN ('pending', 'delivered', 'failed')),
    CONSTRAINT notification_delivery_attempts_attempt_no_check CHECK (attempt_no > 0),
    CONSTRAINT notification_delivery_attempts_unique UNIQUE (notification_id, attempt_no)
);

CREATE INDEX IF NOT EXISTS idx_notification_delivery_attempts_notification_id_created_at
    ON notification_delivery_attempts (notification_id, created_at);

CREATE TABLE IF NOT EXISTS notification_preferences (
    user_id UUID PRIMARY KEY REFERENCES user_accounts(id) ON DELETE CASCADE,
    muted_categories TEXT[] NOT NULL DEFAULT '{}',
    quiet_hours_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    quiet_hours_start SMALLINT NOT NULL DEFAULT 22,
    quiet_hours_end SMALLINT NOT NULL DEFAULT 7,
    channel_in_app_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    channel_email_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    channel_push_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    digest_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT notification_preferences_quiet_hours_start_check CHECK (quiet_hours_start BETWEEN 0 AND 23),
    CONSTRAINT notification_preferences_quiet_hours_end_check CHECK (quiet_hours_end BETWEEN 0 AND 23)
);

CREATE TABLE IF NOT EXISTS notification_runtime_controls (
    id SMALLINT PRIMARY KEY,
    delivery_paused BOOLEAN NOT NULL DEFAULT FALSE,
    digest_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    category_state JSONB NOT NULL DEFAULT '{}'::jsonb,
    channel_state JSONB NOT NULL DEFAULT '{"in_app": true, "email": true, "push": true}'::jsonb,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO notification_runtime_controls (id)
VALUES (1)
ON CONFLICT (id) DO NOTHING;
