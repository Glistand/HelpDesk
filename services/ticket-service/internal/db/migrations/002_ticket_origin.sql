ALTER TABLE tickets ADD COLUMN IF NOT EXISTS source TEXT NOT NULL DEFAULT 'manager';
ALTER TABLE tickets ADD COLUMN IF NOT EXISTS created_by_id TEXT NOT NULL DEFAULT '';
ALTER TABLE tickets ADD COLUMN IF NOT EXISTS conversation_id TEXT NOT NULL DEFAULT '';
ALTER TABLE tickets ADD COLUMN IF NOT EXISTS creation_reason TEXT NOT NULL DEFAULT '';

CREATE INDEX IF NOT EXISTS idx_tickets_conversation ON tickets (conversation_id) WHERE conversation_id <> '';
CREATE UNIQUE INDEX IF NOT EXISTS idx_tickets_bot_conversation
  ON tickets (conversation_id) WHERE source = 'bot' AND conversation_id <> '';
