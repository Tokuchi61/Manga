CREATE TABLE IF NOT EXISTS inventory_item_definitions (
    item_id TEXT PRIMARY KEY,
    item_type TEXT NOT NULL,
    stackable BOOLEAN NOT NULL DEFAULT TRUE,
    equipable BOOLEAN NOT NULL DEFAULT FALSE,
    consumable BOOLEAN NOT NULL DEFAULT FALSE,
    max_stack INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT inventory_item_definitions_max_stack_check CHECK (max_stack > 0)
);

CREATE INDEX IF NOT EXISTS idx_inventory_item_definitions_type_updated_at
    ON inventory_item_definitions (item_type, updated_at DESC);

CREATE TABLE IF NOT EXISTS inventory_entries (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    item_id TEXT NOT NULL REFERENCES inventory_item_definitions(item_id) ON DELETE RESTRICT,
    quantity INTEGER NOT NULL,
    equipped BOOLEAN NOT NULL DEFAULT FALSE,
    last_source_type TEXT NOT NULL DEFAULT '',
    last_source_ref TEXT NOT NULL DEFAULT '',
    request_id TEXT NOT NULL DEFAULT '',
    correlation_id TEXT NOT NULL DEFAULT '',
    expires_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT inventory_entries_quantity_check CHECK (quantity > 0),
    CONSTRAINT inventory_entries_unique_user_item UNIQUE (user_id, item_id)
);

CREATE INDEX IF NOT EXISTS idx_inventory_entries_user_updated_at
    ON inventory_entries (user_id, updated_at DESC);

CREATE INDEX IF NOT EXISTS idx_inventory_entries_user_equipped
    ON inventory_entries (user_id, equipped);

CREATE TABLE IF NOT EXISTS inventory_grant_dedup (
    dedup_key TEXT PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    item_id TEXT NOT NULL REFERENCES inventory_item_definitions(item_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_inventory_grant_dedup_user_created_at
    ON inventory_grant_dedup (user_id, created_at DESC);

CREATE TABLE IF NOT EXISTS inventory_runtime_controls (
    id SMALLINT PRIMARY KEY,
    read_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    claim_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    consume_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    equip_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO inventory_runtime_controls (id)
VALUES (1)
ON CONFLICT (id) DO NOTHING;