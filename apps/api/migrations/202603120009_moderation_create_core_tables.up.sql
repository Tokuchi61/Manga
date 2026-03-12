CREATE TABLE IF NOT EXISTS moderation_cases (
    id UUID PRIMARY KEY,
    source TEXT NOT NULL DEFAULT 'support_report',
    source_ref_id UUID NULL,
    request_id TEXT NOT NULL DEFAULT '',
    correlation_id TEXT NOT NULL DEFAULT '',
    target_type TEXT NOT NULL,
    target_id UUID NOT NULL,
    reporter_user_id UUID NULL REFERENCES user_accounts(id) ON DELETE SET NULL,
    case_status TEXT NOT NULL DEFAULT 'new',
    assignment_status TEXT NOT NULL DEFAULT 'unassigned',
    assigned_moderator_user_id UUID NULL REFERENCES user_accounts(id) ON DELETE SET NULL,
    escalation_status TEXT NOT NULL DEFAULT 'not_escalated',
    escalation_reason TEXT NOT NULL DEFAULT '',
    escalated_at TIMESTAMPTZ NULL,
    action_result TEXT NOT NULL DEFAULT 'none',
    last_action_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT moderation_cases_source_check CHECK (source IN ('support_report', 'manual')),
    CONSTRAINT moderation_cases_target_type_check CHECK (target_type IN ('manga', 'chapter', 'comment')),
    CONSTRAINT moderation_cases_status_check CHECK (case_status IN ('new', 'queued', 'assigned', 'in_review', 'escalated', 'resolved', 'rejected', 'closed')),
    CONSTRAINT moderation_cases_assignment_status_check CHECK (assignment_status IN ('unassigned', 'assigned', 'handoff_pending', 'released')),
    CONSTRAINT moderation_cases_escalation_status_check CHECK (escalation_status IN ('not_escalated', 'pending_admin', 'escalated', 'resolved')),
    CONSTRAINT moderation_cases_action_result_check CHECK (action_result IN ('none', 'content_hidden', 'content_restored', 'warning_sent', 'no_action'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_moderation_cases_source_ref
    ON moderation_cases (source, source_ref_id)
    WHERE source_ref_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_moderation_cases_status ON moderation_cases (case_status);
CREATE INDEX IF NOT EXISTS idx_moderation_cases_target ON moderation_cases (target_type, target_id);
CREATE INDEX IF NOT EXISTS idx_moderation_cases_assignment ON moderation_cases (assignment_status, assigned_moderator_user_id);

CREATE TABLE IF NOT EXISTS moderation_case_notes (
    id UUID PRIMARY KEY,
    case_id UUID NOT NULL REFERENCES moderation_cases(id) ON DELETE CASCADE,
    author_user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    body TEXT NOT NULL,
    internal_only BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_moderation_case_notes_case_id_created_at
    ON moderation_case_notes (case_id, created_at);

CREATE TABLE IF NOT EXISTS moderation_case_actions (
    id UUID PRIMARY KEY,
    case_id UUID NOT NULL REFERENCES moderation_cases(id) ON DELETE CASCADE,
    actor_user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    action_type TEXT NOT NULL,
    reason_code TEXT NOT NULL DEFAULT '',
    summary TEXT NOT NULL DEFAULT '',
    action_result TEXT NOT NULL DEFAULT 'none',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT moderation_case_actions_action_type_check CHECK (action_type IN ('hide', 'unhide', 'lock', 'unlock', 'warning', 'review_complete', 'escalate')),
    CONSTRAINT moderation_case_actions_action_result_check CHECK (action_result IN ('none', 'content_hidden', 'content_restored', 'warning_sent', 'no_action'))
);

CREATE INDEX IF NOT EXISTS idx_moderation_case_actions_case_id_created_at
    ON moderation_case_actions (case_id, created_at);
