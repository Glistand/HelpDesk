CREATE TABLE IF NOT EXISTS conversations (
    id          TEXT PRIMARY KEY,
    site_key    TEXT NOT NULL,
    visitor_id  TEXT NOT NULL,
    status      TEXT NOT NULL DEFAULT 'bot',
    assignee_id TEXT NOT NULL DEFAULT '',
    preview     TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS conversations_status_updated_idx
    ON conversations (status, updated_at DESC);

CREATE TABLE IF NOT EXISTS messages (
    id               TEXT PRIMARY KEY,
    conversation_id  TEXT NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    role             TEXT NOT NULL,
    body             TEXT NOT NULL,
    author_id        TEXT NOT NULL DEFAULT '',
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS messages_conversation_created_idx
    ON messages (conversation_id, created_at);
