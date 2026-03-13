CREATE TABLE IF NOT EXISTS royalpass_seasons (
    season_id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    state TEXT NOT NULL,
    starts_at TIMESTAMPTZ NULL,
    ends_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT royalpass_seasons_state_check CHECK (state IN ('draft', 'active', 'paused', 'ended', 'archived')),
    CONSTRAINT royalpass_seasons_window_check CHECK (starts_at IS NULL OR ends_at IS NULL OR starts_at <= ends_at)
);

CREATE INDEX IF NOT EXISTS idx_royalpass_seasons_state_updated_at
    ON royalpass_seasons (state, updated_at DESC);

CREATE TABLE IF NOT EXISTS royalpass_tiers (
    id UUID PRIMARY KEY,
    season_id TEXT NOT NULL REFERENCES royalpass_seasons(season_id) ON DELETE CASCADE,
    tier_number INTEGER NOT NULL,
    track TEXT NOT NULL,
    required_points INTEGER NOT NULL,
    reward_item_id TEXT NOT NULL DEFAULT '',
    reward_quantity INTEGER NOT NULL DEFAULT 1,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT royalpass_tiers_track_check CHECK (track IN ('free', 'premium')),
    CONSTRAINT royalpass_tiers_required_points_check CHECK (required_points > 0),
    CONSTRAINT royalpass_tiers_reward_quantity_check CHECK (reward_quantity > 0),
    CONSTRAINT royalpass_tiers_unique_season_tier_track UNIQUE (season_id, tier_number, track)
);

CREATE INDEX IF NOT EXISTS idx_royalpass_tiers_season_track_tier
    ON royalpass_tiers (season_id, track, tier_number);

CREATE TABLE IF NOT EXISTS royalpass_user_progress (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    season_id TEXT NOT NULL REFERENCES royalpass_seasons(season_id) ON DELETE CASCADE,
    points INTEGER NOT NULL DEFAULT 0,
    premium_activated BOOLEAN NOT NULL DEFAULT FALSE,
    premium_activation_source TEXT NOT NULL DEFAULT '',
    premium_activation_ref TEXT NOT NULL DEFAULT '',
    claimed_tiers_json JSONB NOT NULL DEFAULT '[]'::jsonb,
    last_request_id TEXT NOT NULL DEFAULT '',
    last_correlation_id TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT royalpass_user_progress_points_check CHECK (points >= 0),
    CONSTRAINT royalpass_user_progress_unique_user_season UNIQUE (user_id, season_id)
);

CREATE INDEX IF NOT EXISTS idx_royalpass_user_progress_user_updated_at
    ON royalpass_user_progress (user_id, updated_at DESC);

CREATE TABLE IF NOT EXISTS royalpass_progress_dedup (
    dedup_key TEXT PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    season_id TEXT NOT NULL REFERENCES royalpass_seasons(season_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_royalpass_progress_dedup_user_created_at
    ON royalpass_progress_dedup (user_id, created_at DESC);

CREATE TABLE IF NOT EXISTS royalpass_claim_dedup (
    dedup_key TEXT PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    season_id TEXT NOT NULL REFERENCES royalpass_seasons(season_id) ON DELETE CASCADE,
    tier_number INTEGER NOT NULL,
    track TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT royalpass_claim_dedup_track_check CHECK (track IN ('free', 'premium'))
);

CREATE INDEX IF NOT EXISTS idx_royalpass_claim_dedup_user_created_at
    ON royalpass_claim_dedup (user_id, created_at DESC);

CREATE TABLE IF NOT EXISTS royalpass_premium_activation_dedup (
    dedup_key TEXT PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    season_id TEXT NOT NULL REFERENCES royalpass_seasons(season_id) ON DELETE CASCADE,
    source_type TEXT NOT NULL,
    activation_ref TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_royalpass_premium_activation_dedup_user_created_at
    ON royalpass_premium_activation_dedup (user_id, created_at DESC);

CREATE TABLE IF NOT EXISTS royalpass_runtime_controls (
    id SMALLINT PRIMARY KEY,
    season_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    claim_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    premium_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO royalpass_runtime_controls (id)
VALUES (1)
ON CONFLICT (id) DO NOTHING;
