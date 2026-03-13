CREATE TABLE IF NOT EXISTS payment_mana_packages (
    package_id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    mana_amount INTEGER NOT NULL,
    price_amount INTEGER NOT NULL,
    price_currency TEXT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    display_order INTEGER NOT NULL DEFAULT 0,
    provider TEXT NOT NULL,
    provider_sku TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT payment_mana_packages_mana_amount_check CHECK (mana_amount > 0),
    CONSTRAINT payment_mana_packages_price_amount_check CHECK (price_amount > 0),
    CONSTRAINT payment_mana_packages_display_order_check CHECK (display_order >= 0)
);

CREATE INDEX IF NOT EXISTS idx_payment_mana_packages_active_display_order
    ON payment_mana_packages (active, display_order, updated_at DESC);

CREATE TABLE IF NOT EXISTS payment_transactions (
    transaction_id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    package_id TEXT NOT NULL REFERENCES payment_mana_packages(package_id) ON DELETE RESTRICT,
    amount_mana INTEGER NOT NULL,
    money_amount INTEGER NOT NULL,
    money_currency TEXT NOT NULL,
    source TEXT NOT NULL,
    status TEXT NOT NULL,
    provider TEXT NOT NULL,
    provider_reference TEXT NOT NULL DEFAULT '',
    request_id TEXT NOT NULL DEFAULT '',
    correlation_id TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT payment_transactions_amount_mana_check CHECK (amount_mana > 0),
    CONSTRAINT payment_transactions_money_amount_check CHECK (money_amount > 0),
    CONSTRAINT payment_transactions_source_check CHECK (source IN ('external_provider', 'recovery_replay')),
    CONSTRAINT payment_transactions_status_check CHECK (status IN ('pending', 'success', 'failed', 'cancelled', 'refunded', 'reversed'))
);

CREATE INDEX IF NOT EXISTS idx_payment_transactions_user_status_created_at
    ON payment_transactions (user_id, status, created_at DESC);

CREATE TABLE IF NOT EXISTS payment_provider_sessions (
    session_id UUID PRIMARY KEY,
    transaction_id UUID NOT NULL REFERENCES payment_transactions(transaction_id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    package_id TEXT NOT NULL REFERENCES payment_mana_packages(package_id) ON DELETE RESTRICT,
    provider TEXT NOT NULL,
    provider_reference TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL,
    expires_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT payment_provider_sessions_status_check CHECK (status IN ('started', 'completed', 'cancelled', 'expired'))
);

CREATE INDEX IF NOT EXISTS idx_payment_provider_sessions_transaction_status
    ON payment_provider_sessions (transaction_id, status);

CREATE TABLE IF NOT EXISTS payment_checkout_dedup (
    dedup_key TEXT PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    transaction_id UUID NOT NULL REFERENCES payment_transactions(transaction_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_payment_checkout_dedup_user_created_at
    ON payment_checkout_dedup (user_id, created_at DESC);

CREATE TABLE IF NOT EXISTS payment_callback_dedup (
    provider_event_id TEXT PRIMARY KEY,
    transaction_id UUID NOT NULL REFERENCES payment_transactions(transaction_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_payment_callback_dedup_transaction
    ON payment_callback_dedup (transaction_id, created_at DESC);

CREATE TABLE IF NOT EXISTS payment_ledger_entries (
    entry_id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES user_accounts(id) ON DELETE CASCADE,
    transaction_id UUID NOT NULL REFERENCES payment_transactions(transaction_id) ON DELETE CASCADE,
    entry_type TEXT NOT NULL,
    amount_mana INTEGER NOT NULL,
    reason_code TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT payment_ledger_entries_entry_type_check CHECK (entry_type IN ('credit', 'debit')),
    CONSTRAINT payment_ledger_entries_amount_mana_check CHECK (amount_mana > 0)
);

CREATE INDEX IF NOT EXISTS idx_payment_ledger_entries_user_created_at
    ON payment_ledger_entries (user_id, created_at DESC);

CREATE TABLE IF NOT EXISTS payment_balance_snapshots (
    user_id UUID PRIMARY KEY REFERENCES user_accounts(id) ON DELETE CASCADE,
    balance_mana INTEGER NOT NULL DEFAULT 0,
    last_ledger_entry_id UUID NULL REFERENCES payment_ledger_entries(entry_id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS payment_runtime_controls (
    id SMALLINT PRIMARY KEY,
    mana_purchase_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    checkout_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    transaction_read_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    callback_intake_paused BOOLEAN NOT NULL DEFAULT FALSE,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO payment_runtime_controls (id)
VALUES (1)
ON CONFLICT (id) DO NOTHING;
