CREATE TABLE IF NOT EXISTS shop_products (
    product_id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    state TEXT NOT NULL,
    inventory_item_id TEXT NOT NULL DEFAULT '',
    slot_id TEXT NOT NULL DEFAULT '',
    single_purchase BOOLEAN NOT NULL DEFAULT FALSE,
    vip_required BOOLEAN NOT NULL DEFAULT FALSE,
    min_level INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT shop_products_state_check CHECK (state IN ('draft', 'active', 'archived')),
    CONSTRAINT shop_products_min_level_check CHECK (min_level >= 0)
);

CREATE INDEX IF NOT EXISTS idx_shop_products_state_updated_at
    ON shop_products (state, updated_at DESC);

CREATE TABLE IF NOT EXISTS shop_offers (
    offer_id TEXT PRIMARY KEY,
    product_id TEXT NOT NULL REFERENCES shop_products(product_id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    visibility TEXT NOT NULL,
    price_mana INTEGER NOT NULL,
    discount_percent INTEGER NOT NULL DEFAULT 0,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    starts_at TIMESTAMPTZ NULL,
    ends_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT shop_offers_visibility_check CHECK (visibility IN ('visible', 'hidden', 'campaign_only')),
    CONSTRAINT shop_offers_price_mana_check CHECK (price_mana > 0),
    CONSTRAINT shop_offers_discount_percent_check CHECK (discount_percent >= 0 AND discount_percent <= 100),
    CONSTRAINT shop_offers_window_check CHECK (starts_at IS NULL OR ends_at IS NULL OR starts_at <= ends_at)
);

CREATE INDEX IF NOT EXISTS idx_shop_offers_product_visibility_active
    ON shop_offers (product_id, visibility, active);

CREATE TABLE IF NOT EXISTS shop_purchase_intents (
    intent_id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    product_id TEXT NOT NULL REFERENCES shop_products(product_id) ON DELETE CASCADE,
    offer_id TEXT NOT NULL REFERENCES shop_offers(offer_id) ON DELETE CASCADE,
    final_price_mana INTEGER NOT NULL,
    currency TEXT NOT NULL DEFAULT 'mana',
    status TEXT NOT NULL,
    request_id TEXT NOT NULL DEFAULT '',
    correlation_id TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT shop_purchase_intents_final_price_check CHECK (final_price_mana > 0),
    CONSTRAINT shop_purchase_intents_currency_check CHECK (currency IN ('mana')),
    CONSTRAINT shop_purchase_intents_status_check CHECK (status IN ('delivery_pending', 'completed', 'recovery_required', 'blocked'))
);

CREATE INDEX IF NOT EXISTS idx_shop_purchase_intents_user_created_at
    ON shop_purchase_intents (user_id, created_at DESC);

CREATE TABLE IF NOT EXISTS shop_purchase_dedup (
    dedup_key TEXT PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    product_id TEXT NOT NULL REFERENCES shop_products(product_id) ON DELETE CASCADE,
    offer_id TEXT NOT NULL REFERENCES shop_offers(offer_id) ON DELETE CASCADE,
    intent_id UUID NOT NULL REFERENCES shop_purchase_intents(intent_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_shop_purchase_dedup_user_created_at
    ON shop_purchase_dedup (user_id, created_at DESC);

CREATE TABLE IF NOT EXISTS shop_runtime_controls (
    id SMALLINT PRIMARY KEY,
    catalog_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    purchase_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    campaign_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO shop_runtime_controls (id)
VALUES (1)
ON CONFLICT (id) DO NOTHING;
