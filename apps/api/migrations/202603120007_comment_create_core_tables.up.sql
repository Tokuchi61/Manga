CREATE TABLE IF NOT EXISTS comment_entries (
    id UUID PRIMARY KEY,
    target_type TEXT NOT NULL,
    target_id UUID NOT NULL,
    author_user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    parent_comment_id UUID NULL REFERENCES comment_entries(id) ON DELETE SET NULL,
    root_comment_id UUID NULL REFERENCES comment_entries(id) ON DELETE SET NULL,
    depth INTEGER NOT NULL DEFAULT 0,
    content TEXT NOT NULL,
    sanitized_content TEXT NOT NULL,
    attachments TEXT[] NOT NULL DEFAULT '{}',
    spoiler BOOLEAN NOT NULL DEFAULT FALSE,
    pinned BOOLEAN NOT NULL DEFAULT FALSE,
    locked BOOLEAN NOT NULL DEFAULT FALSE,
    moderation_status TEXT NOT NULL DEFAULT 'visible',
    shadowbanned BOOLEAN NOT NULL DEFAULT FALSE,
    spam_risk_score INTEGER NOT NULL DEFAULT 0,
    like_count BIGINT NOT NULL DEFAULT 0,
    reply_count BIGINT NOT NULL DEFAULT 0,
    edit_count INTEGER NOT NULL DEFAULT 0,
    edited_at TIMESTAMPTZ NULL,
    deleted_at TIMESTAMPTZ NULL,
    delete_reason TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT comment_entries_target_type_check CHECK (target_type IN ('manga', 'chapter')),
    CONSTRAINT comment_entries_depth_check CHECK (depth >= 0 AND depth <= 3),
    CONSTRAINT comment_entries_moderation_status_check CHECK (moderation_status IN ('visible', 'hidden', 'flagged')),
    CONSTRAINT comment_entries_spam_risk_score_check CHECK (spam_risk_score >= 0 AND spam_risk_score <= 100),
    CONSTRAINT comment_entries_like_count_check CHECK (like_count >= 0),
    CONSTRAINT comment_entries_reply_count_check CHECK (reply_count >= 0),
    CONSTRAINT comment_entries_edit_count_check CHECK (edit_count >= 0)
);

CREATE INDEX IF NOT EXISTS idx_comment_entries_target ON comment_entries (target_type, target_id);
CREATE INDEX IF NOT EXISTS idx_comment_entries_parent_comment_id ON comment_entries (parent_comment_id);
CREATE INDEX IF NOT EXISTS idx_comment_entries_root_comment_id ON comment_entries (root_comment_id);
CREATE INDEX IF NOT EXISTS idx_comment_entries_created_at ON comment_entries (created_at);
CREATE INDEX IF NOT EXISTS idx_comment_entries_moderation_status ON comment_entries (moderation_status);
CREATE INDEX IF NOT EXISTS idx_comment_entries_deleted_at ON comment_entries (deleted_at);
CREATE INDEX IF NOT EXISTS idx_comment_entries_author_user_id ON comment_entries (author_user_id);
