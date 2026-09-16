CREATE TABLE IF NOT EXISTS timeline_events (
    id           TEXT PRIMARY KEY,
    ticket_id    TEXT NOT NULL,
    event_id     TEXT NOT NULL UNIQUE,
    event_type   TEXT NOT NULL,
    title        TEXT NOT NULL,
    detail       TEXT NOT NULL DEFAULT '',
    actor        TEXT NOT NULL DEFAULT '',
    occurred_at  TIMESTAMPTZ NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_timeline_ticket_occurred
    ON timeline_events (ticket_id, occurred_at);
