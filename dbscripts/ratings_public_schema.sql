CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS ratings(
  record_id TEXT NOT NULL,
  record_type TEXT NOT NULL,
  user_id TEXT,
  value INTEGER NOT NULL
);