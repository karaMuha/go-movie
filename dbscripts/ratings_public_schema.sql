CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS ratings(
  record_id TEXT NOT NULL,
  record_type TEXT NOT NULL,
  user_id TEXT,
  value INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS aggregated_ratings(
  id UUID NOT NULL,
  record_type TEXT NOT NULL,
  rating DECIMAL(2, 1) DEFAULT 0,
  amount_ratings INTEGER DEFAULT 0,
  PRIMARY KEY (id, record_type)
);

CREATE TABLE IF NOT EXISTS metadata_created_event(
  id UUID NOT NULL,
  record_type TEXT NOT NULL,
  event_type TEXT NOT NULL,
  PRIMARY KEY (id, record_type)
);