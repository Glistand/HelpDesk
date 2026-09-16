CREATE TABLE IF NOT EXISTS ticket_sla (
    ticket_id            TEXT PRIMARY KEY,
    state                TEXT NOT NULL DEFAULT 'ok',
    policy               TEXT NOT NULL DEFAULT 'default',
    first_response_due   TIMESTAMPTZ NOT NULL,
    resolve_due          TIMESTAMPTZ NOT NULL,
    warned_at            TIMESTAMPTZ,
    breached_at          TIMESTAMPTZ,
    cancelled_at         TIMESTAMPTZ,
    correlation_id       TEXT NOT NULL DEFAULT '',
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS processed_events (
    event_id     TEXT PRIMARY KEY,
    processed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS fired_timers (
    ticket_id TEXT NOT NULL,
    kind      TEXT NOT NULL,
    fired_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (ticket_id, kind)
);
