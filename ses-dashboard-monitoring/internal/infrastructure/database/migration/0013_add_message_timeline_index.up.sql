CREATE INDEX IF NOT EXISTS idx_ses_events_message_timestamp ON ses_events(message_id, event_timestamp DESC);
