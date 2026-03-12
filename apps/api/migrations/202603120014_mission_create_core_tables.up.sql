CREATE TABLE IF NOT EXISTS mission_definitions (
    mission_id TEXT PRIMARY KEY,
    category TEXT NOT NULL,
    title TEXT NOT NULL,
    objective_type TEXT NOT NULL,
    target_count INTEGER NOT NULL,
    reward_item_id TEXT NOT NULL DEFAULT '',
    reward_quantity INTEGER NOT NULL DEFAULT 1,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    starts_at TIMESTAMPTZ NULL,
    ends_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT mission_definitions_target_count_check CHECK (target_count > 0),
    CONSTRAINT mission_definitions_reward_quantity_check CHECK (reward_quantity > 0),
    CONSTRAINT mission_definitions_window_check CHECK (starts_at IS NULL OR ends_at IS NULL OR starts_at <= ends_at)
);

CREATE INDEX IF NOT EXISTS idx_mission_definitions_category_active_updated_at
    ON mission_definitions (category, active, updated_at DESC);

CREATE TABLE IF NOT EXISTS mission_user_progress (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    mission_id TEXT NOT NULL REFERENCES mission_definitions(mission_id) ON DELETE CASCADE,
    period_key TEXT NOT NULL,
    progress_count INTEGER NOT NULL DEFAULT 0,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    claimed BOOLEAN NOT NULL DEFAULT FALSE,
    last_request_id TEXT NOT NULL DEFAULT '',
    last_correlation_id TEXT NOT NULL DEFAULT '',
    completed_at TIMESTAMPTZ NULL,
    claimed_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT mission_user_progress_progress_count_check CHECK (progress_count >= 0),
    CONSTRAINT mission_user_progress_unique_user_mission_period UNIQUE (user_id, mission_id, period_key)
);

CREATE INDEX IF NOT EXISTS idx_mission_user_progress_user_period_updated_at
    ON mission_user_progress (user_id, period_key, updated_at DESC);

CREATE TABLE IF NOT EXISTS mission_progress_dedup (
    dedup_key TEXT PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    mission_id TEXT NOT NULL REFERENCES mission_definitions(mission_id) ON DELETE CASCADE,
    period_key TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_mission_progress_dedup_user_created_at
    ON mission_progress_dedup (user_id, created_at DESC);

CREATE TABLE IF NOT EXISTS mission_claim_dedup (
    dedup_key TEXT PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    mission_id TEXT NOT NULL REFERENCES mission_definitions(mission_id) ON DELETE CASCADE,
    period_key TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_mission_claim_dedup_user_created_at
    ON mission_claim_dedup (user_id, created_at DESC);

CREATE TABLE IF NOT EXISTS mission_runtime_controls (
    id SMALLINT PRIMARY KEY,
    read_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    claim_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    progress_ingest_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    daily_reset_hour_utc SMALLINT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT mission_runtime_controls_daily_reset_hour_utc_check CHECK (daily_reset_hour_utc >= 0 AND daily_reset_hour_utc <= 23)
);

INSERT INTO mission_runtime_controls (id)
VALUES (1)
ON CONFLICT (id) DO NOTHING;
