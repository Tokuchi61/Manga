CREATE TABLE IF NOT EXISTS app_outbox_events (
    event_id UUID PRIMARY KEY,
    topic TEXT NOT NULL,
    payload JSONB NOT NULL,
    request_id TEXT NOT NULL,
    correlation_id TEXT,
    causation_id TEXT,
    status TEXT NOT NULL DEFAULT 'pending',
    attempts INTEGER NOT NULL DEFAULT 0,
    next_attempt_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    leased_until TIMESTAMPTZ,
    last_error TEXT,
    dead_letter BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_app_outbox_events_pending
    ON app_outbox_events (status, dead_letter, next_attempt_at);