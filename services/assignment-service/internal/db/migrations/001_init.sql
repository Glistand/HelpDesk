CREATE TABLE IF NOT EXISTS agents (
    id    TEXT PRIMARY KEY,
    name  TEXT NOT NULL,
    role  TEXT NOT NULL DEFAULT 'L1 Support'
);

CREATE TABLE IF NOT EXISTS assignments (
    ticket_id     TEXT PRIMARY KEY,
    assignee_id   TEXT NOT NULL,
    assignee_name TEXT NOT NULL,
    assigned_at   TIMESTAMPTZ NOT NULL,
    source_event  TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS processed_events (
    event_id    TEXT PRIMARY KEY,
    processed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS rr_state (
    queue TEXT PRIMARY KEY,
    idx   INT NOT NULL DEFAULT 0
);

INSERT INTO agents (id, name, role) VALUES
    ('a-1', 'Алексей К.', 'L1 Support')
ON CONFLICT (id) DO NOTHING;

INSERT INTO rr_state (queue, idx) VALUES ('l1', 0)
ON CONFLICT (queue) DO NOTHING;
