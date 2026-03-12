CREATE TABLE IF NOT EXISTS chapter_entries (
    id UUID PRIMARY KEY,
    manga_id UUID NOT NULL REFERENCES manga_entries(id) ON DELETE CASCADE,
    slug TEXT NOT NULL,
    title TEXT NOT NULL,
    summary TEXT NOT NULL DEFAULT '',
    sequence_no INTEGER NOT NULL,
    display_number TEXT NOT NULL DEFAULT '',
    publish_state TEXT NOT NULL DEFAULT 'draft',
    read_access_level TEXT NOT NULL DEFAULT 'authenticated',
    inherit_access_from_manga BOOLEAN NOT NULL DEFAULT FALSE,
    vip_only BOOLEAN NOT NULL DEFAULT FALSE,
    early_access_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    early_access_level TEXT NOT NULL DEFAULT 'none',
    early_access_start_at TIMESTAMPTZ NULL,
    early_access_end_at TIMESTAMPTZ NULL,
    early_access_fallback_access TEXT NOT NULL DEFAULT 'authenticated',
    preview_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    preview_page_count INTEGER NOT NULL DEFAULT 0,
    media_health_status TEXT NOT NULL DEFAULT 'healthy',
    integrity_status TEXT NOT NULL DEFAULT 'unknown',
    page_count INTEGER NOT NULL DEFAULT 0,
    scheduled_at TIMESTAMPTZ NULL,
    published_at TIMESTAMPTZ NULL,
    archived_at TIMESTAMPTZ NULL,
    deleted_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chapter_entries_manga_slug_unique UNIQUE (manga_id, slug),
    CONSTRAINT chapter_entries_manga_sequence_unique UNIQUE (manga_id, sequence_no),
    CONSTRAINT chapter_entries_publish_state_check CHECK (publish_state IN ('draft', 'scheduled', 'published', 'archived', 'unpublished')),
    CONSTRAINT chapter_entries_read_access_check CHECK (read_access_level IN ('guest', 'authenticated', 'vip')),
    CONSTRAINT chapter_entries_early_access_level_check CHECK (early_access_level IN ('none', 'vip')),
    CONSTRAINT chapter_entries_early_access_fallback_check CHECK (early_access_fallback_access IN ('guest', 'authenticated', 'vip')),
    CONSTRAINT chapter_entries_media_health_status_check CHECK (media_health_status IN ('healthy', 'degraded', 'broken')),
    CONSTRAINT chapter_entries_integrity_status_check CHECK (integrity_status IN ('unknown', 'passed', 'failed')),
    CONSTRAINT chapter_entries_sequence_no_check CHECK (sequence_no > 0),
    CONSTRAINT chapter_entries_page_count_check CHECK (page_count >= 0),
    CONSTRAINT chapter_entries_preview_page_count_check CHECK (preview_page_count >= 0)
);

CREATE INDEX IF NOT EXISTS idx_chapter_entries_manga_id ON chapter_entries (manga_id);
CREATE INDEX IF NOT EXISTS idx_chapter_entries_publish_state ON chapter_entries (publish_state);
CREATE INDEX IF NOT EXISTS idx_chapter_entries_deleted_at ON chapter_entries (deleted_at);
CREATE INDEX IF NOT EXISTS idx_chapter_entries_sequence_no ON chapter_entries (sequence_no);

CREATE TABLE IF NOT EXISTS chapter_pages (
    id UUID PRIMARY KEY,
    chapter_id UUID NOT NULL REFERENCES chapter_entries(id) ON DELETE CASCADE,
    page_number INTEGER NOT NULL,
    media_url TEXT NOT NULL,
    width INTEGER NOT NULL,
    height INTEGER NOT NULL,
    long_strip BOOLEAN NOT NULL DEFAULT FALSE,
    checksum TEXT NOT NULL DEFAULT '',
    cdn_healthy BOOLEAN NOT NULL DEFAULT TRUE,
    missing BOOLEAN NOT NULL DEFAULT FALSE,
    broken BOOLEAN NOT NULL DEFAULT FALSE,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chapter_pages_chapter_page_unique UNIQUE (chapter_id, page_number),
    CONSTRAINT chapter_pages_page_number_check CHECK (page_number > 0),
    CONSTRAINT chapter_pages_width_check CHECK (width > 0),
    CONSTRAINT chapter_pages_height_check CHECK (height > 0)
);

CREATE INDEX IF NOT EXISTS idx_chapter_pages_chapter_id ON chapter_pages (chapter_id);
