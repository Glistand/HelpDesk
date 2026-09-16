CREATE TABLE IF NOT EXISTS escalations (
    ticket_id      TEXT PRIMARY KEY,
    reason         TEXT NOT NULL,
    from_assignee  TEXT NOT NULL DEFAULT '',
    to_assignee    TEXT NOT NULL,
    escalated_at   TIMESTAMPTZ NOT NULL,
    source_event   TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS processed_events (
    event_id     TEXT PRIMARY KEY,
    processed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
