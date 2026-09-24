CREATE TABLE IF NOT EXISTS project_profile (
  singleton BOOLEAN PRIMARY KEY DEFAULT TRUE CHECK (singleton),
  name TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  domain TEXT NOT NULL DEFAULT '',
  website_url TEXT NOT NULL DEFAULT '',
  support_email TEXT NOT NULL DEFAULT '',
  support_phone TEXT NOT NULL DEFAULT '',
  timezone TEXT NOT NULL DEFAULT 'UTC',
  site_key TEXT NOT NULL,
  bot_instructions TEXT NOT NULL DEFAULT '',
  auto_create_ticket BOOLEAN NOT NULL DEFAULT TRUE,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
