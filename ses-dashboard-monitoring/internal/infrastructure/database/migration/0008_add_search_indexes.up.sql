-- Enable pg_trgm extension for better text search
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Add indexes for search optimization
CREATE INDEX IF NOT EXISTS idx_ses_events_email_gin ON ses_events USING gin(email gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_ses_events_subject_gin ON ses_events USING gin(subject gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_ses_events_source_gin ON ses_events USING gin(source gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_ses_events_timestamp ON ses_events(event_timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_ses_events_type ON ses_events(event_type);
CREATE INDEX IF NOT EXISTS idx_ses_events_timestamp_type ON ses_events(event_timestamp DESC, event_type);