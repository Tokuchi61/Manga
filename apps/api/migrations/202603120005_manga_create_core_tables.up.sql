CREATE TABLE IF NOT EXISTS manga_entries (
    id UUID PRIMARY KEY,
    slug TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    alternative_titles TEXT[] NOT NULL DEFAULT '{}',
    summary TEXT NOT NULL,
    short_summary TEXT NOT NULL DEFAULT '',
    cover_image_url TEXT NOT NULL DEFAULT '',
    banner_image_url TEXT NOT NULL DEFAULT '',
    seo_title TEXT NOT NULL DEFAULT '',
    seo_description TEXT NOT NULL DEFAULT '',
    genres TEXT[] NOT NULL DEFAULT '{}',
    tags TEXT[] NOT NULL DEFAULT '{}',
    themes TEXT[] NOT NULL DEFAULT '{}',
    content_warnings TEXT[] NOT NULL DEFAULT '{}',
    publish_state TEXT NOT NULL DEFAULT 'draft',
    visibility TEXT NOT NULL DEFAULT 'public',
    featured BOOLEAN NOT NULL DEFAULT FALSE,
    recommended BOOLEAN NOT NULL DEFAULT FALSE,
    collection_keys TEXT[] NOT NULL DEFAULT '{}',
    default_read_access_level TEXT NOT NULL DEFAULT 'authenticated',
    default_early_access_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    default_early_access_level TEXT NOT NULL DEFAULT 'none',
    release_schedule TEXT NOT NULL DEFAULT '',
    translation_group TEXT NOT NULL DEFAULT '',
    view_count BIGINT NOT NULL DEFAULT 0,
    comment_count BIGINT NOT NULL DEFAULT 0,
    chapter_count BIGINT NOT NULL DEFAULT 0,
    content_version INTEGER NOT NULL DEFAULT 1,
    scheduled_at TIMESTAMPTZ NULL,
    published_at TIMESTAMPTZ NULL,
    archived_at TIMESTAMPTZ NULL,
    deleted_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT manga_entries_publish_state_check CHECK (publish_state IN ('draft', 'scheduled', 'published', 'archived', 'unpublished')),
    CONSTRAINT manga_entries_visibility_check CHECK (visibility IN ('public', 'hidden')),
    CONSTRAINT manga_entries_default_read_access_check CHECK (default_read_access_level IN ('guest', 'authenticated')),
    CONSTRAINT manga_entries_default_early_access_level_check CHECK (default_early_access_level IN ('none', 'vip')),
    CONSTRAINT manga_entries_view_count_check CHECK (view_count >= 0),
    CONSTRAINT manga_entries_comment_count_check CHECK (comment_count >= 0),
    CONSTRAINT manga_entries_chapter_count_check CHECK (chapter_count >= 0),
    CONSTRAINT manga_entries_content_version_check CHECK (content_version > 0)
);

CREATE INDEX IF NOT EXISTS idx_manga_entries_slug ON manga_entries (slug);
CREATE INDEX IF NOT EXISTS idx_manga_entries_publish_state ON manga_entries (publish_state);
CREATE INDEX IF NOT EXISTS idx_manga_entries_visibility ON manga_entries (visibility);
CREATE INDEX IF NOT EXISTS idx_manga_entries_featured ON manga_entries (featured);
CREATE INDEX IF NOT EXISTS idx_manga_entries_recommended ON manga_entries (recommended);
CREATE INDEX IF NOT EXISTS idx_manga_entries_deleted_at ON manga_entries (deleted_at);
