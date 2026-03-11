CREATE TABLE IF NOT EXISTS system_heartbeat (
    id SMALLINT PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO system_heartbeat (id, name)
VALUES (1, 'stage0-bootstrap')
ON CONFLICT (id) DO NOTHING;
