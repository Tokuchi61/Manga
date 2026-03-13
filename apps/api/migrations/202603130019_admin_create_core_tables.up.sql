CREATE TABLE IF NOT EXISTS admin_runtime_controls (
    id SMALLINT PRIMARY KEY,
    maintenance_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO admin_runtime_controls (id)
VALUES (1)
ON CONFLICT (id) DO NOTHING;

CREATE TABLE IF NOT EXISTS admin_actions (
    action_id UUID PRIMARY KEY,
    action_type TEXT NOT NULL,
    actor_user_id UUID NULL REFERENCES user_accounts(id) ON DELETE SET NULL,
    target_user_id UUID NULL REFERENCES user_accounts(id) ON DELETE SET NULL,
    target_module TEXT NOT NULL DEFAULT '',
    target_type TEXT NOT NULL DEFAULT '',
    target_id TEXT NOT NULL DEFAULT '',
    request_id TEXT NOT NULL DEFAULT '',
    correlation_id TEXT NOT NULL DEFAULT '',
    reason TEXT NOT NULL,
    result TEXT NOT NULL,
    risk_level TEXT NOT NULL,
    requires_double_confirmation BOOLEAN NOT NULL DEFAULT FALSE,
    double_confirmed BOOLEAN NOT NULL DEFAULT FALSE,
    confirmation_token TEXT NOT NULL DEFAULT '',
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT admin_actions_action_type_check CHECK (action_type IN ('setting_changed', 'override_applied', 'user_reviewed', 'impersonation_started', 'impersonation_stopped')),
    CONSTRAINT admin_actions_result_check CHECK (result IN ('success')),
    CONSTRAINT admin_actions_risk_level_check CHECK (risk_level IN ('low', 'medium', 'high', 'critical'))
);

CREATE INDEX IF NOT EXISTS idx_admin_actions_type_created
    ON admin_actions (action_type, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_admin_actions_actor_created
    ON admin_actions (actor_user_id, created_at DESC);

CREATE TABLE IF NOT EXISTS admin_action_dedup (
    dedup_key TEXT PRIMARY KEY,
    action_id UUID NOT NULL REFERENCES admin_actions(action_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_admin_action_dedup_created
    ON admin_action_dedup (created_at DESC);

CREATE TABLE IF NOT EXISTS admin_overrides (
    override_id UUID PRIMARY KEY,
    action_id UUID NOT NULL REFERENCES admin_actions(action_id) ON DELETE CASCADE,
    target_module TEXT NOT NULL,
    target_type TEXT NOT NULL,
    target_id TEXT NOT NULL,
    decision TEXT NOT NULL,
    reason TEXT NOT NULL,
    risk_level TEXT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    expires_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT admin_overrides_decision_check CHECK (decision IN ('allow', 'deny', 'freeze', 'reopen')),
    CONSTRAINT admin_overrides_risk_level_check CHECK (risk_level IN ('low', 'medium', 'high', 'critical'))
);

CREATE INDEX IF NOT EXISTS idx_admin_overrides_target_module_created
    ON admin_overrides (target_module, created_at DESC);

CREATE TABLE IF NOT EXISTS admin_user_reviews (
    review_id UUID PRIMARY KEY,
    action_id UUID NOT NULL REFERENCES admin_actions(action_id) ON DELETE CASCADE,
    target_user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    decision TEXT NOT NULL,
    reason TEXT NOT NULL,
    risk_level TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT admin_user_reviews_decision_check CHECK (decision IN ('warning', 'restriction', 'suspend', 'ban', 'clear')),
    CONSTRAINT admin_user_reviews_risk_level_check CHECK (risk_level IN ('low', 'medium', 'high', 'critical'))
);

CREATE INDEX IF NOT EXISTS idx_admin_user_reviews_target_created
    ON admin_user_reviews (target_user_id, created_at DESC);

CREATE TABLE IF NOT EXISTS admin_impersonation_sessions (
    session_id UUID PRIMARY KEY,
    action_id UUID NOT NULL REFERENCES admin_actions(action_id) ON DELETE CASCADE,
    actor_user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    target_user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    reason TEXT NOT NULL,
    risk_level TEXT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ended_at TIMESTAMPTZ NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT admin_impersonation_sessions_risk_level_check CHECK (risk_level IN ('low', 'medium', 'high', 'critical')),
    CONSTRAINT admin_impersonation_sessions_window_check CHECK (ended_at IS NULL OR ended_at >= started_at)
);

CREATE INDEX IF NOT EXISTS idx_admin_impersonation_sessions_active_expires
    ON admin_impersonation_sessions (active, expires_at);
