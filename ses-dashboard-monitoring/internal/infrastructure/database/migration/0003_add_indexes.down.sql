-- +migrate Up

-- +migrate Down
DROP INDEX IF EXISTS idx_ses_events_event_timestamp;
DROP INDEX IF EXISTS idx_ses_events_event_type;
DROP INDEX IF EXISTS idx_ses_events_email;