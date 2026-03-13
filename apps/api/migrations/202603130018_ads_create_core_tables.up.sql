CREATE TABLE IF NOT EXISTS ads_placements (
    placement_id TEXT PRIMARY KEY,
    surface TEXT NOT NULL,
    target_type TEXT NOT NULL,
    target_id TEXT NOT NULL DEFAULT '',
    visible BOOLEAN NOT NULL DEFAULT TRUE,
    priority INTEGER NOT NULL DEFAULT 0,
    frequency_cap INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT ads_placements_surface_check CHECK (surface IN ('home', 'listing', 'manga', 'chapter')),
    CONSTRAINT ads_placements_target_type_check CHECK (target_type IN ('none', 'manga', 'chapter')),
    CONSTRAINT ads_placements_priority_check CHECK (priority >= 0),
    CONSTRAINT ads_placements_frequency_cap_check CHECK (frequency_cap >= 0)
);

CREATE INDEX IF NOT EXISTS idx_ads_placements_surface_visible_priority
    ON ads_placements (surface, visible, priority DESC, updated_at DESC);

CREATE TABLE IF NOT EXISTS ads_campaigns (
    campaign_id TEXT PRIMARY KEY,
    placement_id TEXT NOT NULL REFERENCES ads_placements(placement_id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    state TEXT NOT NULL,
    creative_url TEXT NOT NULL,
    click_url TEXT NOT NULL,
    weight INTEGER NOT NULL DEFAULT 1,
    starts_at TIMESTAMPTZ NULL,
    ends_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT ads_campaigns_state_check CHECK (state IN ('draft', 'active', 'paused', 'ended')),
    CONSTRAINT ads_campaigns_weight_check CHECK (weight > 0),
    CONSTRAINT ads_campaigns_window_check CHECK (starts_at IS NULL OR ends_at IS NULL OR starts_at <= ends_at)
);

CREATE INDEX IF NOT EXISTS idx_ads_campaigns_placement_state_weight
    ON ads_campaigns (placement_id, state, weight DESC, updated_at DESC);

CREATE INDEX IF NOT EXISTS idx_ads_campaigns_state_window
    ON ads_campaigns (state, starts_at, ends_at);

CREATE TABLE IF NOT EXISTS ads_impressions (
    impression_id UUID PRIMARY KEY,
    request_id TEXT NOT NULL,
    placement_id TEXT NOT NULL REFERENCES ads_placements(placement_id) ON DELETE CASCADE,
    campaign_id TEXT NOT NULL REFERENCES ads_campaigns(campaign_id) ON DELETE CASCADE,
    session_id TEXT NOT NULL DEFAULT '',
    user_id UUID NULL REFERENCES user_accounts(id) ON DELETE SET NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT ads_impressions_status_check CHECK (status IN ('accepted', 'ignored'))
);

CREATE INDEX IF NOT EXISTS idx_ads_impressions_session_placement_created
    ON ads_impressions (session_id, placement_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_ads_impressions_campaign_created
    ON ads_impressions (campaign_id, created_at DESC);

CREATE TABLE IF NOT EXISTS ads_impression_dedup (
    dedup_key TEXT PRIMARY KEY,
    impression_id UUID NOT NULL REFERENCES ads_impressions(impression_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ads_impression_dedup_created
    ON ads_impression_dedup (created_at DESC);

CREATE TABLE IF NOT EXISTS ads_clicks (
    click_id UUID PRIMARY KEY,
    request_id TEXT NOT NULL,
    placement_id TEXT NOT NULL REFERENCES ads_placements(placement_id) ON DELETE CASCADE,
    campaign_id TEXT NOT NULL REFERENCES ads_campaigns(campaign_id) ON DELETE CASCADE,
    session_id TEXT NOT NULL DEFAULT '',
    user_id UUID NULL REFERENCES user_accounts(id) ON DELETE SET NULL,
    status TEXT NOT NULL,
    invalid_traffic BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT ads_clicks_status_check CHECK (status IN ('accepted', 'ignored'))
);

CREATE INDEX IF NOT EXISTS idx_ads_clicks_session_campaign_created
    ON ads_clicks (session_id, campaign_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_ads_clicks_campaign_created
    ON ads_clicks (campaign_id, created_at DESC);

CREATE TABLE IF NOT EXISTS ads_click_dedup (
    dedup_key TEXT PRIMARY KEY,
    click_id UUID NOT NULL REFERENCES ads_clicks(click_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ads_click_dedup_created
    ON ads_click_dedup (created_at DESC);

CREATE TABLE IF NOT EXISTS ads_campaign_aggregates (
    campaign_id TEXT PRIMARY KEY REFERENCES ads_campaigns(campaign_id) ON DELETE CASCADE,
    impression_count BIGINT NOT NULL DEFAULT 0,
    click_count BIGINT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT ads_campaign_aggregates_impression_count_check CHECK (impression_count >= 0),
    CONSTRAINT ads_campaign_aggregates_click_count_check CHECK (click_count >= 0)
);

CREATE TABLE IF NOT EXISTS ads_runtime_controls (
    id SMALLINT PRIMARY KEY,
    surface_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    placement_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    campaign_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    click_intake_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO ads_runtime_controls (id)
VALUES (1)
ON CONFLICT (id) DO NOTHING;
