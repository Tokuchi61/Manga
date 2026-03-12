CREATE TABLE IF NOT EXISTS support_entries (
    id UUID PRIMARY KEY,
    requester_user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    support_kind TEXT NOT NULL,
    category TEXT NOT NULL,
    priority TEXT NOT NULL DEFAULT 'normal',
    reason_code TEXT NOT NULL DEFAULT '',
    reason_text TEXT NOT NULL,
    target_type TEXT NULL,
    target_id UUID NULL,
    status TEXT NOT NULL DEFAULT 'open',
    duplicate_of_support_id UUID NULL REFERENCES support_entries(id) ON DELETE SET NULL,
    request_id TEXT NOT NULL DEFAULT '',
    spam_risk_score INTEGER NOT NULL DEFAULT 0,
    attachments TEXT[] NOT NULL DEFAULT '{}',
    resolution_note TEXT NOT NULL DEFAULT '',
    assignee_user_id UUID NULL REFERENCES user_accounts(id) ON DELETE SET NULL,
    reviewed_by_user_id UUID NULL REFERENCES user_accounts(id) ON DELETE SET NULL,
    resolved_at TIMESTAMPTZ NULL,
    closed_at TIMESTAMPTZ NULL,
    moderation_handoff_requested_at TIMESTAMPTZ NULL,
    linked_moderation_case_id UUID NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT support_entries_support_kind_check CHECK (support_kind IN ('communication', 'ticket', 'report')),
    CONSTRAINT support_entries_priority_check CHECK (priority IN ('low', 'normal', 'high', 'urgent')),
    CONSTRAINT support_entries_status_check CHECK (status IN ('open', 'triaged', 'waiting_user', 'waiting_team', 'resolved', 'rejected', 'closed', 'spam')),
    CONSTRAINT support_entries_target_type_check CHECK (target_type IS NULL OR target_type IN ('manga', 'chapter', 'comment')),
    CONSTRAINT support_entries_target_pair_check CHECK ((target_type IS NULL AND target_id IS NULL) OR (target_type IS NOT NULL AND target_id IS NOT NULL)),
    CONSTRAINT support_entries_spam_risk_score_check CHECK (spam_risk_score >= 0 AND spam_risk_score <= 100)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_support_entries_requester_request_id
    ON support_entries (requester_user_id, request_id)
    WHERE request_id <> '';

CREATE INDEX IF NOT EXISTS idx_support_entries_requester_user_id ON support_entries (requester_user_id);
CREATE INDEX IF NOT EXISTS idx_support_entries_status ON support_entries (status);
CREATE INDEX IF NOT EXISTS idx_support_entries_priority_created_at ON support_entries (priority, created_at);
CREATE INDEX IF NOT EXISTS idx_support_entries_target ON support_entries (target_type, target_id);

CREATE TABLE IF NOT EXISTS support_replies (
    id UUID PRIMARY KEY,
    support_id UUID NOT NULL REFERENCES support_entries(id) ON DELETE CASCADE,
    author_user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    sanitized_body TEXT NOT NULL,
    visibility TEXT NOT NULL DEFAULT 'public_to_requester',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT support_replies_visibility_check CHECK (visibility IN ('public_to_requester', 'internal_only'))
);

CREATE INDEX IF NOT EXISTS idx_support_replies_support_id_created_at ON support_replies (support_id, created_at);