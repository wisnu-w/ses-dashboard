-- Drop search optimization indexes
DROP INDEX IF EXISTS idx_ses_events_email_gin;
DROP INDEX IF EXISTS idx_ses_events_subject_gin;
DROP INDEX IF EXISTS idx_ses_events_source_gin;
DROP INDEX IF EXISTS idx_ses_events_timestamp;
DROP INDEX IF EXISTS idx_ses_events_type;
DROP INDEX IF EXISTS idx_ses_events_timestamp_type;