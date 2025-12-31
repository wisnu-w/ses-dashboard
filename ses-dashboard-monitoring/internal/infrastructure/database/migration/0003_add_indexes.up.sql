-- +migrate Up
CREATE INDEX IF NOT EXISTS idx_ses_events_event_timestamp ON ses_events(event_timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_ses_events_event_type ON ses_events(event_type);
CREATE INDEX IF NOT EXISTS idx_ses_events_email ON ses_events(email);

-- +migrate Down
DROP INDEX IF EXISTS idx_ses_events_event_timestamp;
DROP INDEX IF EXISTS idx_ses_events_event_type;
DROP INDEX IF EXISTS idx_ses_events_email;