CREATE TABLE IF NOT EXISTS access_roles (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    priority INTEGER NOT NULL DEFAULT 0,
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    is_super_admin BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS access_permissions (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    module TEXT NOT NULL,
    surface TEXT NOT NULL,
    action TEXT NOT NULL,
    audience_kind TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT access_permissions_audience_kind_check CHECK (audience_kind IN ('all', 'guest', 'authenticated', 'authenticated_non_vip', 'vip'))
);

CREATE TABLE IF NOT EXISTS access_role_permissions (
    role_id UUID NOT NULL REFERENCES access_roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES access_permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE IF NOT EXISTS access_user_roles (
    user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES access_roles(id) ON DELETE CASCADE,
    expires_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, role_id)
);
CREATE INDEX IF NOT EXISTS idx_access_user_roles_user_id ON access_user_roles (user_id);
CREATE INDEX IF NOT EXISTS idx_access_user_roles_expires_at ON access_user_roles (expires_at);

CREATE TABLE IF NOT EXISTS access_policy_rules (
    id UUID PRIMARY KEY,
    key TEXT NOT NULL,
    effect TEXT NOT NULL,
    audience_kind TEXT NOT NULL,
    audience_selector TEXT NOT NULL,
    scope_kind TEXT NOT NULL,
    scope_selector TEXT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    version INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT access_policy_rules_effect_check CHECK (effect IN ('allow', 'deny', 'emergency_deny')),
    CONSTRAINT access_policy_rules_audience_kind_check CHECK (audience_kind IN ('all', 'guest', 'authenticated', 'authenticated_non_vip', 'vip')),
    CONSTRAINT access_policy_rules_scope_kind_check CHECK (scope_kind IN ('site', 'module', 'feature', 'resource/context'))
);
CREATE INDEX IF NOT EXISTS idx_access_policy_rules_key ON access_policy_rules (key);
CREATE INDEX IF NOT EXISTS idx_access_policy_rules_active ON access_policy_rules (active);
CREATE UNIQUE INDEX IF NOT EXISTS idx_access_policy_rules_active_conflict
    ON access_policy_rules (key, audience_kind, audience_selector, scope_kind, scope_selector)
    WHERE active = TRUE;

CREATE TABLE IF NOT EXISTS access_temporary_grants (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES access_permissions(id) ON DELETE CASCADE,
    reason TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_access_temporary_grants_user_id ON access_temporary_grants (user_id);
CREATE INDEX IF NOT EXISTS idx_access_temporary_grants_expires_at ON access_temporary_grants (expires_at);
