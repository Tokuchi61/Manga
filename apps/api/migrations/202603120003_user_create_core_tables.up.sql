CREATE TABLE IF NOT EXISTS user_accounts (
    id UUID PRIMARY KEY,
    credential_id UUID NOT NULL UNIQUE REFERENCES auth_credentials(id) ON DELETE CASCADE,
    username TEXT NOT NULL UNIQUE,
    display_name TEXT NOT NULL DEFAULT '',
    bio TEXT NOT NULL DEFAULT '',
    avatar_url TEXT NOT NULL DEFAULT '',
    banner_url TEXT NOT NULL DEFAULT '',
    profile_visibility TEXT NOT NULL DEFAULT 'public',
    history_visibility_preference TEXT NOT NULL DEFAULT 'private',
    account_state TEXT NOT NULL DEFAULT 'active',
    vip_active BOOLEAN NOT NULL DEFAULT FALSE,
    vip_frozen BOOLEAN NOT NULL DEFAULT FALSE,
    vip_started_at TIMESTAMPTZ NULL,
    vip_ends_at TIMESTAMPTZ NULL,
    vip_frozen_at TIMESTAMPTZ NULL,
    vip_freeze_reason TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT user_accounts_profile_visibility_check CHECK (profile_visibility IN ('public', 'private')),
    CONSTRAINT user_accounts_history_visibility_check CHECK (history_visibility_preference IN ('public', 'private')),
    CONSTRAINT user_accounts_account_state_check CHECK (account_state IN ('active', 'deactivated', 'banned'))
);
CREATE INDEX IF NOT EXISTS idx_user_accounts_credential_id ON user_accounts (credential_id);
CREATE INDEX IF NOT EXISTS idx_user_accounts_username ON user_accounts (username);
CREATE INDEX IF NOT EXISTS idx_user_accounts_account_state ON user_accounts (account_state);
