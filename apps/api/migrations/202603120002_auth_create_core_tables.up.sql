CREATE TABLE IF NOT EXISTS auth_credentials (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    suspended BOOLEAN NOT NULL DEFAULT FALSE,
    banned BOOLEAN NOT NULL DEFAULT FALSE,
    failed_login_attempts INTEGER NOT NULL DEFAULT 0,
    login_cooldown_until TIMESTAMPTZ NULL,
    verification_resend_available_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS auth_sessions (
    id UUID PRIMARY KEY,
    credential_id UUID NOT NULL REFERENCES auth_credentials(id) ON DELETE CASCADE,
    device TEXT NOT NULL DEFAULT '',
    ip TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_seen_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ NULL
);
CREATE INDEX IF NOT EXISTS idx_auth_sessions_credential_id ON auth_sessions (credential_id);
CREATE INDEX IF NOT EXISTS idx_auth_sessions_revoked_at ON auth_sessions (revoked_at);

CREATE TABLE IF NOT EXISTS auth_tokens (
    id UUID PRIMARY KEY,
    credential_id UUID NOT NULL REFERENCES auth_credentials(id) ON DELETE CASCADE,
    session_id UUID NULL REFERENCES auth_sessions(id) ON DELETE CASCADE,
    token_type TEXT NOT NULL,
    token_hash TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    consumed_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT auth_tokens_type_check CHECK (token_type IN ('refresh', 'password_reset', 'email_verification')),
    CONSTRAINT auth_tokens_type_hash_unique UNIQUE (token_type, token_hash)
);
CREATE INDEX IF NOT EXISTS idx_auth_tokens_credential_id ON auth_tokens (credential_id);
CREATE INDEX IF NOT EXISTS idx_auth_tokens_session_id ON auth_tokens (session_id);
CREATE INDEX IF NOT EXISTS idx_auth_tokens_expires_at ON auth_tokens (expires_at);

CREATE TABLE IF NOT EXISTS auth_security_events (
    id UUID PRIMARY KEY,
    credential_id UUID NULL REFERENCES auth_credentials(id) ON DELETE SET NULL,
    actor_id UUID NULL,
    target_id UUID NULL,
    action TEXT NOT NULL,
    result TEXT NOT NULL,
    reason TEXT NOT NULL,
    request_id TEXT NOT NULL,
    correlation_id TEXT NOT NULL,
    device TEXT NOT NULL DEFAULT '',
    ip TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_auth_security_events_credential_id ON auth_security_events (credential_id);
CREATE INDEX IF NOT EXISTS idx_auth_security_events_action ON auth_security_events (action);
CREATE INDEX IF NOT EXISTS idx_auth_security_events_created_at ON auth_security_events (created_at);
