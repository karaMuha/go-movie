CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS metadata(
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  title TEXT NOT NULL,
  description TEXT,
  director TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS metadata_created_event(
  id UUID KEY NOT NULL,
  record_type TEXT NOT NULL,
  event_type TEXT NOT NULL
  PRIMARY KEY (id, record_type)
);