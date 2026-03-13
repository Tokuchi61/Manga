CREATE TABLE IF NOT EXISTS history_library_entries (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    manga_id UUID NOT NULL REFERENCES manga_entries(id) ON DELETE CASCADE,
    last_chapter_id UUID NULL REFERENCES chapter_entries(id) ON DELETE SET NULL,
    last_page_number INTEGER NOT NULL DEFAULT 0,
    page_count INTEGER NOT NULL DEFAULT 0,
    reading_status TEXT NOT NULL DEFAULT 'in_progress',
    bookmarked BOOLEAN NOT NULL DEFAULT FALSE,
    favorited BOOLEAN NOT NULL DEFAULT FALSE,
    share_public BOOLEAN NOT NULL DEFAULT FALSE,
    last_read_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT history_library_entries_user_manga_unique UNIQUE (user_id, manga_id),
    CONSTRAINT history_library_entries_reading_status_check CHECK (reading_status IN ('in_progress', 'completed', 'dropped')),
    CONSTRAINT history_library_entries_last_page_number_check CHECK (last_page_number >= 0),
    CONSTRAINT history_library_entries_page_count_check CHECK (page_count >= 0)
);

CREATE INDEX IF NOT EXISTS idx_history_library_entries_user_updated_at
    ON history_library_entries (user_id, updated_at DESC);

CREATE INDEX IF NOT EXISTS idx_history_library_entries_user_last_read_at
    ON history_library_entries (user_id, last_read_at DESC);

CREATE INDEX IF NOT EXISTS idx_history_library_entries_user_share_public
    ON history_library_entries (user_id, share_public);

CREATE TABLE IF NOT EXISTS history_timeline_events (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    manga_id UUID NOT NULL REFERENCES manga_entries(id) ON DELETE CASCADE,
    chapter_id UUID NOT NULL REFERENCES chapter_entries(id) ON DELETE CASCADE,
    event TEXT NOT NULL,
    page_number INTEGER NOT NULL DEFAULT 0,
    page_count INTEGER NOT NULL DEFAULT 0,
    request_id TEXT NOT NULL DEFAULT '',
    correlation_id TEXT NOT NULL DEFAULT '',
    dedup_key TEXT NOT NULL DEFAULT '',
    occurred_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT history_timeline_events_event_check CHECK (event IN ('chapter.read.started', 'chapter.read.checkpoint', 'chapter.read.finished')),
    CONSTRAINT history_timeline_events_page_number_check CHECK (page_number >= 0),
    CONSTRAINT history_timeline_events_page_count_check CHECK (page_count >= 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_history_timeline_events_dedup_key
    ON history_timeline_events (dedup_key)
    WHERE dedup_key <> '';

CREATE INDEX IF NOT EXISTS idx_history_timeline_events_user_occurred_at
    ON history_timeline_events (user_id, occurred_at DESC);

CREATE TABLE IF NOT EXISTS history_runtime_controls (
    id SMALLINT PRIMARY KEY,
    continue_reading_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    library_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    timeline_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    bookmark_write_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO history_runtime_controls (id)
VALUES (1)
ON CONFLICT (id) DO NOTHING;
