-- Basic indexes for ses_events table
CREATE INDEX IF NOT EXISTS idx_ses_events_created_at ON ses_events(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_ses_events_message_id ON ses_events(message_id);

-- +migrate Down
DROP INDEX IF EXISTS idx_ses_events_event_type;
DROP INDEX IF EXISTS idx_ses_events_email;
